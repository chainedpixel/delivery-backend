package server

import (
	"bootstrap"
	"config"
	"github.com/gorilla/mux"
	"infrastructure/api/routes"
	"net/http"
	"shared/logs"
	"time"
)

type Server struct {
	router    *mux.Router
	container *bootstrap.Container
	config    *config.EnvConfig

	privatePath string
	publicPath  string
}

func NewAPIServer(container *bootstrap.Container, config *config.EnvConfig) *Server {
	return &Server{
		router:      mux.NewRouter(),
		container:   container,
		config:      config,
		privatePath: "/api/v1",
		publicPath:  "/api/v1",
	}
}

func (s *Server) Start() error {
	err := s.container.Initialize()
	if err != nil {
		return err
	}

	s.configureRoutes()
	server := &http.Server{
		Handler:      s.router,
		Addr:         ":" + s.config.Server.Port,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logs.Info("Server started successfully", map[string]interface{}{
		"port": s.config.Server.Port,
	})
	if err = server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) configureRoutes() {
	s.configureGlobalMiddlewares(s.router)
	routes.RegisterSwaggerRoutes(s.router)

	public := s.router.PathPrefix(s.publicPath).Subrouter()
	private := s.router.PathPrefix(s.privatePath).Subrouter()

	s.configurePublicRoutes(public)
	s.configureProtectedRoutes(private)

	logs.Info("Routes configured successfully", map[string]interface{}{
		"publicPath":    s.publicPath,
		"protectedPath": s.privatePath,
	})
}

func (s *Server) configurePublicRoutes(router *mux.Router) {
	routes.RegisterPublicAuthRoutes(router, s.container.GetHandlerContainer().GetAuthHandler())

}

func (s *Server) configureProtectedRoutes(router *mux.Router) {
	s.configureProtectedMiddlewares(router)

	routes.RegisterProtectedAuthRoutes(router, s.container.GetHandlerContainer().GetAuthHandler())
	routes.RegisterUserRoutes(router, s.container.GetHandlerContainer().GetUserHandler())
	routes.RegisterOrderRoutes(router, s.container.GetHandlerContainer().GetOrderHandler())
	routes.RegisterRoleRoutes(router, s.container.GetHandlerContainer().GetRoleHandler())
}

func (s *Server) configureGlobalMiddlewares(router *mux.Router) {
	router.Use(s.container.GetMiddlewareContainer().GetErrorMiddleware().Handler)
	router.Use(s.container.GetMiddlewareContainer().GetCorsMiddleware().Handler)
}

func (s *Server) configureProtectedMiddlewares(router *mux.Router) {
	router.Use(s.container.GetMiddlewareContainer().GetAuthMiddleware().Handle)
	router.Use(s.container.GetMiddlewareContainer().GetTokenExtractor().ExtractToken)
}
