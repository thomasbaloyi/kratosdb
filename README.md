# KratosDB

KratosDB is a storage engine built from scratch in Go, following *Database Internals* by Alex Petrov chapter by chapter. It uses an **LSM-tree** (Log-Structured Merge-tree) architecture, the same foundation as RocksDB and LevelDB. The long-term target is to match or exceed RocksDB's write-read throughput (benchmarks to be defined in a future milestone).

---

## Current State

The project is in its earliest stage. The in-memory component of the LSM-tree write path is being built first.

**Implemented:**
- `MemTable` — an in-memory key/value store with soft-delete (tombstone) semantics via an `entry` struct

**Known gaps (in progress):**
- `NewMemTable` uses an undeclared variable (`dict`) — does not compile yet
- `Delete` assigns directly to a map field — does not compile yet
- No concurrency protection (`sync.RWMutex`)
- Underlying map is unordered — needs to be replaced with a sorted structure (e.g. a skip list or red-black tree) for efficient SSTable flushing
- No tests yet

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

#### Milestone 1 — MemTable (Chapter 3: File Formats / Chapter 6: LSM Trees)
The write buffer that absorbs all incoming writes before they are flushed to disk.

- [ ] Fix `NewMemTable` and `Delete` compile errors
- [ ] Add `sync.RWMutex` for concurrent access
- [ ] Replace `map` with a sorted in-memory structure (skip list)
- [ ] Define a size threshold that triggers a flush
- [ ] Write table-driven unit tests

#### Milestone 2 — Write-Ahead Log / WAL (Chapter 3: File Formats)
Durability guarantee — every write is appended to the WAL before being applied to the MemTable, so writes survive crashes.

- [ ] Binary log format with fixed-size record framing
- [ ] Sequential append-only writes
- [ ] WAL replay on startup to recover an unflushed MemTable
- [ ] Log segment rotation

#### Milestone 3 — SSTable (Chapter 3: File Formats, Chapter 6: LSM Trees)
The on-disk, immutable, sorted representation of a flushed MemTable.

- [ ] Sort and serialise a MemTable to disk on flush
- [ ] SSTable file format: data blocks + index block + footer
- [ ] Point lookup via binary search on the index
- [ ] Iterator for sequential scans
- [ ] Handle tombstones during reads

#### Milestone 4 — Bloom Filters (Chapter 7: Log-Structured Storage)
Probabilistic filter attached to each SSTable to short-circuit disk reads for keys that do not exist.

- [ ] Implement a basic bloom filter
- [ ] Attach a per-SSTable bloom filter at flush time
- [ ] Skip SSTable reads for keys that fail the filter

#### Milestone 5 — LSM-Tree Read Path (Chapter 6 & 7)
Assemble the full read path: MemTable → immutable MemTables → SSTable levels.

- [ ] Multi-level read path with correct recency ordering
- [ ] Tombstone propagation across levels
- [ ] Iterator merge across levels (merge iterator)

#### Milestone 6 — Compaction (Chapter 7: Log-Structured Storage)
Background process that merges and garbage-collects SSTables to bound read amplification and space amplification.

- [ ] Size-tiered compaction strategy (initial)
- [ ] Leveled compaction strategy (RocksDB-style)
- [ ] Tombstone resolution and key deduplication during compaction
- [ ] Compaction scheduling and concurrency

#### Milestone 7 — Block Cache & Read Optimisations (Chapter 5: Transactions and Recovery)
- [ ] LRU block cache for hot SSTable blocks
- [ ] Key prefix compression within data blocks
- [ ] Configurable block size

#### Milestone 8 — Recovery & Crash Consistency (Chapter 5)
- [ ] Full WAL-based crash recovery
- [ ] Manifest file tracking live SSTable versions
- [ ] Atomic version transitions (MVCC groundwork)

---

## References

- *Database Internals* — Alex Petrov (O'Reilly, 2019)
- [RocksDB Architecture](https://github.com/facebook/rocksdb/wiki/RocksDB-Overview)
- [LevelDB Implementation Notes](https://github.com/google/leveldb/blob/main/doc/impl.md)
