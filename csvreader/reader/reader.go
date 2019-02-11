package reader

import (
	"bufio"
	"context"
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/theskyinflames/quiz/csvreader/pkg/domain"
)

type (
	CSVReaderFunc func(ctx context.Context, filePath string, recordsBlockSize int64, ch chan []domain.Record)
)

var (
	// This function opens the csv file, and reads it packaging its recors in blocks.
	// Each block is sent by the output channel to be sent to the CRM Integrator
	DefaultCSVReaderFunc = func(ctx context.Context, filePath string, recordsBlockSize int64, ch chan []domain.Record) {

		csvFile, _ := os.Open(filePath)
		reader := csv.NewReader(bufio.NewReader(csvFile))

		var (
			block []domain.Record
			err   error
			line  []string
		)
		z := recordsBlockSize

		for err == nil {
			select {
			case <-ctx.Done():
				break
			default:
				line, err = reader.Read()
				if err != nil {
					switch err {
					case io.EOF:
						log.Println("the end of the file has been achieved")
					default:
						log.Printf("ERROR: some went wrong when reading the csv file: %s", err.Error())
					}
					break
				}
				if z == recordsBlockSize {
					if block != nil { // First iteration
						ch <- block
					}
					z = 0
					block = make([]domain.Record, recordsBlockSize)
				}
				block[z] = domain.Record{
					ID:        line[0],
					FirstName: line[1],
					LastName:  line[2],
					Email:     line[3],
					Phone:     line[4],
				}
				z++
			}
		}
		close(ch)
		return
	}
)
