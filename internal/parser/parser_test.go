package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkJSONStreamParser_ObjectProperties(b *testing.B) {
	file, err := os.Open("../testfiles/ports.json")
	require.NoError(b, err, "Can't open test input data")

	var memBefore, memAfter runtime.MemStats
	for i := 0; i < b.N; i++ {
		// rewind file before each test
		_, _ = file.Seek(0, 0)
		runtime.ReadMemStats(&memBefore)

		testSuccess(b, file)

		runtime.ReadMemStats(&memAfter)

		fmt.Println("total:", memAfter.TotalAlloc-memBefore.TotalAlloc)
		fmt.Println("mallocs:", memAfter.Mallocs-memBefore.Mallocs)
	}
}

func testSuccess(b *testing.B, file io.ReadCloser) {
	streamer := NewJSONStreamParser()
	err := streamer.ObjectProperties(file, func(decoder *json.Decoder, property json.Token) error { return nil })
	assert.NoError(b, err)
}
