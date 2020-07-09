package exporter

import (
	"context"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
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
	ctx := context.Background()

	for _, dbCollection := range d.collections {
		parts := strings.Split(dbCollection, ".")
		if len(parts) != 2 {
			continue
		}

		database := parts[0]
		collection := parts[1]

		aggregation := bson.D{
			{"$collStats", bson.M{"latencyStats": bson.E{"histograms", true}}}, //nolint
		}

		cursor, err := d.client.Database(database).Collection(collection).Aggregate(ctx, mongo.Pipeline{aggregation})
		if err != nil {
			logrus.Errorf("cannot get $collstats for collection %s.%s: %s", database, collection, err)
			continue
		}

		var stats []bson.M
		if err = cursor.All(ctx, &stats); err != nil {
			panic(err)
		}

		for _, m := range stats {
			for _, metric := range buildMetrics(m) {
				ch <- metric
			}
		}
	}
}
