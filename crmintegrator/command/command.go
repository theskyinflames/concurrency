package command

import (
	csvrecord "github.com/theskyinflames/quiz/csvreader/pkg/domain"
	"github.com/theskyinflames/quiz/csvreader/pkg/rpc"
)

type (
	CRMIntegrator interface {
		ProcessRecord(record *csvrecord.Record) error
	}

	ProcessRecordCommand struct {
		crmIntegrator CRMIntegrator
	}
)

func NewProcessRecordCommand(crmIntegrator CRMIntegrator) *ProcessRecordCommand {
	return &ProcessRecordCommand{
		crmIntegrator: crmIntegrator,
	}
}

func (p *ProcessRecordCommand) Process(b []byte, retry *bool) error {

	item, err := rpc.FromGobToItem(b, &csvrecord.Record{})
	if err != nil {
		return err
	}

	err = p.crmIntegrator.ProcessRecord(item.(*csvrecord.Record))
	if err != nil {
		*retry = true
		return err
	}
	*retry = false
	return nil
}
