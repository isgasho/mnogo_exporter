package exporter

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDiagnosticDataCollector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := getTestClient(t)

	c := &diagnosticDataCollector{
		ctx:    ctx,
		client: client,
	}

	expected := strings.NewReader(`
	`)
	err := testutil.CollectAndCompare(c, expected)
	assert.NoError(t, err)
}

func TestMakeMetrics(t *testing.T) {
	hostname := "127.0.0.1"
	username := os.Getenv("TEST_MONGODB_ADMIN_USERNAME")
	password := os.Getenv("TEST_MONGODB_ADMIN_PASSWORD")
	port := os.Getenv("TEST_MONGODB_S1_PRIMARY_PORT")
	username = "admin"
	password = "admin123456"
	port = "17001"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("mongodb://%s:%s@%s:%s/admin", username, password, hostname, port)
	client, err := connect(ctx, dsn)
	assert.NoError(t, err)
	defer client.Disconnect(ctx)

	c := &diagnosticDataCollector{
		ctx:    ctx,
		client: client,
	}

	cmd := bson.D{{"getDiagnosticData", "1"}}
	res := c.client.Database("admin").RunCommand(c.ctx, cmd)
	var m bson.M
	err = res.Decode(&m)
	assert.NoError(t, err)

	buf, err := bson.Marshal(m["data"])
	assert.NoError(t, err)
	err = ioutil.WriteFile(path.Join("../internal", "renamer", "testdata", "test001_src.json"), buf, os.ModePerm)
	assert.NoError(t, err)

	var obj bson.M

	buf, err = ioutil.ReadFile(path.Join("../internal", "renamer", "testdata", "test001_src.json"))
	assert.NoError(t, err)
	err = bson.Unmarshal(buf, &obj)
	assert.NoError(t, err)

	mm := c.makeMetrics("", obj)

	pretty.Println(mm)

}
