// Package portsloader is used for reading ports list from file and store this data to ports domain service
package portsloader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dehimb/shrimp/proto/port"
	pb "github.com/dehimb/shrimp/proto/port"
	"google.golang.org/grpc"
)

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

// Load method read data from given file and pass it to ports domain service
func Load(ctx context.Context, file *os.File) (*Results, error) {
	// close file when all done
	defer file.Close()
	// Set up a connection to port domain service.
	conn, err := grpc.Dial(os.Getenv("PORT_DOMAIN_SERVICE_GRPC_URL"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	// close connection when all done
	defer conn.Close()
	// register grpc client
	client := pb.NewPortDomainServiceClient(conn)

	// use decoder to read file by parts for minimizing memory consumption
	decoder := json.NewDecoder(file)
	// reading first "{"
	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}
	results := &Results{}
	// read file intil the end
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
		// call port save via rpc
		_, err = client.AddPort(ctx, port)
		if err != nil {
			return nil, err
		}
		results.Count++
	}
	return results, nil
}
