package exporter

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type collstatsCollector struct {
	ctx         context.Context
	client      *mongo.Client
	collections []string
}

func (d *collstatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(d, ch)
}

func (d *collstatsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, dbCollection := range d.collections {
		parts := strings.Split(dbCollection, ".")
		if len(parts) != 2 {
			continue
		}

		cmd := bson.D{{Key: "collStats", Value: parts[1]}}
		res := d.client.Database(parts[0]).RunCommand(d.ctx, cmd)

		var m bson.M
		if err := res.Decode(&m); err != nil {
			ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
			continue
		}

		for _, metric := range buildMetrics(m) {
			ch <- metric
		}
	}
}
