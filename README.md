# Distributed Key-Value Store

This project implements a simple distributed key-value server inspired by MIT 6.824 labs. It provides a reliable RPC-based key-value storage service with versioned `Put` and `Get` operations, supporting both single-client and concurrent access scenarios.

## Features

- Thread-safe key-value store using mutex locking.
- Version-controlled `Put` operations to support idempotency and consistency.
- RPC-based interface using Go’s `net/rpc`.
- Minimal client retry logic (for reliable network setup).
- Fully passing correctness and race-condition tests under reliable network assumptions.