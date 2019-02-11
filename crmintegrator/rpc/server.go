package rpc

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/theskyinflames/quiz/crmintegrator/config"
)

type (
	ProcessRecordCommand interface {
		Process(b []byte, retry *bool) error
	}

	RPCServer struct {
		cfg                  *config.Config
		processRecordCommand ProcessRecordCommand
	}
)

func NewRPCServer(cfg *config.Config, processRecordCommand ProcessRecordCommand) *RPCServer {
	return &RPCServer{
		cfg:                  cfg,
		processRecordCommand: processRecordCommand,
	}
}

func (r *RPCServer) Start() error {
	err := rpc.Register(r.processRecordCommand)
	if err != nil {
		return err
	}

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", r.cfg.RPCServerAddr)
	if err != nil {
		return err
	}
	log.Printf("Serving RPC server on port %s", r.cfg.RPCServerAddr)

	return http.Serve(listener, nil)
}
