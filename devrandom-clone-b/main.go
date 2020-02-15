package main

import (
	"devrandom-clone/entropy_source"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	HttpClientTimout = 20
)

func main() {
	entropyChan := make(chan []byte)
	errChan := make(chan error)
	entropySource := entropy_source.NewQrngEntropySource(&http.Client{Timeout: HttpClientTimout * time.Second})

	// Will retrieve random numbers eternally until timeout or error happened.
	go entropySource.Entropy(entropyChan, errChan)

	// Handing random numbers and write them to the stdout.
	go func() {
		for {
			if entropy, ok := <-entropyChan; !ok {
				return
			} else if len(entropy) > 0 {
				dst := make([]byte, hex.DecodedLen(len(entropy)))
				if _, err := hex.Decode(dst, entropy); err != nil {
					errChan <- err
				}
				if err := binary.Write(os.Stdout, binary.LittleEndian, dst); err != nil {
					errChan <- errors.New("unable to write to stdout")
				}
			} else {
				errChan <- errors.New("entropy source exhausted")
			}
		}
	}()

	// Block until error received.
	err := <-errChan
	close(entropyChan)
	close(errChan)
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(-1)
}
