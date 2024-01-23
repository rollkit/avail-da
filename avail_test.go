package avail_test

import (
	"context"
	"sync"
	"testing"

	"github.com/rollkit/avail-da"
	dummy "github.com/rollkit/avail-da/test"
)

func TestAvailDA(t *testing.T) {
	config := avail.Config{
		AppID: 1,
		LcURL: "http://localhost:9000/v2",
	}
	ctx := context.Background()

	da := avail.NewAvailDA(config, ctx)

	var wg sync.WaitGroup
	wg.Add(1)

	// Start the mock server in a separate goroutine
	go func() {
		defer wg.Done()
		dummy.StartMockServer()
	}()

	// Wait for the mock server to start
	wg.Wait()

	dummy.RunDATestSuite(t, da)
}
