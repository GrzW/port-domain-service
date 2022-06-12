package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GrzW/port-domain-service/internal/parser"
	"github.com/GrzW/port-domain-service/internal/storage"

	ta "github.com/stretchr/testify/assert"
	tr "github.com/stretchr/testify/require"
)

func TestPortsUpdatedHandler_HandlePortsUpdated(t *testing.T) {
	assert := ta.New(t)
	require := tr.New(t)

	file, err := ioutil.ReadFile("../testfiles/ports.json")
	require.NoError(err, "Can't open test input data")

	streamParser := parser.NewJSONStreamParser()

	srv := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.WriteHeader(http.StatusOK)

		_, err = resp.Write(file)
		require.NoError(err)
	}))

	fakeWriter := &mockPortsWriter{writeReturns: nil}
	handler := NewPortsUpdatedHandler(fakeWriter, streamParser)

	err = handler.HandlePortsUpdated(context.Background(), srv.URL)
	assert.NoError(err)
	assert.Equal(1596, fakeWriter.writeCalledCount)
}

type mockPortsWriter struct {
	writeReturns     error
	writeCalledCount int
}

func (mpw *mockPortsWriter) Put(string, storage.Port) error {
	mpw.writeCalledCount++

	return mpw.writeReturns
}
