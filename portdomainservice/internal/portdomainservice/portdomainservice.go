// Package portdomainservice used to work with remote databse exposing grpc interface
// TODO Change using direct field names to data binding approach
package portdomainservice

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/dehimb/shrimp/proto/port"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port       = ":50051"
	keyspace   = "ports"
	portsTable = "ports"
)

type PortService struct {
	grpcServer *grpc.Server
	dbSession  *gocql.Session
	logger     *logrus.Logger
}

// AddPort used to store port information to database
func (service *PortService) AddPort(ctx context.Context, port *pb.Port) (*pb.AddPortResponse, error) {
	err := service.dbSession.Query(fmt.Sprintf(`
		INSERT INTO %s.%s (ID, Name, Lat, Long, Country, Province, City, Timezone, Code, Alias, Region, Unlocs)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, keyspace, portsTable), port.Id, port.Name, port.Lat, port.Long, port.Country, port.Province, port.City, port.Timezone, port.Code, port.Alias, port.Region, port.Unlocs).Exec()

	if err != nil {
		return nil, err
	}
	return &pb.AddPortResponse{}, nil
}

// GetPort used to obtain information about port from database
func (service *PortService) GetPort(ctx context.Context, portID *pb.PortID) (*pb.GetPortResponse, error) {
	service.logger.Info("Req port with id ", portID.GetId())
	var (
		ID       string
		Name     string
		Lat      float32
		Long     float32
		Country  string
		Province string
		City     string
		Timezone string
		Code     string
		Alias    []string
		Region   []string
		Unlocs   []string
	)
	err := service.dbSession.Query(fmt.Sprintf(`
		SELECT ID, Name, Lat, Long, Country, Province, City, Timezone, Code, Alias, Region, Unlocs FROM %s.%s WHERE ID = ?
	`, keyspace, portsTable), portID.GetId()).Scan(&ID, &Name, &Lat, &Long, &Country, &Province, &City, &Timezone, &Code, &Alias, &Region, &Unlocs)
	if err != nil {
		return nil, err
	}
	return &pb.GetPortResponse{
		Port: &pb.Port{
			Id:       ID,
			Name:     Name,
			Lat:      Lat,
			Long:     Long,
			Country:  Country,
			Province: Province,
			City:     City,
			Timezone: Timezone,
			Code:     Code,
			Alias:    Alias,
			Region:   Region,
			Unlocs:   Unlocs,
		},
	}, nil
}

// Establish connection to database and start grpc server
func StartPortDomainService(ctx context.Context, logger *logrus.Logger) {
	service := &PortService{logger: logger}
	service.setupDatabase()
	// Set-up our gRPC server.
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service.grpcServer = grpc.NewServer()

	// Register our service with the gRPC server
	pb.RegisterPortDomainServiceServer(service.grpcServer, service)

	// Register reflection service on gRPC server.
	reflection.Register(service.grpcServer)

	go func() {
		if err := service.grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	log.Println("Started grpc on port:", port)

	waitForShutdown(ctx, service, logger)
}

// Initialize database
func (service *PortService) setupDatabase() {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_CLUSTER_URL"))
	// TODO move credentials to secrets store
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	var err error
	service.dbSession, err = cluster.CreateSession()
	if err != nil {
		service.logger.Fatal(err)
	}
	// Create keyspace if not exists
	err = service.dbSession.Query(fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class': 'SimpleStrategy','replication_factor': 1};", keyspace), nil).Exec()
	if err != nil {
		service.logger.Fatal(err)
	}
	// Create table if not exists
	err = service.dbSession.Query(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.%s (
		ID text,
		Name text,
		Lat float,
		Long float,
		Country text,
		Province text,
		City text,
		Timezone text,
		Code text,
		Alias list<text>,
		Region list<text>,
		Unlocs list<text>,
		PRIMARY KEY (ID))
	`, keyspace, portsTable), nil).Exec()
	if err != nil {
		service.logger.Fatal(err)
	}
	service.logger.Info("DB connection established")
}

// Wait for contex done signal for graceful shutdow grpc server and db session
func waitForShutdown(ctx context.Context, service *PortService, logger *logrus.Logger) {
	<-ctx.Done()
	service.grpcServer.GracefulStop()
	logger.Infoln("grpc service stopped")
	service.dbSession.Close()
	logger.Infoln("db session closed")
}
