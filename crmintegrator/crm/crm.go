package crm

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	csvrecord "github.com/theskyinflames/quiz/csvreader/pkg/domain"
)

type (
	CRMClient struct {
		randomBool RandomBool
	}
)

func NewCRMClient() *CRMClient {
	return &CRMClient{randomBool: NewRandomBool()}
}

func (c *CRMClient) ProcessRecord(record *csvrecord.Record) error {

	b, _ := json.Marshal(record)
	log.Printf("received record %s", string(b))

	// Simulate CRM API failing randomly
	if c.randomBool.Bool() {
		// API fails
		return errors.New("CRM API fail")
	}
	time.Sleep(10 * time.Millisecond) // Process time simulation

	return nil
}
