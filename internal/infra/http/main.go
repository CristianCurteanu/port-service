package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	server   *http.Server
	listener net.Listener
}

func NewApp(listener net.Listener, handler http.Handler) *App {
	return &App{
		listener: listener,
		server: &http.Server{
			Handler: handler,
		},
	}
}

func (a *App) Run() error {
	return a.server.Serve(a.listener)
}

func (a *App) CloseWithContext(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.CloseWithContext(ctx)
}

type DomainHandler struct {
	Path        string
	Middlewares []gin.HandlerFunc
	Routes      []Route
}

type Route struct {
	Path        string
	Method      string
	Middlewares []gin.HandlerFunc
	Handler     gin.HandlerFunc
}

func NewRouter(groups ...DomainHandler) http.Handler {
	router := gin.Default()

	for _, group := range groups {
		newGroup := router.Group(group.Path)
		newGroup.Use(group.Middlewares...)

		for _, route := range group.Routes {
			newGroup.Use(route.Middlewares...)
			newGroup.Handle(route.Method, route.Path, route.Handler)
		}
	}

	return router
}

func BuildApp(port int, modules ...DomainHandler) (*App, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	handler := NewRouter(modules...)
	app := NewApp(listener, handler)
	return app, nil
}
