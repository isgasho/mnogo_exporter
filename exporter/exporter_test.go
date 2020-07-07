package exporter

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getEnvDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getTestClient(ctx context.Context, t *testing.T) *mongo.Client {
	hostname := "127.0.0.1"
	port := getEnvDefault("TEST_MONGODB_S1_PRIMARY_PORT", "17001") // standalone instance
	direct := true
	to := time.Second
	co := &options.ClientOptions{
		ConnectTimeout: &to,
		Hosts:          []string{net.JoinHostPort(hostname, port)},
		Direct:         &direct,
	}

	client, err := mongo.Connect(ctx, co)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := client.Disconnect(ctx)
		assert.NoError(t, err)
	})

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	return client
}

func TestConnect(t *testing.T) {
	hostname := "127.0.0.1"
	ctx := context.Background()

	ports := map[string]string{
		"standalone":          getEnvDefault("TEST_MONGODB_STANDALONE_PORT", "27017"),
		"shard-1 primary":     getEnvDefault("TEST_MONGODB_S1_PRIMARY_PORT", "17001"),
		"shard-1 secondary-1": getEnvDefault("TEST_MONGODB_S1_SECONDARY1_PORT", "17002"),
		"shard-1 secondary-2": getEnvDefault("TEST_MONGODB_S1_SECONDARY2_PORT", "17003"),
		"shard-2 primary":     getEnvDefault("TEST_MONGODB_S2_PRIMARY_PORT", "17004"),
		"shard-2 secondary-1": getEnvDefault("TEST_MONGODB_S2_SECONDARY1_PORT", "17005"),
		"shard-2 secondary-2": getEnvDefault("TEST_MONGODB_S2_SECONDARY2_PORT", "17006"),
		"config server 1":     getEnvDefault("TEST_MONGODB_CONFIGSVR1_PORT", "17007"),
		"mongos":              getEnvDefault("TEST_MONGODB_MONGOS_PORT", "17000"),
	}

	t.Run("Connect without SSL", func(t *testing.T) {
		for name, port := range ports {
			dsn := fmt.Sprintf("mongodb://%s:%s/admin", hostname, port)
			client, err := connect(ctx, dsn)
			assert.NoError(t, err, name)
			err = client.Disconnect(ctx)
			assert.NoError(t, err, name)
		}
	})
}
