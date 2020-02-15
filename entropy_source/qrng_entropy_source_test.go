package entropy_source

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	ValidMockResponseBody = `{"type":"string","length":10,"size":2,
					"data":["f7b5","e4d1","c842","5551","f117","e75d","30a1","32b8","8e40","9ba3"],
					"success":true}`
	InvalidMockResponseBody = `{"type":"string","length":10,"size":2,
					"data":"f7b5","e4d1","c842","5551","f117","e75d","30a1","32b8","8e40","9ba3"],
					"success":true}`
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func MockHttpClient(body string, withTimeout bool) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			var err error = nil
			body := ioutil.NopCloser(bytes.NewBufferString(body))
			statusCode := http.StatusOK
			if withTimeout {
				body = nil
				statusCode = http.StatusRequestTimeout
				err = errors.New("dial tcp: lookup qrng.anu.edu.au: no such host")
			}
			return &http.Response{
				StatusCode: statusCode,
				Body:       body,
				Header:     make(http.Header)}, err
		})}
}

func TestEntropySource_Timeout(t *testing.T) {
	entropyChan := make(chan []byte)
	errChan := make(chan error)
	mockedClient := MockHttpClient(ValidMockResponseBody, true)
	entropySource := NewQrngEntropySource(mockedClient)
	go entropySource.Entropy(entropyChan, errChan)

	select {
	case err := <-errChan:
		require.Equal(t, fmt.Sprint(err), "Get https://qrng.anu.edu.au/API/jsonI.php?"+
			"length=1024&size=1024&type=hex16: "+
			"dial tcp: lookup qrng.anu.edu.au: no such host")
		return
	case entropy := <-entropyChan:
		t.Errorf("no error thrown. got entropy: %s", string(entropy))
		return
	}
}

func TestEntropySource_InvalidJsonBody(t *testing.T) {
	entropyChan := make(chan []byte)
	errChan := make(chan error)
	mockedClient := MockHttpClient(InvalidMockResponseBody, false)
	entropySource := NewQrngEntropySource(mockedClient)
	go entropySource.Entropy(entropyChan, errChan)

	select {
	case err := <-errChan:
		require.Equal(t, fmt.Sprint(err), "invalid character ',' after object key")
		return
	case entropy := <-entropyChan:
		t.Errorf("no error thrown. got entropy: %s", string(entropy))
		return
	}
}

func TestEntropySource_Valid(t *testing.T) {
	entropyChan := make(chan []byte)
	errChan := make(chan error)
	mockedClient := MockHttpClient(ValidMockResponseBody, false)
	entropySource := NewQrngEntropySource(mockedClient)
	go entropySource.Entropy(entropyChan, errChan)

	select {
	case err := <-errChan:
		require.NoError(t, err)
		return
	case entropy := <-entropyChan:
		require.Equal(t, string(entropy), "f7b5")
		return
	}
}
