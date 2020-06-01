package exporter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestReplsetStatusCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := getTestClient(t)

	c := &replSetGetStatusCollector{
		ctx:    ctx,
		client: client,
	}

	// The last \n at the end of this string is important
	expected := strings.NewReader(`
# HELP mongodb_configServerState_opTime_t TODO
# TYPE mongodb_configServerState_opTime_t untyped
mongodb_configServerState_opTime_t 1
# HELP mongodb_electionCandidateMetrics_electionTerm TODO
# TYPE mongodb_electionCandidateMetrics_electionTerm untyped
mongodb_electionCandidateMetrics_electionTerm 1
# HELP mongodb_electionCandidateMetrics_numCatchUpOps TODO
# TYPE mongodb_electionCandidateMetrics_numCatchUpOps untyped
mongodb_electionCandidateMetrics_numCatchUpOps 0
# HELP mongodb_members_id TODO
# TYPE mongodb_members_id untyped
mongodb_members_id{member_idx="127.0.0.1:17001"} 0
mongodb_members_id{member_idx="127.0.0.1:17002"} 1
mongodb_members_id{member_idx="127.0.0.1:17003"} 2
`)
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_configServerState_opTime_t",
		"mongodb_electionCandidateMetrics_electionTerm",
		"mongodb_electionCandidateMetrics_numCatchUpOps",
		"mongodb_members_id",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}