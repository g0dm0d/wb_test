package main_test

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/g0dm0d/wbtest/internal/server/req"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/g0dm0d/wbtest/internal/service/nats"
	"github.com/g0dm0d/wbtest/internal/service/order"
	"github.com/g0dm0d/wbtest/internal/store"

	"github.com/g0dm0d/wbtest/internal/store/postgres"
	"github.com/g0dm0d/wbtest/pkg/cache"
	"github.com/golang-migrate/migrate/v4"
	pg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB
var sc stan.Conn

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resourcePG, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.1-alpine3.18",
		Env: []string{
			"POSTGRES_PASSWORD=12345",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	resourceNS, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "nats-streaming",
		Tag:        "0.25.6",
		Env: []string{
			"NATS_STREAMING_ID=streaming-server",
			"NATS_CLUSTER_ID=service-test-cluster",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPortNS := resourceNS.GetHostPort("4222/tcp")

	log.Println("Connecting to nats on url: ", fmt.Sprint("nats://", hostAndPortNS))

	resourceNS.Expire(120)
	resourcePG.Expire(120)

	if err != nil {
		log.Fatal(err)
	}

	if err = pool.Retry(func() error {
		sc, err = stan.Connect("test-cluster", "client-test", stan.NatsURL(fmt.Sprint("nats://", hostAndPortNS)))
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer sc.Close()

	hostAndPortPG := resourcePG.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:12345@%s/postgres?sslmode=disable", hostAndPortPG)

	log.Println("Connecting to database on url: ", databaseUrl)

	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	driver, err := pg.WithInstance(db, &pg.Config{})
	gom, err := migrate.NewWithDatabaseInstance("file://internal/store/migrations/", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = gom.Up()
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := pool.Purge(resourcePG); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestNatsHandler(t *testing.T) {
	orderStore := postgres.NewOrderStore(db)
	cacheMap := cache.NewCacheMap()

	nats := nats.New(orderStore, cacheMap)

	rightResp, err := os.ReadFile("model.json")
	if err != nil {
		t.Error(err)
	}

	done := make(chan bool)
	sub, _ := sc.Subscribe("order.pipeline", func(msg *stan.Msg) {
		defer func() { done <- true }()
		nats.HandleData(msg)
	})

	err = sc.Publish("order.pipeline", rightResp)
	if err != nil {
		t.Error(err)
	}

	<-done

	defer sub.Unsubscribe()

	_, ok := cacheMap.Get("b563feb7b2b84b6test")
	if !ok {
		t.Error("value in cache not found")
	}

	_, err = db.Exec("DELETE FROM Orders")
	if err != nil {
		t.Error(err)
	}
}

func TestGetOrderFromDB(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/{orderID}", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orderID", "b563feb7b2b84b6test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	rightResp, err := os.ReadFile("model.json")
	if err != nil {
		t.Error(err)
	}

	orderStore := postgres.NewOrderStore(db)

	cacheMap := cache.NewCacheMap()

	orderStore.SaveOrder(store.SaveOrderOpts{Jsonb: rightResp})

	order := order.New(orderStore, cacheMap)
	order.GetOrder(&req.Ctx{
		Writer:  w,
		Request: r,
	})

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}

	require.JSONEq(t, string(rightResp), string(data))

	_, err = db.Exec("DELETE FROM Orders")
	if err != nil {
		t.Error(err)
	}
}
func TestGetOrderFromCache(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/{orderID}", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orderID", "b563feb7b2b84b6test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	rightResp, err := os.ReadFile("model.json")
	if err != nil {
		t.Error(err)
	}

	orderStore := postgres.NewOrderStore(db)

	cacheMap := cache.NewCacheMap()

	cacheMap.Set("b563feb7b2b84b6test", rightResp)

	order := order.New(orderStore, cacheMap)
	order.GetOrder(&req.Ctx{
		Writer:  w,
		Request: r,
	})

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.JSONEq(t, string(rightResp), string(data))
}

func TestGetOrderFromCacheNoDB(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/{orderID}", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("orderID", "b563feb7b2b84b6test")

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	rightResp, err := os.ReadFile("model.json")
	if err != nil {
		t.Error(err)
	}

	var orderStore store.OrderStore

	cacheMap := cache.NewCacheMap()

	cacheMap.Set("b563feb7b2b84b6test", rightResp)

	order := order.New(orderStore, cacheMap)
	order.GetOrder(&req.Ctx{
		Writer:  w,
		Request: r,
	})

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	require.JSONEq(t, string(rightResp), string(data))
}
