// clientapi package provides http api to work with ports domain service
package clientapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	pb "github.com/dehimb/shrimp/proto/port"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type handler struct {
	router               *mux.Router
	logger               *logrus.Logger
	grpcConnection       *grpc.ClientConn
	portDomainGRPCClient pb.PortDomainServiceClient
}

func (h *handler) initRouter(ctx context.Context) error {
	// add logging middleware
	h.router.Use(h.logRequest)
	// add handler functions
	h.router.HandleFunc("/ping", h.ping).Methods("GET", "POST", "PUT")
	h.router.HandleFunc("/getPort/{portID}", h.getPort).Methods("GET")
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
		<-ctx.Done()
		h.grpcConnection.Close()
	}()

	return nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// Test handler to check service availability
func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Handler for retriving port information
func (h *handler) getPort(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := h.portDomainGRPCClient.GetPort(ctx, &pb.PortID{
		Id: mux.Vars(r)["portID"],
	})
	// TODO split grpc errors with no record found errors
	if err != nil {
		h.logger.Error("Error when requesting port: ", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	h.sendResponse(w, http.StatusOK, result)
}

// Handle all unsuported requests
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

func (h *handler) sendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Error when tryibg marshal response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(status)
	w.Write(json)
}
