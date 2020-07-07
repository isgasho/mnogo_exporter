package exporter

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	replicationNotEnabled = 76
)

type replSetGetStatusCollector struct {
	ctx    context.Context
	client *mongo.Client
}

func (d *replSetGetStatusCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(d, ch)
}

func (d *replSetGetStatusCollector) Collect(ch chan<- prometheus.Metric) {
	cmd := bson.D{{Key: "replSetGetStatus", Value: "1"}}
	res := d.client.Database("admin").RunCommand(d.ctx, cmd)
	var m bson.M
	if err := res.Decode(&m); err != nil {
		if e, ok := err.(mongo.CommandError); ok && e.Code == replicationNotEnabled {
			return
		}
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
		return
	}

	for _, metric := range buildMetrics(m) {
		ch <- metric
	}
}
