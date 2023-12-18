package server

import (
	"fmt"
	"net/http"

	"github.com/g0dm0d/wbtest/internal/server/req"
	"github.com/g0dm0d/wbtest/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nats-io/stan.go"
)

type Server struct {
	server *http.Server
	router chi.Router

	service *service.Service
	sc      stan.Conn
}

type Config struct {
	Addr     string
	Port     int
	Service  *service.Service
	StanConn stan.Conn
}

func NewServer(config *Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprint(config.Addr, ":", config.Port),
			Handler: http.NotFoundHandler(),
		},
		router: chi.NewRouter(),

		service: config.Service,
		sc:      config.StanConn,
	}
}

func (s *Server) SetupRouter() {
	s.setupCors()

	s.router.Route("/", func(r chi.Router) {
		r.Method("GET", "/{orderID}", req.NewHandler(s.service.Order.GetOrder))
	})

	s.server.Handler = s.router
}

func (s *Server) RunServer() error {
	sub, _ := s.sc.Subscribe("order.pipeline", s.service.Nats.HandleData)

	defer sub.Unsubscribe()
	defer s.sc.Close()

	return s.server.ListenAndServe()
}

func (s *Server) setupCors() {
	s.router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)
}
