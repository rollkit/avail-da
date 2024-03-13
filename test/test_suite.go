package dummy

import (
	"context"
	"encoding/binary"
	"sync"
	"testing"

	"github.com/rollkit/go-da"

	"github.com/stretchr/testify/assert"
)

// RunDATestSuite runs all tests against given DA
func RunDATestSuite(t *testing.T, d da.DA) {
	t.Run("Basic DA test", func(t *testing.T) {
		BasicDATest(t, d)
	})
	t.Run("Get IDs and all data", func(t *testing.T) {
		GetIDsTest(t, d)
	})
	t.Run("Check Errors", func(t *testing.T) {
		CheckErrors(t, d)
	})
	t.Run("Concurrent read/write test", func(t *testing.T) {
		ConcurrentReadWriteTest(t, d)
	})
}

// Blob is a type alias
type Blob = da.Blob

// ID is a type alias
type ID = da.ID

// BasicDATest tests round trip of messages to DA and back.
func BasicDATest(t *testing.T, da da.DA) {
	ctx := context.TODO()
	msg1 := []byte("MockedData")
	msg2 := []byte("MockedData2")
	id1, err := da.Submit(ctx, []Blob{msg1}, -1, nil)
	assert.NoError(t, err)

	expID1 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID1, 42)

	assert.GreaterOrEqual(t, len(id1), 1)
	assert.Equal(t, id1[0], expID1)

	id2, err := da.Submit(ctx, []Blob{msg2}, -1, nil)
	assert.NoError(t, err)

	expID2 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID2, 43)

	assert.GreaterOrEqual(t, len(id2), 1)
	assert.Equal(t, id2[0], expID2)

	ret, err := da.Get(ctx, id1, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	assert.Equal(t, []Blob{msg1}, ret)

	ret, err = da.Get(ctx, id2, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	assert.Equal(t, []Blob{msg2}, ret)

	ids, err := da.Submit(ctx, []Blob{msg1, msg2}, -1, nil)
	assert.NoError(t, err)
	assert.Contains(t, ids, expID1)
	assert.Contains(t, ids, expID2)
}

// CheckErrors ensures that errors are handled properly by DA.
func CheckErrors(t *testing.T, da da.DA) {
	ctx := context.TODO()
	blob, err := da.Get(ctx, []ID{[]byte("invalid")}, nil)
	assert.Error(t, err)
	assert.Empty(t, blob)
}

// GetIDsTest tests iteration over DA
func GetIDsTest(t *testing.T, da da.DA) {
	ctx := context.TODO()

	msg1 := []byte("MockedData")
	msg2 := []byte("MockedData2")

	expID1 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID1, 42)

	expID2 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID2, 43)

	ids, err := da.Submit(ctx, []Blob{msg1, msg2}, -1, nil)
	assert.NoError(t, err)
	assert.Contains(t, ids, expID1)
	assert.Contains(t, ids, expID2)

	var height [][]byte
	var i uint64

	var allBlobs [][]byte

	for i = 42; i < 44; i++ {
		height, err = da.GetIDs(ctx, i, nil)
		if err != nil {
			t.Error("failed to get height:", err)
		}
		blobs, err := da.Get(ctx, height, nil)
		assert.NoError(t, err)
		allBlobs = append(allBlobs, blobs...)
	}
	assert.Contains(t, allBlobs, msg1)
	assert.Contains(t, allBlobs, msg2)
}

// ConcurrentReadWriteTest tests the use of mutex lock in DummyDA by calling separate methods that use `d.data` and making sure there's no race conditions.
func ConcurrentReadWriteTest(t *testing.T, da da.DA) {
	ctx := context.TODO()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, err := da.GetIDs(ctx, i, nil)
			assert.NoError(t, err)
		}
	}()

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, err := da.Submit(ctx, [][]byte{[]byte("MockedData")}, -1, nil)
			assert.NoError(t, err)
		}
	}()

	wg.Wait()
}
