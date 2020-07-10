// Package tu has Test Util functions
package tu

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// MongoDBS1PrimaryPort MongoDB Shard 1 Primary Port
	MongoDBS1PrimaryPort = "17001"
	// MongoDBStandAlonePort MongoDB stand alone instance Port
	MongoDBStandAlonePort = "27017"
)

// GetenvDefault gets a variable from the environment and returns its value or the
// spacified default if the variable is empty.
func GetenvDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

// DefaultTestClient returns the default MongoDB connection used for tests. It is a direct
// connection to the primary server of replicaset 1.
func DefaultTestClient(ctx context.Context, t *testing.T) *mongo.Client {
	return TestClient(ctx, MongoDBS1PrimaryPort, t)
}

// TestClient returns a new MongoDB connection to the specified server port.
func TestClient(ctx context.Context, port string, t *testing.T) *mongo.Client {
	if port == "" {
		port = MongoDBS1PrimaryPort
	}

	hostname := "127.0.0.1"
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
