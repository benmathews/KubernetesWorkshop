package visibilityworkshop

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger/options"
	_ "source.vivint.com/pl/messagetypes/nullable"
	_ "source.vivint.com/pl/messagetypes/objectid"
	_ "source.vivint.com/pl/messagetypes/time"
	_ "source.vivint.com/pl/protoc-gen-voter/v2/definition"
)

// Generate the GRPC server, REST server (GRPC gateway), and swagger documentation
// To validate your swagger json you can either install swagger UI locally or copy your json file into the editor at
// https://editor.swagger.io/
//go:generate bash -c "docker run -v $(pwd):/app -w /app tp-artifactory.vivint.com:5000/gobuild:1.14.6 ./generate.sh"
