package main

import (
	"fmt"
	"heydp/TgiInterface/internal/core"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel/trace"
)

type WebServer struct {
	router      *chi.Mux
	initialized bool
}

func NewWebServer() *WebServer {
	router := chi.NewRouter()

	return &WebServer{
		router: router,
	}
}

type ServerResources struct {
	tracer trace.Tracer
	ctrls  []Controllers
	itcrls []Controllers
}

type Controllers interface {
	MountRoutes(r chi.Router)
}

func (ws *WebServer) InitRouter(sr *ServerResources) {
	internalRoute := fmt.Sprintf("/internal")
	ws.router.Route(internalRoute, func(r chi.Router) {
		r.Use(middleware.Recoverer)
		r.Use(core.ContentTypeSetter)

		for _, c := range sr.itcrls {
			c.MountRoutes(r)
		}
	})

	externalRoute := fmt.Sprintf("/app")
	ws.router.Route(externalRoute, func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(core.ContentTypeSetter)
		r.Use(core.RequestIdHeader)

		for _, c := range sr.ctrls {
			c.MountRoutes(r)
		}
	})
}

func (ws *WebServer) Start(appconfig AppConfig) error {
	host := appconfig.Host
	port := appconfig.Port
	if ws.initialized == false {
		return fmt.Errorf("server not yet initialized, can't start")
	}

	serverAddr := fmt.Sprintf("%v:%v", host, port)
	fmt.Println("starting the webserver @ ", serverAddr)
	err := http.ListenAndServe(serverAddr, ws.router)
	if err != nil {
		errMsg := fmt.Sprintf("unable to start the server, err - %v", err.Error())
		fmt.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}
