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

	client := getTestClient(t)

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
# HELP mongodb_count TODO
# TYPE mongodb_count untyped
mongodb_count 1
# HELP mongodb_indexDetails_id_LSM_bloom_filter_hits TODO
# TYPE mongodb_indexDetails_id_LSM_bloom_filter_hits untyped
mongodb_indexDetails_id_LSM_bloom_filter_hits 0
# HELP mongodb_indexDetails_id_btree_overflow_pages TODO
# TYPE mongodb_indexDetails_id_btree_overflow_pages untyped
mongodb_indexDetails_id_btree_overflow_pages 0
`)
	// Filter metrics for 2 reasons:
	// 1. The result is huge
	// 2. We need to check against know values. Don't use metrics that return counters like uptime
	//    or counters like the number of transactions because they won't return a known value to compare
	filter := []string{
		"mongodb_count",
		"mongodb_indexDetails_id_LSM_bloom_filter_hits",
		"mongodb_indexDetails_id_btree_overflow_pages",
	}
	err = testutil.CollectAndCompare(c, expected, filter...)
	assert.NoError(t, err)
}
