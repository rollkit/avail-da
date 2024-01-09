package avail_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rollkit/go-da"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/rollkit/avail-da"
	"github.com/rollkit/go-da/proxy"
	goDATest "github.com/rollkit/go-da/test"
)

func TestMain(m *testing.M) {
	srv := startMockGRPCServ()
	if srv == nil {
		os.Exit(1)
	}
	exitCode := m.Run()

	// teardown servers
	srv.GracefulStop()

	os.Exit(exitCode)
}

func startMockGRPCServ() *grpc.Server {
	srv := proxy.NewServer(goDATest.NewDummyDA(), grpc.Creds(insecure.NewCredentials()))
	lis, err := net.Listen("tcp", "127.0.0.1"+":"+strconv.Itoa(5000))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	go func() {
		_ = srv.Serve(lis)
	}()
	return srv
}

// Blob is a type alias
type Blob = da.Blob

// ID is a type alias
type ID = da.ID

func BasicDATest(t *testing.T, da da.DA) {
	msg1 := []byte("message 1")
	msg2 := []byte("message 2")

	id1, proof1, err := da.Submit([]Blob{msg1})
	assert.NoError(t, err)
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, proof1)

	id2, proof2, err := da.Submit([]Blob{msg2})
	assert.NoError(t, err)
	assert.NotEmpty(t, id2)
	assert.NotEmpty(t, proof2)

	id3, proof3, err := da.Submit([]Blob{msg1})
	assert.NoError(t, err)
	assert.NotEmpty(t, id3)
	assert.NotEmpty(t, proof3)

	assert.NotEqual(t, id1, id2)
	assert.NotEqual(t, id1, id3)

	ret, err := da.Get(id1)
	assert.NoError(t, err)
	assert.Equal(t, []Blob{msg1}, ret)

	ret, err = da.Get(id2)
	assert.NoError(t, err)
	assert.Equal(t, []Blob{msg2}, ret)
}

func TestAvailDA(t *testing.T) {
	config := avail.Config{
		AppID: 1,
		LcURL: "http://localhost:8000/v2",
	}
	ctx := context.Background()

	da := avail.NewAvailDA(config, ctx)
	BasicDATest(t, da)
}
