// Package portsloader is used for reading ports list from file and store this data to ports domain service
package portsloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dehimb/shrimp/clientapi/internal/clientapi"
	"github.com/dehimb/shrimp/proto/port"
)

type Loader struct {
	dialer clientapi.Dialer
}

type Results struct {
	Count uint
}

// Representation of port item in json format
type Item struct {
	Name        string
	City        string
	Country     string
	Province    string
	Timezone    string
	Code        string
	Alias       []string
	Regions     []string
	Unlocs      []string
	Coordinates []float32
}

func New(ctx context.Context) (*Loader, error) {
	loader := &Loader{
		dialer: &clientapi.GRPCDialer{},
	}
	err := loader.dialer.Dial(ctx)
	if err != nil {
		return nil, err
	}
	return loader, nil
}

// Load method read data from given file and pass it to ports domain service
func (l *Loader) Load(ctx context.Context, file io.Reader) (*Results, error) {
	// use decoder to read file by parts for minimizing memory consumption
	decoder := json.NewDecoder(file)
	// reading first "{"
	_, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	results := &Results{}
	// read file until the end
	for decoder.More() {
		// read port id
		key, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		// read port data
		var item Item
		err = decoder.Decode(&item)
		if err != nil {
			return nil, err
		}
		// check coordinates
		var (
			lat  float32
			long float32
		)
		if len(item.Coordinates) == 2 {
			lat = item.Coordinates[0]
			long = item.Coordinates[1]
		}
		port := &port.Port{
			Id:       fmt.Sprintf("%v", key),
			Name:     item.Name,
			Lat:      lat,
			Long:     long,
			Country:  item.Country,
			Province: item.Province,
			City:     item.City,
			Timezone: item.Timezone,
			Code:     item.Code,
			Alias:    item.Alias,
			Region:   item.Regions,
			Unlocs:   item.Unlocs,
		}
		// call port save via grpc
		_, err = l.dialer.GetClient().AddPort(ctx, port)
		if err != nil {
			return nil, err
		}
		results.Count++
	}
	return results, nil
}
