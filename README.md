# KratosDB

KratosDB is a storage engine built from scratch in Go, following *Database Internals* by Alex Petrov chapter by chapter. It uses an **LSM-tree** (Log-Structured Merge-tree) architecture, the same foundation as RocksDB and LevelDB. The long-term target is to match or exceed RocksDB's write-read throughput (benchmarks to be defined in a future milestone).

---

## Current State

Currently reading Chapter 4. The MemTable is implemented ahead of its chapter (Chapter 7) as a project scaffold — it will be fully contextualised when Chapter 7 is reached.

**Implemented:**
- `memtable.Table` — in-memory key/value store with `Put`, `Get`, and soft-delete (`Delete`) via tombstone semantics
- Key-ordered storage via `github.com/egregors/sortedmap` — keys maintained in lexicographic order
- `sync.RWMutex` concurrency protection — safe for concurrent goroutine access
- Table-driven unit tests covering: basic put/get, missing keys, overwrites, delete, delete of non-existent key, put-after-delete, and concurrent access
- All tests pass under `go test -race`

**Project layout (target):**
```
kratosdb/
  go.mod
  internal/
    memtable/      # in-memory sorted write buffer
    wal/           # write-ahead log
    sstable/       # on-disk sorted string tables
    bloom/         # bloom filters for read optimisation
    compaction/    # compaction strategies
  storage/         # public-facing DB API
  cmd/kratosdb/    # binary entry point
```

---

## Milestones

Progress is tracked chapter by chapter through *Database Internals*.

### Part I — Storage Engines

#### Milestone 1 — MemTable 
The in-memory write buffer. Implemented early to have something to build on — the book formally introduces this in Chapter 7.

- [x] `Put`, `Get`, `Delete` with tombstone semantics
- [x] Correct project layout (`internal/memtable`)
- [x] Table-driven unit tests — all passing under `go test -race`
- [x] `sync.RWMutex` for concurrent access
- [x] Key-ordered storage (`egregors/sortedmap`)

#### Milestone 2 — File Format Primitives (Chapter 3: File Formats)
Binary encoding of records on disk. Foundation for everything that writes to or reads from disk (SSTable, WAL).

- [ ] Fixed-size record encoding/decoding
- [ ] Variable-size record encoding with length prefixes
- [ ] Key-value cell serialisation

#### Milestone 3 — Write-Ahead Log / WAL (Chapter 5: Transaction Processing and Recovery)
Durability guarantee — every write is appended to the WAL before being applied to the MemTable, so writes survive crashes.

- [ ] Binary log format using Chapter 3 record framing
- [ ] Sequential append-only writes
- [ ] WAL replay on startup to recover an unflushed MemTable
- [ ] Log segment rotation

#### Milestone 4 — SSTable (Chapter 7: Log-Structured Storage)
The on-disk, immutable, sorted representation of a flushed MemTable.

- [ ] Serialise a MemTable to disk using Chapter 3 file format primitives
- [ ] SSTable file format: data blocks + index block + footer
- [ ] Point lookup via binary search on the index
- [ ] Iterator for sequential scans
- [ ] Handle tombstones during reads

#### Milestone 5 — Full LSM-Tree (Chapter 7: Log-Structured Storage)
Wire up the complete LSM-tree: flush threshold, immutable MemTable handoff, and the full read path.

- [ ] Size threshold that triggers a flush
- [ ] Immutable MemTable promotion and background flush
- [ ] Multi-level read path: MemTable → immutable MemTables → SSTables
- [ ] Tombstone propagation across levels
- [ ] Merge iterator across levels

#### Milestone 6 — Bloom Filters (Chapter 7: Log-Structured Storage)
Probabilistic filter per SSTable to short-circuit disk reads for absent keys.

- [ ] Implement a basic bloom filter
- [ ] Attach a per-SSTable bloom filter at flush time
- [ ] Skip SSTable reads for keys that fail the filter

#### Milestone 7 — Compaction (Chapter 7: Log-Structured Storage)
Background process that merges SSTables to bound read amplification and reclaim space from tombstones.

- [ ] Size-tiered compaction strategy
- [ ] Leveled compaction strategy (RocksDB-style)
- [ ] Tombstone resolution and key deduplication during compaction
- [ ] Compaction scheduling and concurrency

#### Milestone 8 — Recovery & Crash Consistency (Chapter 5: Transaction Processing and Recovery)
- [ ] Full WAL-based crash recovery
- [ ] Manifest file tracking live SSTable versions
- [ ] Atomic version transitions

---

## References

- *Database Internals* — Alex Petrov (O'Reilly, 2019)
- [RocksDB Architecture](https://github.com/facebook/rocksdb/wiki/RocksDB-Overview)
- [LevelDB Implementation Notes](https://github.com/google/leveldb/blob/main/doc/impl.md)
