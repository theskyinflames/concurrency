package service

import (
	csvrecord "github.com/theskyinflames/quiz/csvreader/pkg/domain"
)

type (
	CRMClient interface {
		ProcessRecord(record *csvrecord.Record) error
	}

	CRMIntegrator struct {
		crmClient CRMClient
	}
)

func NewCRMIntegrator(crmClient CRMClient) *CRMIntegrator {
	return &CRMIntegrator{
		crmClient: crmClient,
	}
}

func (c *CRMIntegrator) ProcessRecord(record *csvrecord.Record) error {
	return c.crmClient.ProcessRecord(record)
}
