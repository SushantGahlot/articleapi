package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	srv *http.Server
}

func GetServer(addr string, logger *log.Logger, router *httprouter.Router) (*Server, error) {
	if len(addr) == 0 {
		return nil, errors.New("Server address missing")
	}

	if router == nil {
		return nil, errors.New("Server handlers can not be nil")
	}

	srv := http.Server{
		Addr:     ":" + addr,
		ErrorLog: logger,
		Handler:  router,
	}

	return &Server{
		srv: &srv,
	}, nil
}

func (srv *Server) Start() error {
	return srv.srv.ListenAndServe()
}

func (srv *Server) Close() error {
	return srv.srv.Close()
}
