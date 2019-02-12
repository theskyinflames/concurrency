package rpc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	Item struct {
		Name string
		Age  int
	}
)

func TestStorage_ItemToGob(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "given a storage, when a item is serialized, then all works fine",
			args:    args{Item{Name: "Bartolo", Age: 22}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ItemToGob(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ItemToGob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.True(t, len(got) > 0)
		})
	}
}

func TestStorage_FromGobToItem(t *testing.T) {

	item := Item{Name: "Bartolo", Age: 22}

	bin, err := ItemToGob(item)
	if err != nil {
		t.Fail()
	}

	type args struct {
		b    []byte
		item *Item
	}
	tests := []struct {
		name    string
		args    args
		want    Item
		wantErr bool
	}{
		{
			name:    "given a serialized item, when it's deserialized, then all works fine",
			args:    args{b: bin, item: &Item{}},
			want:    item,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromGobToItem(tt.args.b, tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.FromGobToItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*(got.(*Item)), tt.want) {
				t.Errorf("Storage.FromGobToItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
