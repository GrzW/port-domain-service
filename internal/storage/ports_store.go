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

// Put saves key-port data pair in a database, creates new documents and overwrites existing ones.
func (p *PortsStore) Put(key string, port Port) error {
	if err := p.dw.Write(key, port); err != nil {
		return fmt.Errorf("writing port projection: %w", err)
	}

	return nil
}
