package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"

	"github.com/theskyinflames/quiz/csvreader/pkg/domain"
	serializer "github.com/theskyinflames/quiz/csvreader/pkg/rpc"
	"github.com/theskyinflames/quiz/csvreader/reader"

	"github.com/theskyinflames/quiz/csvreader/config"
)

const remoteRPCMethodName = "ProcessRecordCommand.Process"

type (
	Repository interface {
		InsertBlock(block []domain.Record) error
	}

	Service struct {
		ctx               context.Context
		cfg               *config.Config
		client            *rpc.Client
		repository        Repository
		concurrentSenders int32
	}
)

func NewService(ctx context.Context, cfg *config.Config, repository Repository) *Service {
	return &Service{
		ctx:        ctx,
		cfg:        cfg,
		repository: repository,
	}
}

func (r *Service) Start(reader reader.CSVReaderFunc) (err error) {

	// Start the RPC client
	r.client, err = rpc.DialHTTP("tcp", r.cfg.RPCCRMServerAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			log.Printf("concurrent senders %d \n", r.concurrentSenders)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Start the reader
	wg := sync.WaitGroup{}

	// If the CRM Integrator is broken down, as a maximum,
	// the taken memory space in the server will be
	// (cfg.MaxConcurrentSenders * cfg.r.cfg.CSVRecordsBlock * [record size in bytes])
	ch := make(chan []domain.Record, int(r.cfg.CSVRecordsBlock/int64(r.cfg.MaxConcurrentSenders))+1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader(r.ctx, r.cfg.FilePath, r.cfg.CSVRecordsBlock, ch)
		log.Println("csv reader function finishes")
		return
	}()

	// Start the senders. Record sending is a synchronized operation. It waits for
	// the CRM integrator response for each record. This response may be directly an error,
	// or a retry flag. As well as a error is returned, or the retry flag returns true, the
	// record will be resent. Until it works fine, or the maximum number of process attemps
	// is achieved. If this occurs, the record will be logged and discarded
	wg.Add(1)
	go func() {
		defer wg.Done()
		keepOn := true
		for keepOn {
			select {
			case <-r.ctx.Done():
				keepOn = false
				break
			case block, ok := <-ch:
				if !ok {
					keepOn = false
					break
				}
				for {
					// It may be this check for concurrent sender must be
					// synchonized. But this would have an important cost in
					// performance. To put it here with a previous tunning
					// work will be over engineering. So, this limit may be not
					// exactly
					if r.concurrentSenders < int32(r.cfg.MaxConcurrentSenders) {
						wg.Add(1)
						go r.processBlock(block, &wg)
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
		log.Println("all senders has finished")
		return
	}()

	log.Println("wating ....")
	wg.Wait()
	log.Println("all records processed")
	return nil
}

func (r *Service) processBlock(block []domain.Record, wg *sync.WaitGroup) (err error) {
	defer wg.Done()
	err = r.sendRecordToCrmIntegrator(block) // Send the records's block to CRMIntegrator
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = r.repository.InsertBlock(block) // Save the sent records to PostgreSQL db
	if err != nil {
		log.Printf("ERROR: something went wrong when trying to save records to db: %s\n", err.Error())
	}
	return
}

func (r *Service) sendRecordToCrmIntegrator(block []domain.Record) error {

	// Update concurrent senders counter
	atomic.AddInt32(&r.concurrentSenders, 1)
	defer func() {
		atomic.AddInt32(&r.concurrentSenders, -1)
	}()

	// Tries to send the records to be processed by the CRM Integrator
	sent := false
	for _, record := range block {
		b, err := serializer.ItemToGob(record)
		if err != nil {
			log.Printf("ERROR: the record's %s deserializing has failed by: %s", record.ID, err.Error())
			break
		}

		pendingAttemps := r.cfg.MaxAttempsPerRecord
		for !sent {
			pendingAttemps--
			var retry bool
			err = r.client.Call(remoteRPCMethodName, b, &retry)
			if err != nil {
				log.Printf("ERROR: some when wrong when trying to process the record %s: %s\n", record.ID, err.Error())
			} else {
				sent = true
			}

			if pendingAttemps > 0 && (err != nil || retry) {
				continue
			}
			break
		}
		if !sent {
			return errors.New(fmt.Sprintf("ERROR: the records block started with id %s has not been sent", block[0].ID))
		}
		log.Printf("record %s sent \n", record.ID)
	}
	return nil
}
