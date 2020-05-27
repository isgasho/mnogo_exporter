package exporter

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type serverStatusCollector struct {
	ctx    context.Context
	client *mongo.Client
}

func (d *serverStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(d, ch)
}

func (d *serverStatusCollector) Collect(ch chan<- prometheus.Metric) {
	cmd := bson.D{{Key: "serverStatus", Value: "1"}}
	res := d.client.Database("admin").RunCommand(d.ctx, cmd)
	var m bson.M
	if err := res.Decode(&m); err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}

	m, ok := m["data"].(bson.M)
	if !ok {
		err := fmt.Errorf("unexpected %T for data", m["data"])
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}

	for _, metric := range makeMetrics("", m, nil) {
		ch <- metric
	}
}
