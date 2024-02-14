package dummy

import (
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
	msg1 := []byte("MockedData")
	msg2 := []byte("MockedData2")
	id1, proof1, err := da.Submit([]Blob{msg1})
	assert.NoError(t, err)

	expID1 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID1, 42)

	assert.GreaterOrEqual(t, len(id1), 1)
	assert.Equal(t, id1[0], expID1)
	assert.GreaterOrEqual(t, len(proof1), 1)
	assert.Equal(t, proof1[0], []byte("mocked_transaction_hash"))

	id2, proof2, err := da.Submit([]Blob{msg2})
	assert.NoError(t, err)

	expID2 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID2, 43)

	assert.GreaterOrEqual(t, len(id2), 1)
	assert.Equal(t, id2[0], expID2)
	assert.GreaterOrEqual(t, len(proof2), 1)
	assert.Equal(t, proof2[0], []byte("mocked_transaction_hash2"))

	ret, err := da.Get(id1)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	assert.Equal(t, []Blob{msg1}, ret)

	ret, err = da.Get(id2)
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	assert.Equal(t, []Blob{msg2}, ret)

	ids, proofs, err := da.Submit([]Blob{msg1, msg2})
	assert.NoError(t, err)
	assert.Contains(t, ids, expID1)
	assert.Contains(t, ids, expID2)
	assert.Contains(t, proofs, []byte("mocked_transaction_hash"))
	assert.Contains(t, proofs, []byte("mocked_transaction_hash2"))
}

// CheckErrors ensures that errors are handled properly by DA.
func CheckErrors(t *testing.T, da da.DA) {
	blob, err := da.Get([]ID{[]byte("invalid")})
	assert.Error(t, err)
	assert.Empty(t, blob)
}

// GetIDsTest tests iteration over DA
func GetIDsTest(t *testing.T, da da.DA) {
	msg1 := []byte("MockedData")
	msg2 := []byte("MockedData2")

	expID1 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID1, 42)

	expID2 := make([]byte, 8)
	binary.BigEndian.PutUint32(expID2, 43)

	ids, proofs, err := da.Submit([]Blob{msg1, msg2})
	assert.NoError(t, err)
	assert.Contains(t, ids, expID1)
	assert.Contains(t, ids, expID2)
	assert.Contains(t, proofs, []byte("mocked_transaction_hash"))
	assert.Contains(t, proofs, []byte("mocked_transaction_hash2"))

	var height [][]byte
	var i uint64

	var allBlobs [][]byte

	for i = 42; i < 44; i++ {
		height, err = da.GetIDs(i)
		if err != nil {
			t.Error("failed to get height:", err)
		}
		blobs, err := da.Get(height)
		assert.NoError(t, err)
		allBlobs = append(allBlobs, blobs...)
	}
	assert.Contains(t, allBlobs, msg1)
	assert.Contains(t, allBlobs, msg2)
}

// ConcurrentReadWriteTest tests the use of mutex lock in DummyDA by calling separate methods that use `d.data` and making sure there's no race conditions.
func ConcurrentReadWriteTest(t *testing.T, da da.DA) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, err := da.GetIDs(i)
			assert.NoError(t, err)
		}
	}()

	go func() {
		defer wg.Done()
		for i := uint64(1); i <= 100; i++ {
			_, _, err := da.Submit([][]byte{[]byte("MockedData")})
			assert.NoError(t, err)
		}
	}()

	wg.Wait()
}
