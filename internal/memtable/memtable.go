package memtable

type MemTable struct {
	table map[string]entry // TODO: Add concurrency protection + use a sorted map
}

type entry struct {
	value 	string
	deleted bool
}

func NewMemTable() *MemTable {
	dict := make(map[string]entry)
	return &MemTable{table: dict}
}

func (m *MemTable) Put(key string, value string) {
	m.table[key] = entry{value: value, deleted: false}
}

func (m *MemTable) Get(key string) (string, bool) {
	entry, exists := m.table[key]
	if exists && !entry.deleted {
		return entry.value, exists
	}
	 
	return "", false
}

func (m *MemTable) Delete(key string) {
	entry, exists := m.table[key] 
	if exists {
		entry.deleted = true
		m.table[key] = entry
	}
}