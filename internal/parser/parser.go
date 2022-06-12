package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// JSONDataHandler defines the function signature used by JSONStreamParser to process the returned JSON data.
type JSONDataHandler func(decoder *json.Decoder, property json.Token) error

// JSONStreamParser allows JSON tokens to be read from io.Reader without first loading the entire data into memory.
type JSONStreamParser struct {
	source io.Reader
}

// NewJSONStreamParser returns a pointer to a new JSONStreamParser.
func NewJSONStreamParser() *JSONStreamParser {
	return &JSONStreamParser{}
}

// ObjectProperties reads the properties of the JSON object one by one and passes their values to the supplied
// JSONDataHandler function. Returns an error when the data in the source is not a JSON object.
func (p *JSONStreamParser) ObjectProperties(source io.Reader, dataHandler JSONDataHandler) error {
	reader := bufio.NewReader(source)
	decoder := json.NewDecoder(reader)

	openingToken, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("reading opening JSON token: %w", err)
	}

	if openingToken != json.Delim('{') {
		return fmt.Errorf("the source data is not a JSON object")
	}

	var propName json.Token
	for decoder.More() {
		propName, err = decoder.Token()
		if err != nil {
			return fmt.Errorf("reading token: %w", err)
		}

		if err = dataHandler(decoder, propName); err != nil {
			return fmt.Errorf("handling data: %w", err)
		}
	}

	return nil
}
