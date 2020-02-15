package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"
)

func Test_Main(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	// Redefining stdout.
	os.Stdout = w

	go func() {
		main()
	}()

	outChan := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		require.NoError(t, err)
		outChan <- buf.Bytes()
	}()

	// Waiting for 30 seconds for collecting main() program output.
	time.Sleep(30 * time.Second)
	w.Close()
	out := <-outChan

	// Restoring stdout.
	os.Stdout = oldStdout
	assert.NotEmpty(t, out)
}
