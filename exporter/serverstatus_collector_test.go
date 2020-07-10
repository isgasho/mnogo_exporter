package exporter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/Percona-Lab/mnogo_exporter/internal/tu"
)

func TestServerStatusDataCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := tu.DefaultTestClient(ctx, t)

	c := &serverStatusCollector{
		client: client,
	}

	// The last \n at the end of this string is important
	expected := strings.NewReader(`
# HELP mongodb_mem_bits mem.
# TYPE mongodb_mem_bits untyped
mongodb_mem_bits 64
# HELP mongodb_metrics_commands_cloneCollection_failed metrics.commands.cloneCollection.
# TYPE mongodb_metrics_commands_cloneCollection_failed untyped
mongodb_metrics_commands_cloneCollection_failed 0
# HELP mongodb_metrics_commands_connPoolSync_failed metrics.commands.connPoolSync.
# TYPE mongodb_metrics_commands_connPoolSync_failed untyped
mongodb_metrics_commands_connPoolSync_failed 0
# HELP mongodb_wiredTiger_log_slot_join_calls_yielded wiredTiger.log.
# TYPE mongodb_wiredTiger_log_slot_join_calls_yielded untyped
mongodb_wiredTiger_log_slot_join_calls_yielded 0` + "\n")
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_mem_bits",
		"mongodb_metrics_commands_cloneCollection_failed",
		"mongodb_metrics_commands_connPoolSync_failed",
		"mongodb_wiredTiger_log_slot_join_calls_yielded",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
