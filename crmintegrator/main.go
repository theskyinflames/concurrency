package main

import (
	"github.com/theskyinflames/quiz/crmintegrator/command"
	"github.com/theskyinflames/quiz/crmintegrator/config"
	"github.com/theskyinflames/quiz/crmintegrator/crm"
	"github.com/theskyinflames/quiz/crmintegrator/rpc"
	"github.com/theskyinflames/quiz/crmintegrator/service"
)

func main() {

	cfg := &config.Config{}
	cfg.Load()

	crmclient := crm.NewCRMClient()
	crmintegrator := service.NewCRMIntegrator(crmclient)
	processRecordCommand := command.NewProcessRecordCommand(crmintegrator)
	server := rpc.NewRPCServer(cfg, processRecordCommand)
	server.Start()
}
