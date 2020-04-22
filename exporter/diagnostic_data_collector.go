package exporter

import (
	"context"
	"fmt"

	// "github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type diagnosticDataCollector struct {
	ctx    context.Context
	client *mongo.Client
	// l      log.Logger
}

func (d *diagnosticDataCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(d, ch)
}

func (d *diagnosticDataCollector) Collect(ch chan<- prometheus.Metric) {
	// cmd := bson.D{{Key: "replSetGetStatus", Value: "1"}}
	cmd := bson.D{{Key: "getDiagnosticData", Value: "1"}}
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

	for _, metric := range d.makeMetrics("", m, nil) {
		ch <- metric
	}
}

func (d *diagnosticDataCollector) makeMetrics(prefix string, m bson.M, labels map[string]string) []prometheus.Metric {
	var res []prometheus.Metric
	//prefix = prometheusize(prefix)
	if prefix != "" {
		prefix += "."
	}

	for k, val := range m {
		switch v := val.(type) {
		case bson.M:
			res = append(res, d.makeMetrics(prefix+k, v, labels)...)
		case map[string]interface{}:
			res = append(res, d.makeMetrics(prefix+k, v, labels)...)
		case primitive.A:
			v = []interface{}(v)
			res = append(res, d.processSlice(prefix, k, v)...)
		case []interface{}:
			continue
			// res = append(res, d.processSlice(prefix, k, v)...)

		default:
			metric, err := makeRawMetric(prefix, k, v, labels)
			if err != nil {
				// TODO
				panic(err)
			}
			if metric != nil {
				res = append(res, metric)
			}
		}
	}

	return res
}

func (d *diagnosticDataCollector) processSlice(prefix, k string, v []interface{}) []prometheus.Metric {
	metrics := make([]prometheus.Metric, 0)
	labels := make(map[string]string)

	for _, item := range v {
		var s map[string]interface{}

		switch i := item.(type) {
		case map[string]interface{}:
			s = i
		case primitive.M:
			s = map[string]interface{}(i)
		}

		if name, ok := s["name"].(string); ok {
			labels["member_idx"] = name
		}

		metrics = append(metrics, d.makeMetrics(prefix+k, s, labels)...)
	}

	return metrics
}

// check interface
var _ prometheus.Collector = (*diagnosticDataCollector)(nil)
