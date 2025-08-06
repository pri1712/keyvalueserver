# Distributed Lock with Fault-Tolerant KV Store

This project implements a **distributed lock system** using a custom-built key-value store with support for **unreliable network conditions**.

## Features

- Distributed mutual exclusion lock (`Lock`, `Unlock`) built on top of KVStore
- Versioned writes to ensure linearizability (verified using porcupine linearizability checker)
- Supports:
    - Reliable networks
    - Unreliable message delivery (loss, delay, reordering)
- Client-side retries with idempotency
- Safe concurrent access via versioning
- Lock release uses compare-and-swap semantics
- Extensible to support per-client locks


## Unreliable Network Handling

- KVServer simulates an unreliable network using random drops and reordering
- client implements retry logic for robustness
- Lock operations tested under:
