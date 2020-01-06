package portsloader

import (
	"context"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
	pb "github.com/dehimb/shrimp/proto/port"
	"google.golang.org/grpc"
)

type mockGRPCClient struct{}

func (c *mockGRPCClient) GetPort(ctx context.Context, in *pb.PortID, opts ...grpc.CallOption) (*pb.GetPortResponse, error) {
	return nil, nil
}

func (c *mockGRPCClient) AddPort(ctx context.Context, in *pb.Port, opts ...grpc.CallOption) (*pb.AddPortResponse, error) {
	return &pb.AddPortResponse{}, nil
}

type mockGRPCDialer struct {
}

func (d *mockGRPCDialer) Dial(ctx context.Context) error {
	return nil
}

func (d *mockGRPCDialer) GetClient() pb.PortDomainServiceClient {
	return &mockGRPCClient{}
}

var testLoader *Loader

func init() {
	testLoader = &Loader{
		dialer: &mockGRPCDialer{},
	}
	testLoader.dialer.Dial(context.Background())
}

func TestGetPort(t *testing.T) {
	// TODO add error type check
	testCases := []struct {
		name        string
		data        string
		expectError bool
	}{
		{
			name:        "Valid json",
			data:        "{}",
			expectError: false,
		},
		{
			name:        "Empty reader",
			expectError: true,
		},
		{
			name:        "Malformed json",
			data:        "Malformed json",
			expectError: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := testLoader.Load(context.Background(), strings.NewReader(testCase.data))
			assert.Equal(t, testCase.expectError, err != nil)
		})
	}
}
