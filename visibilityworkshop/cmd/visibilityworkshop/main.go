package main

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"source.vivint.com/pl/flag"

	vgrpc "source.vivint.com/pl/grpc/v3/service"
	"source.vivint.com/pl/log"
	proto "source.vivint.com/pl/visibilityworkshop/generated"
	visibilityworkshop "source.vivint.com/pl/visibilityworkshop/server"
)

var (
	logLevel                        = flag.LogLevelFlag()
	address                         = flag.AddressFlag()
	httpPort                        = flag.Int("httpPort", 8080, "HTTP Listening port", "HTTP_PORT")
	grpcPort                        = flag.Int("grpcPort", 9090, "gRPC Listening port", "GRPC_PORT")
	grpcReflection                  = flag.EnableGrpcReflection()
	prometheusPort                  = flag.PrometheusPort(2112)
	prometheusHandlingTimeHistogram = flag.EnablePrometheusHandlingTimeHistogram()
	maximumConnectionAge            = flag.MaxGrpcConnectionAge()
	maximumConnectionAgeGracePeriod = flag.MaxGrpcConnectionAgeGracePeriod()
)

const (
	defaultLogLevel = log.INFO
)

func main() {
	flag.Parse()

	log.SetupHandlerLogging(visibilityworkshop.ServiceName, *logLevel, defaultLogLevel)

	httpAddr := fmt.Sprintf("%v:%v", *address, *httpPort)
	grpcAddr := fmt.Sprintf("%v:%v", *address, *grpcPort)

	server := visibilityworkshop.NewServer()
	grpcServerConf := vgrpc.ServerConf{
		Name:        visibilityworkshop.ServiceName,
		GrpcAddress: grpcAddr,
		HTTPAddress: httpAddr,
		GrpcRegister: func(grpcServer *grpc.Server) {
			proto.RegisterVisibilityWorkshopServer(grpcServer, server)
		},
		HTTPFunc: proto.RegisterVisibilityWorkshopHandlerFromEndpoint,
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			vgrpc.ValidateInterceptor,
		},
		Options: []vgrpc.Option{
			vgrpc.GrpcReflection(*grpcReflection),
			vgrpc.PrometheusPort(int32(*prometheusPort)),
			vgrpc.PrometheusHandlingTimeHistogram(*prometheusHandlingTimeHistogram),
		},
		GrpcOptions: []grpc.ServerOption{
			grpc.KeepaliveParams(keepalive.ServerParameters{
				MaxConnectionAge:      *maximumConnectionAge,
				MaxConnectionAgeGrace: *maximumConnectionAgeGracePeriod,
			}),
		},
	}

	grpcServer, err := vgrpc.NewServer(grpcServerConf)
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(); err != nil {
		log.Error("Server closed unexpectedly", log.Fields{"err": err})
	}
}
