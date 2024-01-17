package dummy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// BlocksResponse represents the structure of a response containing blocks information.
type BlocksResponse struct {
	BlockNumber      uint32             `json:"block_number"`
	DataTransactions []DataTransactions `json:"data_transactions"`
}

// DataTransactions represents data transactions within the block.
type DataTransactions struct {
	Data      string `json:"data"`
	Extrinsic string `json:"extrinsic"`
}

// SubmitRequest represents a request to submit data.
type SubmitRequest struct {
	Data string `json:"data"`
}

// SubmitResponse represents the response after submitting data.
type SubmitResponse struct {
	BlockNumber      uint32 `json:"block_number"`
	BlockHash        string `json:"block_hash"`
	TransactionHash  string `json:"hash"`
	TransactionIndex uint32 `json:"index"`
}

func mockGetEndpoint(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	blockNoStr := parts[3]
	blockNo, err := strconv.ParseUint(blockNoStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid block number", http.StatusBadRequest)
		return
	}

	if blockNo == 42 {
		response := BlocksResponse{
			BlockNumber: uint32(blockNo),
			DataTransactions: []DataTransactions{
				{
					Data:      "TW9ja2VkRGF0YQ==",
					Extrinsic: "MockedExtrinsic",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	} else if blockNo == 43 {
		response := BlocksResponse{
			BlockNumber: uint32(blockNo),
			DataTransactions: []DataTransactions{
				{
					Data:      "TW9ja2VkRGF0YTI=",
					Extrinsic: "MockedExtrinsic2",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	} else if blockNo == 44 {
		response := BlocksResponse{
			BlockNumber: uint32(blockNo),
			DataTransactions: []DataTransactions{
				{
					Data:      "TW9ja2VkRGF0YTM=",
					Extrinsic: "MockedExtrinsic3",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	}
}

func mockSubmitEndpoint(w http.ResponseWriter, r *http.Request) {
	// Extract data from the request body
	var submitReq SubmitRequest
	err := json.NewDecoder(r.Body).Decode(&submitReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if submitReq.Data == "TW9ja2VkRGF0YQ==" {
		response := SubmitResponse{
			BlockNumber:      42,
			BlockHash:        "mocked_block_hash",
			TransactionHash:  "mocked_transaction_hash",
			TransactionIndex: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	} else if submitReq.Data == "TW9ja2VkRGF0YTI=" {
		response := SubmitResponse{
			BlockNumber:      43,
			BlockHash:        "mocked_block_hash2",
			TransactionHash:  "mocked_transaction_hash2",
			TransactionIndex: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	} else if submitReq.Data == "TW9ja2VkRGF0YTM=" {
		response := SubmitResponse{
			BlockNumber:      44,
			BlockHash:        "mocked_block_hash3",
			TransactionHash:  "mocked_transaction_hash3",
			TransactionIndex: 3,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	}
}

// StartMockServer starts a mock server for testing purpose
func StartMockServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/blocks/", mockGetEndpoint)
	mux.HandleFunc("/v2/submit", mockSubmitEndpoint)

	// Create an HTTP server with timeouts
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
}
