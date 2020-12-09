
FROM tp-artifactory.vivint.com:5000/gobuild:1.15 as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod vendor
RUN ./generate.sh
RUN CGO_ENABLED=0 GOOS=linux go build -i  -a -o binary ./cmd/visibilityworkshop

FROM tp-artifactory.vivint.com:5001/alpine:3.12.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/binary visibilityworkshop
CMD ["./visibilityworkshop"]

