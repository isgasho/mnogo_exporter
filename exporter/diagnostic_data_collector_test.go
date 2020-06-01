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
# HELP mongodb_oplog_stats_max TODO
# TYPE mongodb_oplog_stats_max untyped
mongodb_oplog_stats_max -1
# HELP mongodb_oplog_stats_nindexes TODO
# TYPE mongodb_oplog_stats_nindexes untyped
mongodb_oplog_stats_nindexes 0
# HELP mongodb_rs_members_configVersion TODO
# TYPE mongodb_rs_members_configVersion untyped
mongodb_rs_members_configVersion{member_idx="127.0.0.1:17001"} 1
mongodb_rs_members_configVersion{member_idx="127.0.0.1:17002"} 1
mongodb_rs_members_configVersion{member_idx="127.0.0.1:17003"} 1
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
		"mongodb_oplog_stats_max",
		"mongodb_oplog_stats_nindexes",
		"mongodb_rs_members_configVersion",
		"mongodb_rs_members_state",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
