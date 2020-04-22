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
# HELP rs_members_configVersion TODO
# TYPE rs_members_configVersion untyped
rs_members_configVersion{member_idx="127.0.0.1:17001"} 1
rs_members_configVersion{member_idx="127.0.0.1:17002"} 1
rs_members_configVersion{member_idx="127.0.0.1:17003"} 1
# HELP rs_members_syncSourceId TODO
# TYPE rs_members_syncSourceId untyped
rs_members_syncSourceId{member_idx="127.0.0.1:17001"} -1
rs_members_syncSourceId{member_idx="127.0.0.1:17002"} 0
rs_members_syncSourceId{member_idx="127.0.0.1:17003"} 0
# HELP ss_metrics_commands_count_failed TODO
# TYPE ss_metrics_commands_count_failed untyped
ss_metrics_commands_count_failed 0
# HELP ss_wt_session_table_salvage_failed_calls TODO
# TYPE ss_wt_session_table_salvage_failed_calls untyped
ss_wt_session_table_salvage_failed_calls 0
`)
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"ss_metrics_commands_count_failed",
		"rs_members_syncSourceId",
		"rs_members_configVersion",
		"ss_wt_session_table_salvage_failed_calls",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
