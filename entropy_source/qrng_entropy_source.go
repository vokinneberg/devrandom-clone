package entropy_source

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Entropy source configuration.
const (
	EntropySourceUrl   = "https://qrng.anu.edu.au/API/jsonI.php"
	EntropyDataType    = "hex16"
	EntropyArrayLength = "1024"
	EntropyBlockSize   = "1024"
)

// Response of ANU Quantum Random Numbers Server generator.
type QrngEntropyResponse struct {
	Type    string   `json:"type"`
	Length  int      `json:"length"`
	Size    int      `json:"size"`
	Data    []string `json:"data"`
	Success bool     `json:"success"`
}

// Main abstraction for entropy source.
type EntropySource interface {
	// Eternally provides entropy (random numbers) until error happened or entropy source exhausted.
	Entropy(entropy chan []byte, error chan error)
}

type qrngEntropySource struct {
	HTTPClient *http.Client
}

// Create instance of qrngEntropySource based on ANU Quantum Random Numbers Server generator.
func NewQrngEntropySource(httpClient *http.Client) EntropySource {
	return &qrngEntropySource{
		HTTPClient: httpClient,
	}
}

func (s *qrngEntropySource) Entropy(entropy chan []byte, error chan error) {
	v := url.Values{}
	v.Add("length", EntropyArrayLength)
	v.Add("type", EntropyDataType)
	v.Add("size", EntropyBlockSize)

	for {
		resp, err := s.HTTPClient.Get(EntropySourceUrl + "?" + v.Encode())
		if err != nil {
			error <- err
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			error <- err
			return
		}

		var target QrngEntropyResponse
		if err = json.Unmarshal(body, &target); err != nil {
			error <- err
			return
		}

		if !target.Success {
			error <- errors.New("unable to retrieve entropy")
			return
		}

		for _, block := range target.Data {
			entropy <- []byte(block)
		}
	}
}
