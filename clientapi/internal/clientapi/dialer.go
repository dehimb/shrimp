// Dialer used for establish grpc connection and
// add ability to mock this part for tests

package clientapi

import (
	"context"
	"os"

	pb "github.com/dehimb/shrimp/proto/port"
	"google.golang.org/grpc"
)

type Dialer interface {
	Dial(ctx context.Context) error
	GetClient() pb.PortDomainServiceClient
}

type GRPCDialer struct {
	grpcConnection       *grpc.ClientConn
	portDomainGRPCClient pb.PortDomainServiceClient
}

func (d *GRPCDialer) Dial(ctx context.Context) error {
	// Set up a connection to the server.
	var err error
	d.grpcConnection, err = grpc.Dial(os.Getenv("PORT_DOMAIN_SERVICE_GRPC_URL"), grpc.WithInsecure())
	if err != nil {
		return err
	}
	// register grpc client
	d.portDomainGRPCClient = pb.NewPortDomainServiceClient(d.grpcConnection)
	// listen for context done and close connection
	go func() {
		<-ctx.Done()
		d.grpcConnection.Close()
	}()
	return nil
}

func (d *GRPCDialer) GetClient() pb.PortDomainServiceClient {
	return d.portDomainGRPCClient
}
