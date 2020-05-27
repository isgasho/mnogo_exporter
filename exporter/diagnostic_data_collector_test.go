package exporter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDiagnosticDataCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := getTestClient(t)

	c := &diagnosticDataCollector{
		ctx:    ctx,
		client: client,
	}

	// The last \n at the end of this string is important
	expected := strings.NewReader(`
# HELP mongodb_oplog_stats_ok TODO
# TYPE mongodb_oplog_stats_ok untyped
mongodb_oplog_stats_ok 1
# HELP mongodb_oplog_stats_sleepCount TODO
# TYPE mongodb_oplog_stats_sleepCount untyped
mongodb_oplog_stats_sleepCount 0
# HELP mongodb_oplog_stats_wt_LSM_bloom_filter_misses TODO
# TYPE mongodb_oplog_stats_wt_LSM_bloom_filter_misses untyped
mongodb_oplog_stats_wt_LSM_bloom_filter_misses 0
# HELP mongodb_oplog_stats_wt_btree_fixed_record_size TODO
# TYPE mongodb_oplog_stats_wt_btree_fixed_record_size untyped
mongodb_oplog_stats_wt_btree_fixed_record_size 0
# HELP mongodb_rs_members_id TODO
# TYPE mongodb_rs_members_id untyped
mongodb_rs_members_id{member_idx="127.0.0.1:17001"} 0
mongodb_rs_members_id{member_idx="127.0.0.1:17002"} 1
mongodb_rs_members_id{member_idx="127.0.0.1:17003"} 2
# HELP mongodb_rs_members_state TODO
# TYPE mongodb_rs_members_state untyped
mongodb_rs_members_state{member_idx="127.0.0.1:17001"} 1
mongodb_rs_members_state{member_idx="127.0.0.1:17002"} 2
mongodb_rs_members_state{member_idx="127.0.0.1:17003"} 2
`)
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_oplog_stats_ok",
		"mongodb_oplog_stats_sleepCount",
		"mongodb_oplog_stats_wt_LSM_bloom_filter_misses",
		"mongodb_oplog_stats_wt_btree_fixed_record_size",
		"mongodb_rs_members_id",
		"mongodb_rs_members_state",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
