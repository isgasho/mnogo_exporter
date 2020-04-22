package exporter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	exporterPrefix = "mongodb_"
)

var (
	// Rules to shrink metric names
	nodeToPDMetrics = map[string]string{
		"collStats.storageStats.indexDetails.":            "index_name",
		"globalLock.activeQueue.":                         "count_type",
		"globalLock.locks.":                               "lock_type",
		"serverStatus.asserts.":                           "assert_type",
		"serverStatus.connections.":                       "conn_type",
		"serverStatus.globalLock.currentQueue.":           "count_type",
		"serverStatus.metrics.commands.":                  "cmd_name",
		"serverStatus.metrics.cursor.open.":               "csr_type",
		"serverStatus.metrics.document.":                  "doc_op_type",
		"serverStatus.opLatencies.":                       "op_type",
		"serverStatus.opReadConcernCounters.":             "concern_type",
		"serverStatus.opcounters.":                        "legacy_op_type",
		"serverStatus.opcountersRepl.":                    "legacy_op_type",
		"serverStatus.transactions.commitTypes.":          "commit_type",
		"serverStatus.wiredTiger.concurrentTransactions.": "txn_rw_type",
		"serverStatus.wiredTiger.perf.":                   "perf_bucket",
		"systemMetrics.disks.":                            "device_name",
		/*Following needs to be tested once reportOpWriteConcernCountersInServerStatus*/
		/*  parameter is set*/
		/*"serverStatus.opWriteConcernCounters.":  "cmd_type",*/
		/*"globalLock.locks.<LOCK_TYPE>.acquireCount.":      "lock_mode",*/
		/*"globalLock.locks.<LOCK_TYPE>.acquireWaitCount.":  "lock_mode",*/
		/*"globalLock.locks.<LOCK_TYPE>.deadlockCount.":     "lock_mode",*/
		/*"globalLock.locks.<LOCK_TYPE>.timeAcquiringMicros.":  "lock_mode",*/
	}

	// Regular expressions used to make the metric name Prometheus-compatible
	specialCharsRe      = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	repeatedUnderscores = regexp.MustCompile(`__+`)
	dollarRe            = regexp.MustCompile(`\_$`)
)

func prometheusize(s string) string {
	prefixes := [][]string{
		{"serverStatus.wiredTiger.transaction", "ss_wt_txn"},
		{"serverStatus.wiredTiger", "ss_wt"},
		{"serverStatus", "ss"},
		{"replSetGetStatus", "rs"},
		{"systemMetrics", "sys"},
		{"local.oplog.rs.stats.wiredTiger", "oplog_stats_wt"},
		{"local.oplog.rs.stats", "oplog_stats"},
		{"collstats_storage.wiredTiger", "collstats_storage_wt"},
		{"collstats_storage.indexDetails", "collstats_storage_idx"},
		{"collStats.storageStats", "collstats_storage"},
		{"collStats.latencyStats", "collstats_latency"},
	}
	for _, pair := range prefixes {
		if strings.HasPrefix(s, pair[0]+".") {
			s = pair[1] + strings.TrimPrefix(s, pair[0])
			break
		}
	}

	s = specialCharsRe.ReplaceAllString(s, "_")
	s = dollarRe.ReplaceAllString(s, "")
	s = repeatedUnderscores.ReplaceAllString(s, "_")
	return exporterPrefix + s
}

func makeRawMetric(prefix, name string, value interface{}, labels map[string]string) (prometheus.Metric, error) {
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

	if labels == nil {
		labels = make(map[string]string)
	}

	fqName := prometheusize(prefix + name)
	if label, ok := nodeToPDMetrics[prefix]; ok {
		fqName = prometheusize(prefix)
		labels[label] = name
	}

	help := "TODO"
	typ := prometheus.UntypedValue

	ln := make([]string, 0)
	lv := make([]string, 0)
	for k, v := range labels {
		ln = append(ln, k)
		lv = append(lv, v)
	}

	d := prometheus.NewDesc(fqName, help, ln, nil)
	return prometheus.NewConstMetric(d, typ, f, lv...)
}
