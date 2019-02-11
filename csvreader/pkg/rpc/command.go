package rpc

import (
	"bytes"
	"encoding/gob"
)

// Item serializing
func ItemToGob(item interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Item deserializing
func FromGobToItem(b []byte, item interface{}) (interface{}, error) {
	buff := bytes.Buffer{}
	buff.Write(b)
	d := gob.NewDecoder(&buff)
	err := d.Decode(item) // item must be a pointer
	if err != nil {
		return nil, err
	}
	return item, nil
}
