package exporter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestServerStatusDataCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := getTestClient(ctx, t)

	c := &diagnosticDataCollector{
		ctx:    ctx,
		client: client,
	}

	// The last \n at the end of this string is important
	expected := strings.NewReader(`
# HELP mongodb_oplog_stats_capped local.oplog.rs.stats.
# TYPE mongodb_oplog_stats_capped untyped
mongodb_oplog_stats_capped 1
# HELP mongodb_oplog_stats_nindexes local.oplog.rs.stats.
# TYPE mongodb_oplog_stats_nindexes untyped
mongodb_oplog_stats_nindexes 0
# HELP mongodb_oplog_stats_wt_cache_overflow_pages_read_into_cache local.oplog.rs.stats.wiredTiger.cache.
# TYPE mongodb_oplog_stats_wt_cache_overflow_pages_read_into_cache untyped
mongodb_oplog_stats_wt_cache_overflow_pages_read_into_cache 0
# HELP mongodb_oplog_stats_wt_cursor_remove_calls local.oplog.rs.stats.wiredTiger.cursor.
# TYPE mongodb_oplog_stats_wt_cursor_remove_calls untyped
mongodb_oplog_stats_wt_cursor_remove_calls 0
`)
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_oplog_stats_capped",
		"mongodb_oplog_stats_nindexes",
		"mongodb_oplog_stats_wt_cache_overflow_pages_read_into_cache",
		"mongodb_oplog_stats_wt_cursor_remove_calls",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
