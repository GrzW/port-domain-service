package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GrzW/port-domain-service/internal/parser"
	"github.com/GrzW/port-domain-service/internal/storage"
)

// JSONStreamParser allows parsing JSON data from io.Reader without loading everything into memory first.
type JSONStreamParser interface {
	// ObjectProperties assumes the data in io.Reader is a JSON object and parses its properties one by one.
	ObjectProperties(source io.Reader, handler parser.JSONDataHandler) error
}

// PortsWriter saves Port data in a database.
type PortsWriter interface {
	// Put saves new key-Port pair or overwrites existing one.
	Put(key string, port storage.Port) error
}

// PortsUpdatedHandler handles ports data updates.
type PortsUpdatedHandler struct {
	store  PortsWriter
	parser JSONStreamParser
}

// NewPortsUpdatedHandler return a pointer to a new PortsUpdatedHandler.
func NewPortsUpdatedHandler(store PortsWriter, streamParser JSONStreamParser) *PortsUpdatedHandler {
	return &PortsUpdatedHandler{
		store:  store,
		parser: streamParser,
	}
}

// HandlePortsUpdated downloads data contained in file located at fileURL and, port by port, saves it in a database.
func (puh *PortsUpdatedHandler) HandlePortsUpdated(ctx context.Context, fileURL string) error {
	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var (
		ports  = make(chan storage.Port)
		errors = make(chan error)
	)
	go func() {
		defer close(errors)
		defer close(ports)

		if err = puh.parser.ObjectProperties(resp.Body, getJSONDataHandler(ctx, ports)); err != nil {
			errors <- err
		}
	}()

	for {
		select {
		case port, ok := <-ports:
			if !ok { // closed channel, means that we finished processing the file
				return nil
			}

			if err = puh.store.Put(port.Name, port); err != nil {
				fmt.Printf("Warning: failed to save port %q data: %s", port.Name, err.Error())
			}
		case err = <-errors:
			if err != nil {
				return err
			}
		}
	}
}

func getJSONDataHandler(ctx context.Context, ports chan storage.Port) parser.JSONDataHandler {
	var data storage.Port

	return func(decoder *json.Decoder, property json.Token) error {
		if err := decoder.Decode(&data); err != nil {
			fmt.Printf("Failed to decode %q: %s", fmt.Sprint(property), err.Error())

			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled or deadline exceeded")
		case ports <- data:
		}

		return nil
	}
}
