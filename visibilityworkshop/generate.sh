#!/bin/bash

set -e # important so errors aren't ignored

path=$(pwd)
if [[ ${HOST_DIR} ]]; then
   path=${HOST_DIR}
fi
ROOT_DIR=${ROOT_DIR:-${path}}

MAIN_PB=generated/visibilityworkshop.pb.go

# Remove our main generated Go service proto file so that we can verify if Docker succeeded
if [ -f ${MAIN_PB} ]; then
    rm ${MAIN_PB}
fi

protoc.sh \
    --swagger_out=logtostderr=true:generated \
    --grpc-gateway_out=logtostderr=true,paths=source_relative:generated \
	--govalidators_out=paths=source_relative:generated \
    --gogo_out=plugins=grpc,paths=source_relative:generated \
    --include_imports \
    --descriptor_set_out=generated/visibilityworkshop.protoset \
    --voter_out=paths=source_relative:generated \
    --proto_path=${ROOT_DIR}/vendor \
    visibilityworkshop.proto

# If the Docker image failed to generate our main Go service proto then fail
if [ ! -f ${MAIN_PB} ]; then
    echo "Could not find ${MAIN_PB} after generating protos"
    exit 1
fi

goimports -w $(find . -type f -name '*.go' -not -path './vendor/*')

