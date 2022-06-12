package storage

import "sync"

// MemoryDB is nothing more than a key-value store in memory.
type MemoryDB struct {
	db map[string]interface{}
	mx *sync.Mutex
}

// Write saves key-document pair in MemoryDB.
func (m *MemoryDB) Write(key string, document interface{}) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.db[key] = document

	return nil
}
