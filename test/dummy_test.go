package dummy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	mux := http.NewServeMux()

	mux.HandleFunc("/v2/blocks/", mockGetEndpoint)
	mux.HandleFunc("/v2/submit", mockSubmitEndpoint)

	server := &http.Server{
		Addr:         ":9000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server
	go func() {
		fmt.Println("Mock Server is running on :9000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting mock server: %v\n", err)
		}
	}()

	time.Sleep(1 * time.Second)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestMockGetEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/v2/blocks/42/data?fields=data,extrinsic", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the mockGetEndpoint function with the fake request and recorder
	mockGetEndpoint(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expected := `{"block_number":42,"data_transactions":[{"data":"TW9ja2VkRGF0YQ==","extrinsic":"MockedExtrinsic"}]}`
	if trimmedBody := strings.TrimSpace(rr.Body.String()); trimmedBody != expected {
		t.Errorf("Handler returned unexpected body:\ngot:\n%v\nwant:\n%v", trimmedBody, expected)
	}
}

func TestMockGetEndpointInvalidBlockNumber(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/blocks/invalid/data?fields=data,extrinsic", nil)
	rr := httptest.NewRecorder()

	// Call the mockGetEndpoint function with the fake request and recorder
	mockGetEndpoint(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusBadRequest)
	}
}

func TestMockSubmitEndpoint(t *testing.T) {
	submitReq := SubmitRequest{
		Data: "TW9ja2VkRGF0YQ==",
	}

	reqBody, err := json.Marshal(submitReq)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v2/submit", strings.NewReader(string(reqBody)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Call the mockSubmitEndpoint function with the fake request and recorder
	mockSubmitEndpoint(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expected := `{"block_number":42,"block_hash":"mocked_block_hash","hash":"mocked_transaction_hash","index":1}`
	if trimmedBody := strings.TrimSpace(rr.Body.String()); trimmedBody != expected {
		t.Errorf("Handler returned unexpected body:\ngot:\n%v\nwant:\n%v", trimmedBody, expected)
	}
}
