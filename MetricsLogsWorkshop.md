---
marp: false
inlineSVG: true
header: 'Logs and Metrics'
footer: 'Ben Mathews - February 2021'
backgroundColor: #B9C9A9
---

# Metrics and Logs

## A developers introduction to Prometheus and Elasticsearch

https://source.vivint.com/users/ben.mathews/repos/logsmetricsalertsworkshop

---

# Prerequisites

- [Install and configure](https://confluence.vivint.com/display/PE/Platform+As+A+Service#PlatformAsAService-Kubectl) Kubectl for the TP Kubernetes cluster
- Install [Stern](https://github.com/wercker/stern)
- Install [GHZ](https://ghz.sh/)
- Install [gRPCman](https://confluence.vivint.com/display/PE/gRPCman+--+gRPC+Testing+Tool)

---

# Workshop outline - Metrics

- Best for
- Prometheus
  - Data flow
  - Conventions
  - Adding metrics
  - Querying

---  

# Workshop outline - Logging

- Best for
- Logging flow
- Stern
- Elasticsearch
  - Querying
  - Aggregation
  - Dashboards

---

# Workshop - Assignment

Using a sample app derived from pl/grpc

- Add metrics to measure XXXXXX
- Deploy it to your namespace
- Create a PodMonitor to scrape it
- View the logs in Elasticsearch
- Create an Elasticsearch aggregation that shows log volume
- Measure the rate of queries in Prometheus
- Measure the error rate in Prometheus
- Measure the latency in Prometheus
- Create a Grafana dashboard that shows the above metrics
- Create a alert and trigger it

---

# Metrics

- Number at a point in time
- Composed of
  - Identifier (Name/Labels)
  - Timestamp
  - Number
- Good for trends and overview
- Bad for detail

---
# Logs

- A specific event
- Composed of
  - Timestamp
  - Text (often JSON)
- Good for debugging
- Bad for big picture

---
# Tracing

Coming attraction

---

# Prometheus architecture


![](prometheusarchitecture.png)

---

# Prometheus - Demonstration

- Metrics coming from pod
- Configuration
- Query

<!-- 
In production
k port-forward deployment/sidequeue-nebo 2112:2112
xdg-open http://localhost:2112/metrics
topk(5,sort_desc(sum(rate(grpc_server_msg_received_total[10m])) by (app)))
grpc_method,grpc_service
-->

---

# Prometheus - Scraping configuration - Traditional

Annotations of a pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    prometheus.io/port: "2112"
    prometheus.io/scrape: "true"
```

<!-- How does this all get plumbed together? -->

---

# Prometheus - Scraping configuration - Prometheus Operator

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: example-app
spec:
  selector:
    app: example-app
  ports:
  - name: prometheus
    port: 8080
```

---

# Prometheus - Scraping configuration - Prometheus Operator

[ServiceMonitor](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md) object

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-app
spec:
  selector:
    matchLabels:
      app: example-app
  endpoints:
  - port: prometheus
```
---

# Prometheus metric types

- Counter - monotonically increasing
- Gauge - single numerical value that can arbitrarily go up and down
- Histogram - samples observations and counts them in configurable buckets
- Summary - calculates configurable quantiles over a sliding time window.

<!-- 
https://prometheus.io/docs/concepts/metric_types/

Differences
https://prometheus.io/docs/practices/histograms/#quantiles
-->

---

# Prometheus functions

- min/max
- deriv
- histogram_quantile
- rate/irate
- predict_linear
- sort/sort_desc
- \<aggregation\>_over_time
- lots of math functions

<!-- 
https://prometheus.io/docs/prometheus/latest/querying/functions/
-->

---

# Existing metrics

## Demo

<!-- 
curl localhost:2112/metrics
grpc_server_handling_seconds_count{grpc_method="HelloWorld",grpc_service="visibilityworkshop.VisibilityWorkshop",grpc_type="unary"} 1
-->

---
# Adding custom metrics

https://pkg.go.dev/github.com/prometheus/client_golang/prometheus


<!--
@@ -11,6 +11,9 @@ import (
        vtime "source.vivint.com/pl/messagetypes/time"
        "source.vivint.com/pl/mongo/v4"
        proto "source.vivint.com/pl/visibilityworkshop/generated"
+
+       "github.com/prometheus/client_golang/prometheus"
+       "github.com/prometheus/client_golang/prometheus/promauto"
 )
 
 const (
@@ -25,11 +28,21 @@ func NewServer() *server {
        return &server{}
 }
 
+var (
+       HelloWorldCounter = promauto.NewCounter(prometheus.CounterOpts{
+               Namespace: "viv_visibilityworkshop",
+               Name:      "hello_world_counter",
+               Help:      "Number of times Hello World has been called",
+       })
+)
+
 func (s *server) HelloWorld(ctx context.Context, request *proto.HelloWorldRequest) (*proto.HelloWorldResponse, error) {
        if request.GetName() == "" {
                return nil, vgrpc.MakeError(codes.InvalidArgument, "Name is a required parameter", nil)
        }
 
+       HelloWorldCounter.Inc()
+
        response := &proto.HelloWorldResponse{
-->


---

# Interact and generate some synthetic traffic

- gRPCman
- ghz

``` bash
go run ./cmd/visibilityworkshop/main.go -httpPort 8081
ghz localhost:9090 --call visibilityworkshop.VisibilityWorkshop.HelloWorld --protoset generated/visibilityworkshop.protoset --insecure -d '{"name":"test"}'
curl localhost:2112/metrics|grep -i hello
```

---
# Deploy to paas and add scraping


```bash
docker build -t tp-artifactory.vivint.com:5000/benmathews/visibilityworkshop:1 .
docker push tp-artifactory.vivint.com:5000/benmathews/visibilityworkshop:1 
cd charts/visibilityworkshop
helm dep up .
helm install visibilityworkshop  . --set visibilityworkshop.image.repository=benmathews/visibilityworkshop --set global.image.tag=1 
helm upgrade visibilityworkshop  . --set visibilityworkshop.image.repository=benmathews/visibilityworkshop --set global.image.tag=1 

```

---

# Generate some traffic

``` bash
kubectl port-forward deployment/visibilityworkshop 9090:9090

ghz localhost:9090 \
--call visibilityworkshop.VisibilityWorkshop.HelloWorld \
--protoset generated/visibilityworkshop.protoset \
--insecure -d '{"name":"test"}' \
--load-schedule="step" \
-n 10000 \
--load-start=5 \
--load-step=1 \
--load-step-duration=10s
```

---

# Prometheus functions

Create a dashboard showing metrics

<!-- 
rate(viv_visibilityworkshop_hello_world_counter[1m])
sum(rate(grpc_server_handled_total{namespace="ben", grpc_method="HelloWorld"}[1m]))
grpc_server_handling_seconds_sum{namespace="ben", grpc_method="HelloWorld"}/grpc_server_handling_seconds_count{namespace="ben", grpc_method="HelloWorld"}

https://grafana.platform.vivint.com/d/pz9G0DEGk/visibility-workshop?orgId=1&var-namespace=ben&from=now-1h&to=now
-->

---

![bg](questions.jpg)

---
