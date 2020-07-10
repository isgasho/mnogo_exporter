package exporter

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type serverStatusCollector struct {
	client *mongo.Client
}

func (d *serverStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(d, ch)
}

func (d *serverStatusCollector) Collect(ch chan<- prometheus.Metric) {
	cmd := bson.D{{Key: "serverStatus", Value: "1"}}
	res := d.client.Database("admin").RunCommand(context.Background(), cmd)

	var m bson.M
	if err := res.Decode(&m); err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}

	for _, metric := range buildMetrics(m) {
		ch <- metric
	}
}

var _ prometheus.Collector = (*serverStatusCollector)(nil)
