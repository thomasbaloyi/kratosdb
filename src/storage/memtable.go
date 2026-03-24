package storage

type MemTable struct {
	table map[string]entry // TODO: Add concurrency protection + use a sorted map
}

type entry struct {
	Value 	string
	Deleted bool
}

func NewMemTable() *MemTable {
	dict = make(map[string]entry)
	return &MemTable{table: dict}
}

func (m *MemTable) Put(key string, value string) {
	m.table[key] = entry{Value: value, Deleted: false}
}

func (m *MemTable) Get(key string) (string, bool) {
	entry, exists := m.table[key]
	if exists && !entry.Deleted {
		return entry.Value, exists
	}
	 
	return "", false
}

func (m *MemTable) Delete(key string) {
	_, exists := m.table[key] 
	if exists {
		m.table[key].Deleted = true // Entry.Deleted
	}
}