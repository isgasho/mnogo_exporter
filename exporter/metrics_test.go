package exporter

import (
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test metric renaming and labeling.
func TestMetricName(t *testing.T) {
	tcs := []struct {
		prefix     string
		name       string
		wantMetric string
		wantLabel  string
	}{
		{
			prefix:     "serverStatus.metrics.commands.saslStart.",
			name:       "total",
			wantMetric: "mongodb_ss_metrics_commands_saslStart_total",
		},
		{
			prefix:     "serverStatus.metrics.commands._configsvrShardCollection.",
			name:       "failed",
			wantMetric: "mongodb_ss_metrics_commands_configsvrShardCollection_failed",
		},
		{
			prefix:     "serverStatus.wiredTiger.lock.",
			name:       "metadata lock acquisitions",
			wantMetric: "mongodb_ss_wt_lock_metadata_lock_acquisitions",
		},
		{
			prefix:     "serverStatus.wiredTiger.perf.",
			name:       "file system write latency histogram (bucket 5) - 500-999ms",
			wantMetric: "mongodb_ss_wt_perf",
			wantLabel:  "perf_bucket",
		},
		{
			prefix:     "serverStatus.wiredTiger.transaction.",
			name:       "rollback to stable updates removed from lookaside",
			wantMetric: "mongodb_ss_wt_txn_rollback_to_stable_updates_removed_from_lookaside",
		},
	}

	for _, tc := range tcs {
		metric, label := nameAndLabel(tc.prefix, tc.name)
		assert.Equal(t, tc.wantMetric, metric, tc.prefix+tc.name)
		assert.Equal(t, tc.wantLabel, label, tc.prefix+tc.name)
	}
}

func TestPrometeusize(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "serverStatus.wiredTiger.transaction.transaction checkpoint most recent time (msecs)",
			want: "mongodb_ss_wt_txn_transaction_checkpoint_most_recent_time_msecs",
		},
		{
			in:   "serverStatus.wiredTiger.thread-yield.page acquire time sleeping (usecs)",
			want: "mongodb_ss_wt_thread_yield_page_acquire_time_sleeping_usecs",
		},
		{
			in:   "serverStatus.opLatencies.reads.latency",
			want: "mongodb_ss_opLatencies_reads_latency",
		},
		{
			in:   "replSetGetStatus.optimes.lastCommittedOpTime.t",
			want: "mongodb_rs_optimes_lastCommittedOpTime_t",
		},
		{
			in:   "systemMetrics.memory.Active_kb",
			want: "mongodb_sys_memory_Active_kb",
		},
		{
			in:   "local.oplog.rs.stats.wiredTiger.block-manager.checkpoint size",
			want: "mongodb_oplog_stats_wt_block_manager_checkpoint_size",
		},
		{
			in:   "local.oplog.rs.stats.storageSize",
			want: "mongodb_oplog_stats_storageSize",
		},
		{
			in:   "collstats_storage.wiredTiger.xxx",
			want: "mongodb_collstats_storage_wt_xxx",
		},

		{
			in:   "collstats_storage.indexDetails.xxx",
			want: "mongodb_collstats_storage_idx_xxx",
		},
		{
			in:   "collStats.storageStats.xxx",
			want: "mongodb_collstats_storage_xxx",
		},
		{
			in:   "collStats.latencyStats.xxx",
			want: "mongodb_collstats_latency_xxx",
		},
	}

	for _, test := range tests {
		got := prometheusize(test.in)
		assert.Equal(t, test.want, got)
	}
}

// Test supported value types conversion.
func TestMakeRawMetric(t *testing.T) {
	prefix := "serverStatus.transactions."
	name := "retriedCommandsCount"
	testCases := []struct {
		value     interface{}
		wantVal   *float64
		wantError bool
	}{
		{value: true, wantVal: pointer.ToFloat64(1), wantError: false},
		{value: false, wantVal: pointer.ToFloat64(0), wantError: false},
		{value: int32(1), wantVal: pointer.ToFloat64(1), wantError: false},
		{value: int64(2), wantVal: pointer.ToFloat64(2), wantError: false},
		{value: float32(1.23), wantVal: pointer.ToFloat64(float64(float32(1.23))), wantError: false},
		{value: float64(1.23), wantVal: pointer.ToFloat64(1.23), wantError: false},
		{value: primitive.A{}, wantVal: nil, wantError: false},
		{value: primitive.Timestamp{}, wantVal: nil, wantError: false},
		{value: "zapp", wantVal: nil, wantError: false},
		{value: []byte{}, wantVal: nil, wantError: true},
		{value: time.Date(2020, 06, 15, 0, 0, 0, 0, time.UTC), wantVal: nil, wantError: true},
	}

	ln := make([]string, 0) // needs pre-allocation to accomplish pre-allocation for labels
	lv := make([]string, 0)

	fqName := prometheusize(prefix + name)
	help := metricHelp(prefix, name)
	typ := prometheus.UntypedValue
	d := prometheus.NewDesc(fqName, help, ln, nil)

	for _, tc := range testCases {
		var want prometheus.Metric
		if tc.wantVal != nil {
			want, _ = prometheus.NewConstMetric(d, typ, *tc.wantVal, lv...)
		}

		m, err := makeRawMetric(prefix, name, tc.value, nil)

		assert.Equal(t, tc.wantError, err != nil)
		assert.Equal(t, want, m)
	}
}
