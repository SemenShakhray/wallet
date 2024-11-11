package app

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wallet/internal/config"
	"wallet/internal/server/route"
	"wallet/internal/server/service"
)

type application struct {
	Config  config.Config
	Server  *http.Server
	Service service.Servicer
	Sigint  chan os.Signal
}

func NewApp() (application, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed reading config file: %v", err)
		return application{}, err
	}

	service, err := service.NewService(cfg)
	if err != nil {
		log.Printf("failed init service: %v", err)
		return application{}, err
	}

	route, err := route.NewRouter(service)
	if err != nil {
		log.Printf("failed init route: %v", err)
		return application{}, err
	}

	server := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: route,
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	return application{Config: cfg, Server: server, Service: service, Sigint: sigint}, nil
}
