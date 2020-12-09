package visibilityworkshop

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	vgrpc "source.vivint.com/pl/grpc/v3/service"
	"source.vivint.com/pl/messagetypes/objectid"
	vtime "source.vivint.com/pl/messagetypes/time"
	"source.vivint.com/pl/mongo/v4"
	proto "source.vivint.com/pl/visibilityworkshop/generated"
)

const (
	ServiceName = "VisibilityWorkshop"
)

type server struct {
}

// NewServer creates a new VisibilityWorkshop server for << insert description here >>
func NewServer() *server {
	return &server{}
}

func (s *server) HelloWorld(ctx context.Context, request *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error) {
	if request.GetName() == "" {
		return nil, vgrpc.MakeError(codes.InvalidArgument, "Name is a required parameter", nil)
	}

	response := &proto.HelloWorldResponse{
		Id:        objectid.NewMgoDriverObjectId(mongo.NewObjectID()),
		Text:      fmt.Sprintf("Hello %s!!!", request.GetName()),
		Timestamp: vtime.NewCustomTime(time.Now()),
	}

	return response, nil
}
