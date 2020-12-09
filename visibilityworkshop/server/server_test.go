package visibilityworkshop

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/assert"

	proto "source.vivint.com/pl/visibilityworkshop/generated"

	"github.com/stretchr/testify/suite"
)

type TestVisibilityWorkshopSuite struct {
	suite.Suite
	VisibilityWorkshopServer *server
}

func (s *TestVisibilityWorkshopSuite) SetupTest() {
	s.VisibilityWorkshopServer = &server{}
}

func (s *TestVisibilityWorkshopSuite) TestHelloWorldHappy() {
	name := "Todd"

	request := &proto.HelloWorldRequest{
		Name: name,
	}

	response, err := s.VisibilityWorkshopServer.HelloWorld(context.Background(), request)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), response)
	assert.NotNil(s.T(), response.GetText())
	assert.Equal(s.T(), fmt.Sprintf("Hello %s!!!", name), response.GetText())
	assert.NotNil(s.T(), response.Id)
	assert.NotNil(s.T(), response.Id.ObjectId())
	assert.NotNil(s.T(), response.Timestamp)

	t, err := response.Timestamp.Time()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), t)
}

func (s *TestVisibilityWorkshopSuite) TestHelloWorldMissingName() {
	request := &proto.HelloWorldRequest{}

	response, err := s.VisibilityWorkshopServer.HelloWorld(context.Background(), request)

	assert.Nil(s.T(), response)
	assert.NotNil(s.T(), err)

	st, ok := status.FromError(err)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), codes.InvalidArgument, st.Code())
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestVisibilityWorkshopSuite))
}
