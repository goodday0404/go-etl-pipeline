package api

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerStart(t *testing.T) {
	/* Arrange */
	testServerAddress := "127.0.0.1:0"
	testServer := NewServer(testServerAddress, 30*time.Second)
	listener, err := net.Listen("tcp", testServer.Addr)
	assert.NoError(t, err)
	expectedJSON := `{"status":"ok"}`
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	go func() {
		if err := testServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Error("Failed to run test server")
		}
	}()

	defer testServer.Close()

	networkAddress := listener.Addr()
	testURL := fmt.Sprintf("http://%s/healthz", networkAddress)

	/* Act */
	resp, err := client.Get(testURL)
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	/* Assert */
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, expectedJSON, string(body))
}
