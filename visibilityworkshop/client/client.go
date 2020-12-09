package visibilityworkshop

import (
	"errors"
	"time"

	"source.vivint.com/pl/grpc/v3/service"

	"golang.org/x/net/context"
	"source.vivint.com/pl/flag"
	proto "source.vivint.com/pl/visibilityworkshop/generated"
)

var (
	host           = flag.String("visibilityWorkshopServiceHost", "localhost", "Host for VisibilityWorkshop gRPC server", "VISIBILITY_WORKSHOP_SERVICE_HOST")
	port           = flag.Int("visibilityWorkshopServicePort", 9090, "Port for VisibilityWorkshop gRPC server", "VISIBILITY_WORKSHOP_SERVICE_PORT")
	backoffSeconds = flag.Int("visibilityWorkshopServiceMaxBackoff", 60, "Max elapsed time in seconds for retrying calls before failing", "VISIBILITY_WORKSHOP_SERVICE_MAX_BACKOFF")
)

// Host is a convenient accessor for retrieving the default hostname that the VisibilityWorkshop service runs on
func Host() *string {
	return host
}

// Port is a convenient accessor for retrieving the default port that the VisibilityWorkshop service listens for gRPC on
func Port() *int {
	return port
}

// MaxBackoffTime is a convenient accessor for retrieving the default maximum backoff time used when connecting to the VisibilityWorkshop service over gRPC
func MaxBackoffTime() *time.Duration {
	backoffDuration := time.Duration(*backoffSeconds) * time.Second
	return &backoffDuration
}

// VisibilityWorkshopClienter is the client interface for the VisibilityWorkshop gRPC service
type VisibilityWorkshopClienter interface {
	// HelloWorld is a simplified client call to the service that takes the string name and returns the string response, hiding
	// the underlying grpc request/response details from the client - this will vary per use case in real scenarios
	HelloWorld(name string) (string, error)
}

// Client implements VisibilityWorkshopClienter and leverages the retryable gRPC client implementation for backoff retries
type Client struct {
	retryClient *service.RetryClient
	client      proto.VisibilityWorkshopClient
}

// NewClient creates a new client
func NewClient(host string, port int, contexter func() context.Context, maxBackoffTime time.Duration) *Client {
	retryClient := service.NewRetryClient(host, port, contexter, maxBackoffTime)

	return &Client{
		retryClient: retryClient,
		client:      proto.NewVisibilityWorkshopClient(retryClient.Conn),
	}
}

// Get retrieves the specified account system object from the gRPC service using it's id
func (w *Client) HelloWorld(name string) (string, error) {

	request := proto.HelloWorldRequest{
		Name: name,
	}

	result, err := w.retryClient.Call(func(ctx context.Context) (interface{}, error) {
		return w.client.HelloWorld(ctx, &request)
	})

	if err != nil {
		return "", err
	}

	response, ok := result.(*proto.HelloWorldResponse)

	if !ok {
		return "", errors.New("Unexpected response type!")
	}

	return response.GetText(), nil
}
