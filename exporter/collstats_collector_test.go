package exporter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCollStatsCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := getTestClient(ctx, t)

	database := client.Database("testdb")
	database.Drop(ctx) //nolint
	collection := database.Collection("testcol")
	_, err := collection.InsertOne(ctx, bson.M{"f1": 1, "f2": "2"})
	assert.NoError(t, err)

	c := &collstatsCollector{
		ctx:         ctx,
		client:      client,
		collections: []string{"testdb.testcol"},
	}

	// The last \n at the end of this string is important
	expected := strings.NewReader(`
# HELP mongodb_latencyStats_commands_latency latencyStats.commands.
# TYPE mongodb_latencyStats_commands_latency untyped
mongodb_latencyStats_commands_latency 0
# HELP mongodb_latencyStats_commands_ops latencyStats.commands.
# TYPE mongodb_latencyStats_commands_ops untyped
mongodb_latencyStats_commands_ops 0
# HELP mongodb_latencyStats_writes_ops latencyStats.writes.
# TYPE mongodb_latencyStats_writes_ops untyped
mongodb_latencyStats_writes_ops 1` + "\n")

	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_latencyStats_commands_latency",
		"mongodb_latencyStats_commands_ops",
		"mongodb_latencyStats_writes_ops",
	}
	err = testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
