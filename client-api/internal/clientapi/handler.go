package clientapi

import (
	"context"
	"net/http"
	"os"
	"time"

	pb "github.com/dehimb/shrimp/proto/port"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type handler struct {
	ctx                  context.Context
	router               *mux.Router
	logger               *logrus.Logger
	grpcConnection       *grpc.ClientConn
	portDomainGRPCClient pb.PortDomainServiceClient
}

func (h *handler) initRouter() error {
	// add logging middleware
	h.router.Use(h.logRequest)
	// add handler functions
	h.router.HandleFunc("/ping", h.ping).Methods("GET", "POST", "PUT")
	h.router.HandleFunc("/getPort/{portID}", h.getPort).Methods("GET", "POST", "PUT")
	// add default handler
	h.router.PathPrefix("/").HandlerFunc(h.defaultHandler)
	// Set up a connection to the server.
	var err error
	h.grpcConnection, err = grpc.Dial(os.Getenv("PORT_DOMAIN_SERVICE_GRPC_URL"), grpc.WithInsecure())
	if err != nil {
		return err
	}
	// register grpc client
	h.portDomainGRPCClient = pb.NewPortDomainServiceClient(h.grpcConnection)
	// listen for context done and close connection
	go func() {
		<-h.ctx.Done()
		h.grpcConnection.Close()
	}()

	return nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *handler) getPort(w http.ResponseWriter, r *http.Request) {
	result, err := h.portDomainGRPCClient.GetPort(h.ctx, &pb.PortID{})
	if err != nil {
		h.logger.Error("Error when requesting port: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.logger.Println(result)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// Used for loggin request method, url and execution time.
// Log only when log level set to logrus.InfoLevel or higher.
func (h *handler) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.logger.Level >= logrus.InfoLevel {
			start := time.Now()
			h.logger.Infof("-> %s %s", r.Method, r.URL)
			next.ServeHTTP(w, r)
			h.logger.Infof("<-  %s %s %s", time.Since(start), r.Method, r.URL)
			return
		}
		next.ServeHTTP(w, r)
	})
}
