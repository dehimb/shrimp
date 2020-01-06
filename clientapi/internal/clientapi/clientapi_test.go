package clientapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmizerany/assert"
	pb "github.com/dehimb/shrimp/proto/port"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	validPortID   = "VALID_PORT_ID"
	invalidPortID = "INVALID_PORT_ID"
)
var testHandler *handler

type mockGRPCClient struct{}

func (c *mockGRPCClient) GetPort(ctx context.Context, in *pb.PortID, opts ...grpc.CallOption) (*pb.GetPortResponse, error) {
	if in.GetId() == validPortID {
		return &pb.GetPortResponse{}, nil
	}
	return nil, errors.New("Valid test error")
}

func (c *mockGRPCClient) AddPort(ctx context.Context, in *pb.Port, opts ...grpc.CallOption) (*pb.AddPortResponse, error) {
	return nil, nil
}

type mockGRPCDialer struct {
}

func (d *mockGRPCDialer) Dial(ctx context.Context) error {
	return nil
}

func (d *mockGRPCDialer) GetClient() pb.PortDomainServiceClient {
	return &mockGRPCClient{}
}

func init() {
	testHandler = &handler{
		router: mux.NewRouter(),
		logger: logrus.New(),
		dialer: &mockGRPCDialer{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testHandler.initRouter(ctx)
}

func TestGetPort(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{
			name:         "Valid request",
			path:         fmt.Sprintf("/getPort/%s", validPortID),
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid request",
			path:         fmt.Sprintf("/getPort/%s", invalidPortID),
			expectedCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", testCase.path, nil)
			testHandler.ServeHTTP(rec, req)
			assert.Equal(t, testCase.expectedCode, rec.Code)
		})
	}

}
