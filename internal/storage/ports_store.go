package storage

import "fmt"

// DocumentWriter saves the document using a unique identifier key, if a document with such key already exists,
// it is overwritten.
type DocumentWriter interface {
	Write(key string, document interface{}) error
}

// PortsStore saves key-document pairs in a database.
type PortsStore struct {
	dw DocumentWriter
}

// NewPortsStore returns a pointer to a new PortsStore.
func NewPortsStore(dw DocumentWriter) *PortsStore {
	return &PortsStore{dw: dw}
}

// Put saves key-port data pair in a database, creates new documents and overwrites existing ones.
func (p *PortsStore) Put(key string, port Port) error {
	if err := p.dw.Write(key, port); err != nil {
		return fmt.Errorf("writing port projection: %w", err)
	}

	return nil
}
