package memtable

import "testing"

func TestPut_Get(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"simple key", "foo", "bar"},
		{"empty value", "key", ""},
		{"empty key", "", "value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemTable()
			m.Put(tt.key, tt.value)

			got, ok := m.Get(tt.key)
			if !ok {
				t.Fatalf("Get(%q): expected found=true, got false", tt.key)
			}
			if got != tt.value {
				t.Fatalf("Get(%q): expected %q, got %q", tt.key, tt.value, got)
			}
		})
	}
}

func TestGet_MissingKey(t *testing.T) {
	m := NewMemTable()

	got, ok := m.Get("missing")
	if ok {
		t.Fatal("Get on missing key: expected found=false, got true")
	}
	if got != "" {
		t.Fatalf("Get on missing key: expected empty string, got %q", got)
	}
}

func TestPut_Overwrite(t *testing.T) {
	m := NewMemTable()
	m.Put("key", "first")
	m.Put("key", "second")

	got, ok := m.Get("key")
	if !ok {
		t.Fatal("Get after overwrite: expected found=true, got false")
	}
	if got != "second" {
		t.Fatalf("Get after overwrite: expected %q, got %q", "second", got)
	}
}

func TestDelete_HidesKey(t *testing.T) {
	m := NewMemTable()
	m.Put("key", "value")
	m.Delete("key")

	got, ok := m.Get("key")
	if ok {
		t.Fatal("Get after Delete: expected found=false, got true")
	}
	if got != "" {
		t.Fatalf("Get after Delete: expected empty string, got %q", got)
	}
}

func TestDelete_NonExistentKey(t *testing.T) {
	m := NewMemTable()

	// should not panic
	m.Delete("ghost")

	_, ok := m.Get("ghost")
	if ok {
		t.Fatal("Get after Delete of non-existent key: expected found=false, got true")
	}
}

func TestPut_AfterDelete_RestoresKey(t *testing.T) {
	m := NewMemTable()
	m.Put("key", "original")
	m.Delete("key")
	m.Put("key", "restored")

	got, ok := m.Get("key")
	if !ok {
		t.Fatal("Get after Put-Delete-Put: expected found=true, got false")
	}
	if got != "restored" {
		t.Fatalf("Get after Put-Delete-Put: expected %q, got %q", "restored", got)
	}
}
