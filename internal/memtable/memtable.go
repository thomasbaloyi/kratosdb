package memtable

import "sync"

type Table struct {
	mutex sync.RWMutex
	table map[string]entry
}

type entry struct {
	value   string
	deleted bool
}

func New() *Table {
	return &Table{table: make(map[string]entry)}
}

func (m *Table) Put(key string, value string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.table[key] = entry{value: value, deleted: false}
}

func (m *Table) Get(key string) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	e, exists := m.table[key]
	if exists && !e.deleted {
		return e.value, exists
	}

	return "", false
}

func (m *Table) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	e, exists := m.table[key]
	if exists {
		e.deleted = true
		m.table[key] = e
	}
}
