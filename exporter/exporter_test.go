package exporter

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTestClient(t *testing.T) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	hostname := "127.0.0.1"
	username := os.Getenv("TEST_MONGODB_ADMIN_USERNAME")
	password := os.Getenv("TEST_MONGODB_ADMIN_PASSWORD")
	port := os.Getenv("TEST_MONGODB_S1_PRIMARY_PORT") // standalone instance
	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/admin", username, password, hostname, port)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	require.NoError(t, err)

	t.Cleanup(func() { client.Disconnect(ctx) }) //nolint:errcheck

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	return client
}

func TestConnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	hostname := "127.0.0.1"
	username := os.Getenv("TEST_MONGODB_ADMIN_USERNAME")
	password := os.Getenv("TEST_MONGODB_ADMIN_PASSWORD")

	ports := map[string]string{
		"standalone":          os.Getenv("TEST_MONGODB_STANDALONE_PORT"),
		"shard-1 primary":     os.Getenv("TEST_MONGODB_S1_PRIMARY_PORT"),
		"shard-1 secondary-1": os.Getenv("TEST_MONGODB_S1_SECONDARY1_PORT"),
		"shard-1 secondary-2": os.Getenv("TEST_MONGODB_S1_SECONDARY2_PORT"),
		"shard-2 primary":     os.Getenv("TEST_MONGODB_S2_PRIMARY_PORT"),
		"shard-2 secondary-1": os.Getenv("TEST_MONGODB_S2_SECONDARY1_PORT"),
		"shard-2 secondary-2": os.Getenv("TEST_MONGODB_S2_SECONDARY2_PORT"),
		"config server 1":     os.Getenv("TEST_MONGODB_CONFIGSVR1_PORT"),
		"mongos":              os.Getenv("TEST_MONGODB_MONGOS_PORT"),
	}

	t.Run("Connect without SSL", func(t *testing.T) {
		for name, port := range ports {
			dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/admin", username, password, hostname, port)
			client, err := connect(ctx, dsn)
			assert.NoError(t, err, name)
			err = client.Disconnect(ctx)
			assert.NoError(t, err, name)
		}
	})

	t.Run("Connect with SSL", func(t *testing.T) {
		sslOpts := "ssl=true&tlsInsecure=true&tlsCertificateKeyFile=../docker/test/ssl/client.pem"
		for name, port := range ports {
			dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/admin?%s", username, password, hostname, port, sslOpts)
			client, err := connect(ctx, dsn)
			assert.NoError(t, err, name)
			err = client.Disconnect(ctx)
			assert.NoError(t, err, name)
		}
	})
}
