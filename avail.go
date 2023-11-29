package avail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rollkit/go-da"
)

type SubmitRequest struct {
	Data string `json:"data"`
}

type SubmitResponse struct {
	BlockNumber      uint32 `json:"block_number"`
	BlockHash        string `json:"block_hash"`
	TransactionHash  string `json:"hash"`
	TransactionIndex uint32 `json:"index"`
}

type BlocksResponse struct {
	BlockNumber      uint32             `json:"block_number"`
	DataTransactions []DataTransactions `json:"data_transactions"`
}

type DataTransactions struct {
	Data      string `json:"data"`
	Extrinsic string `json:"extrinsic"`
}

type Config struct {
	AppID uint32 `json:"app_ID"`
	LcURL string `json:"lc_url"`
}

const BlockURL = "/blocks/%d/data?fields=data,extrinsic"

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

// submits each blob to avail data availability layer
func (c *AvailDA) Submit(daBlobs []da.Blob) ([]da.ID, []da.Proof, error) {
	ids := make([]da.ID, len(daBlobs))
	for index, blob := range daBlobs {
		encodedBlob := base64.StdEncoding.EncodeToString(blob)
		requestData := SubmitRequest{
			Data: encodedBlob,
		}
		requestBody, err := json.Marshal(requestData)
		if err != nil {
			return nil, nil, err
		}
		// Make a POST request to the /v2/submit endpoint.
		response, err := http.Post(c.config.LcURL+"/submit", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			return nil, nil, err
		}
		defer response.Body.Close()

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}

		var submitResponse SubmitResponse
		err = json.Unmarshal(responseData, &submitResponse)
		if err != nil {
			return nil, nil, err
		}
		ids[index] = makeID(submitResponse.BlockNumber)
	}
	fmt.Println("succesfully submitted blobs to avail")
	return ids, nil, nil
}

// Get returns Blob for each given ID, or an error
func (c *AvailDA) Get(ids []da.ID) ([]da.Blob, error) {
	var blobs [][]byte
	var blockNumber uint32
	for _, id := range ids {
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
			_ = response.Body.Close()
		}()
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var blocksObject BlocksResponse
		err = json.Unmarshal(responseData, &blocksObject)
		if err != nil {
			return nil, err
		}
		for _, dataTransaction := range blocksObject.DataTransactions {
			blobs = append(blobs, []byte(dataTransaction.Data))
		}
	}
	return blobs, nil
}

func (c *AvailDA) GetIDs(height uint64) ([]da.ID, error) {
	//todo:currently returning height as ID, need to extend avail-light api
	ids := make([]byte, 8)
	binary.BigEndian.PutUint64(ids, height)
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
