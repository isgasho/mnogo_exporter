package exporter

import (
	"fmt"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ssRe                 = regexp.MustCompile(`^serverStatus`)
	wtRe                 = regexp.MustCompile(`^ss\.wiredTiger`)
	txnRe                = regexp.MustCompile(`^ss_wt\.transaction`)
	rsRe                 = regexp.MustCompile(`^replSetGetStatus`)
	sysRe                = regexp.MustCompile(`^systemMetrics`)
	oplogStatsRe         = regexp.MustCompile(`^local\.oplog\.rs\.stats`)
	oplogStatsWtRe       = regexp.MustCompile(`^oplog_stats\.wiredTiger`)
	collstatsLatencyRe   = regexp.MustCompile(`^collStats\.latencyStats`)
	scollstatsStorageRe  = regexp.MustCompile(`^collStats\.storageStats`)
	collstatsStorageWtRe = regexp.MustCompile(`^collstats_storage\.wiredTiger`)
	colstatsStorageIdxRe = regexp.MustCompile(`^collstats_storage\.indexDetails`)
	specialCharsRe       = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	dollarRe             = regexp.MustCompile(`\_$`)
)

func prometheusize(s string) string {
	s = ssRe.ReplaceAllString(s, "ss")
	s = wtRe.ReplaceAllString(s, "ss_wt")
	s = txnRe.ReplaceAllString(s, "ss_wt_txn")
	s = rsRe.ReplaceAllString(s, "rs")
	s = sysRe.ReplaceAllString(s, "sys")
	s = oplogStatsRe.ReplaceAllString(s, "oplog_stats")
	s = oplogStatsWtRe.ReplaceAllString(s, "oplog_stats_wt")
	s = collstatsLatencyRe.ReplaceAllString(s, "collstats_latency")
	s = scollstatsStorageRe.ReplaceAllString(s, "collstats_storage")
	s = collstatsStorageWtRe.ReplaceAllString(s, "collstats_storage_wt")
	s = colstatsStorageIdxRe.ReplaceAllString(s, "collstats_storage_idx")
	s = dollarRe.ReplaceAllString(s, "")

	return s
}

func fqMetricName(s string) string {
	s = specialCharsRe.ReplaceAllString(s, "_")
	s = dollarRe.ReplaceAllString(s, "")
	return s
}

func makeRawMetric(name string, value interface{}, labels map[string]string) (prometheus.Metric, error) {
	var f float64
	switch v := value.(type) {
	case bool:
		if v {
			f = 1
		}
	case int32:
		f = float64(v)
	case int64:
		f = float64(v)
	case float64:
		f = v

	case primitive.DateTime:
		f = float64(v)
	case primitive.Timestamp:
		return nil, nil
	case primitive.ObjectID:
		return nil, nil
	case string:
		return nil, nil

	default:
		return nil, fmt.Errorf("makeRawMetric: unhandled type %T", v)
	}

	fqName := fqMetricName(name)
	help := "TODO"
	typ := prometheus.UntypedValue
	if len(labels) == 0 {
		labels = nil
	}

	ln := make([]string, 0)
	lv := make([]string, 0)
	for k, v := range labels {
		ln = append(ln, k)
		lv = append(lv, v)
	}

	d := prometheus.NewDesc(fqName, help, ln, nil)
	return prometheus.NewConstMetric(d, typ, f, lv...)
}
