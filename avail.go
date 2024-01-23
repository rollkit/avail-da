package avail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	"io"
	"log"
	"net/http"
	"net/url"

	"fmt"

	"github.com/rollkit/go-da"
)

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

// Config represents the configuration structure.
type Config struct {
	AppID             uint32 `json:"app_ID"`
	LcURL             string `json:"lc_url"`
	GRPCServerAddress string `json:"grpc_server_address"`
}

// BlockURL represents the URL pattern for retrieving data and extrinsic information
const BlockURL = "/blocks/%d/data?fields=data,extrinsic"

// BLOCK_NOT_FOUND represents the string indicating that a block is not found.
const BLOCK_NOT_FOUND = "\"Not found\""

// PROCESSING_BLOCK represents the string indicating that a block is still being processed.
const PROCESSING_BLOCK = "\"Processing block\""

// AvailDA implements the avail backend for the DA interface
type AvailDA struct {
	config Config
	ctx    context.Context
}

// NewAvailDA returns an instance of AvailDA
func NewAvailDA(config Config, ctx context.Context) *AvailDA {
	return &AvailDA{
		ctx:    ctx,
		config: Config{LcURL: config.LcURL, AppID: config.AppID},
	}
}

var _ da.DA = &AvailDA{}

// MaxBlobSize returns the max blob size
func (c *AvailDA) MaxBlobSize() (uint64, error) {
	var maxBlobSize uint64 = 64 * 64 * 500
	return maxBlobSize, nil
}

// Submit each blob to avail data availability layer
func (c *AvailDA) Submit(daBlobs []da.Blob) ([]da.ID, []da.Proof, error) {
	resultChan := make(chan SubmitResponse, len(daBlobs))
	errorChan := make(chan error, len(daBlobs))

	var wg sync.WaitGroup

	var mu sync.Mutex

	for _, blob := range daBlobs {
		wg.Add(1)

		// Start a goroutine for each blob
		go func(blob da.Blob) {
			defer wg.Done()
			encodedBlob := base64.StdEncoding.EncodeToString(blob)
			requestData := SubmitRequest{
				Data: encodedBlob,
			}

			requestBody, err := json.Marshal(requestData)
			if err != nil {
				errorChan <- err
				return
			}

			// Make a POST request to the /v2/submit endpoint.
			response, err := http.Post(c.config.LcURL+"/submit", "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				errorChan <- err
				return
			}

			defer func() {
				err = response.Body.Close()
				if err != nil {
					log.Println("error closing response body", err)
				}
			}()

			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				errorChan <- err
				return
			}

			var submitResponse SubmitResponse
			err = json.Unmarshal(responseData, &submitResponse)
			if err != nil {
				errorChan <- err
				return
			}

			// Acquire the mutex before updating slices
			mu.Lock()
			resultChan <- SubmitResponse{
				BlockNumber:      submitResponse.BlockNumber,
				BlockHash:        submitResponse.BlockHash,
				TransactionHash:  submitResponse.TransactionHash,
				TransactionIndex: submitResponse.TransactionIndex,
			}
			mu.Unlock()

		}(blob)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results from channels
	var ids []da.ID
	var proofs []da.Proof

	for result := range resultChan {
		ids = append(ids, makeID(result.BlockNumber))
		proofs = append(proofs, makeProofs(result.TransactionHash))
	}

	// Check for errors
	if err := <-errorChan; err != nil {
		return nil, nil, err
	}

	fmt.Println("successfully submitted blobs to avail")
	return ids, proofs, nil
}

// Get returns Blob for each given ID, or an error
func (c *AvailDA) Get(ids []da.ID) ([]da.Blob, error) {
	var blobs [][]byte
	var blockNumber uint32
	for _, id := range ids {
	Loop:
		blockNumber = binary.BigEndian.Uint32(id)
		blocksURL := fmt.Sprintf(c.config.LcURL+BlockURL, blockNumber)
		parsedURL, err := url.Parse(blocksURL)
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest("GET", parsedURL.String(), nil)
		if err != nil {
			return nil, err
		}
		client := http.DefaultClient
		response, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = response.Body.Close()
			if err != nil {
				log.Println("error closing response body", err)
			}
		}()
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var blocksObject BlocksResponse
		if string(responseData) == BLOCK_NOT_FOUND {
			blocksObject = BlocksResponse{BlockNumber: blockNumber, DataTransactions: []DataTransactions{}}
		} else if string(responseData) == PROCESSING_BLOCK {
			time.Sleep(10 * time.Second)
			goto Loop
		} else {
			err = json.Unmarshal(responseData, &blocksObject)
			if err != nil {
				return nil, err
			}
		}
		for _, dataTransaction := range blocksObject.DataTransactions {
			decodeStr, _ := base64.StdEncoding.DecodeString(dataTransaction.Data)
			blobs = append(blobs, []byte(string(decodeStr)))
		}
	}
	return blobs, nil
}

// GetIDs returns the ID
func (c *AvailDA) GetIDs(height uint64) ([]da.ID, error) {
	// todo:currently returning height as ID, need to extend avail-light api
	heightAsUint32 := uint32(height)
	ids := make([]byte, 8)
	binary.BigEndian.PutUint32(ids, heightAsUint32)
	return [][]byte{ids}, nil
}

// Commit creates a Commitment for each given Blob.
func (c *AvailDA) Commit(daBlobs []da.Blob) ([]da.Commitment, error) {
	return nil, nil
}

// Validate validates Commitments against the corresponding Proofs
func (c *AvailDA) Validate(ids []da.ID, daProofs []da.Proof) ([]bool, error) {
	return nil, nil
}

func makeID(blockNumber uint32) da.ID {
	// IDs are not unique in rollkit context and that this has to be improved
	id := make([]byte, 8)
	binary.BigEndian.PutUint32(id, blockNumber)
	return id
}

func makeProofs(proofs string) da.ID {
	return []byte(proofs)
}
