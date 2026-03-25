package memtable

import (
	"sync"

	sm "github.com/egregors/sortedmap"
)

type Table struct {
	mutex sync.RWMutex
	table *sm.SortedMap[map[string]entry, string, entry]
}

type entry struct {
	value   string
	deleted bool
}

func New() *Table {
	sortedMap := sm.NewFromMap(make(map[string]entry),
		func(i, j sm.KV[string, entry]) bool {
			return i.Key < j.Key
		})
	return &Table{table: sortedMap}
}

func (m *Table) Put(key string, value string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.table.Insert(key, entry{value: value, deleted: false})
}

func (m *Table) Get(key string) (string, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	e, exists := m.table.Get(key)
	if exists && !e.deleted {
		return e.value, exists
	}

	return "", false
}

func (m *Table) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	e, exists := m.table.Get(key)
	if exists {
		e.deleted = true
		m.table.Insert(key, e)
	}
}
