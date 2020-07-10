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

func TestReplsetStatusCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := tu.DefaultTestClient(ctx, t)

	c := &replSetGetStatusCollector{
		client: client,
	}

	// The last \n at the end of this string is important
	expected := strings.NewReader(`
                # HELP mongodb_myState myState
                # TYPE mongodb_myState untyped
                mongodb_myState 1
                # HELP mongodb_ok ok
                # TYPE mongodb_ok untyped
                mongodb_ok 1
                # HELP mongodb_optimes_appliedOpTime_t optimes.appliedOpTime.
                # TYPE mongodb_optimes_appliedOpTime_t untyped
                mongodb_optimes_appliedOpTime_t 1
                # HELP mongodb_optimes_durableOpTime_t optimes.durableOpTime.
                # TYPE mongodb_optimes_durableOpTime_t untyped
                mongodb_optimes_durableOpTime_t 1` + "\n")
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_myState",
		"mongodb_ok",
		"mongodb_optimes_appliedOpTime_t",
		"mongodb_optimes_durableOpTime_t",
	}
	err := testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}

func TestReplsetStatusCollectorNoSharding(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := tu.TestClient(ctx, tu.MongoDBStandAlonePort, t)

	c := &replSetGetStatusCollector{
		client: client,
	}

	expected := strings.NewReader(``)
	err := testutil.CollectAndCompare(c, expected)
	assert.Error(t, err)
}
