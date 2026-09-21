# P2P File Sharing Project Guide

> Generated: 2026-09-21 from commit `ff57fca`. Scope: the v1 source and tests in this checkout.

## What this is

This Go CLI transfers one regular file directly between two reachable peers over TCP ([sender](../../internal/sender/sender.go#L40)). The sender listens for one receiver, and the receiver chooses a local output path, checks a whole-file CRC32, and publishes the verified file without overwriting an existing destination ([receiver](../../internal/receiver/receiver.go#L85)). It is intended for trusted peers: v1 has no encryption, authentication, discovery, resume, or receiver acknowledgement ([documented limits](../../README.md#L127)).

## Run it

Use Go 1.24.6 or a compatible newer toolchain, plus Make if using the build target ([manifest](../../go.mod#L3), [build](../../Makefile#L3)). From the repository root:

```bash
make build
# Terminal 1: choose an existing source file and a reachable port.
./bin/p2p-share send source.bin 3001
# Terminal 2 on the same machine; the output path must not exist.
./bin/p2p-share receive received.bin 3001
# For a different machine, supply the sender address.
./bin/p2p-share receive 192.168.1.42 received.bin 3001
# Behavioral checks, including local TCP integration.
go test -race ./...
```

The sender must finish its checksum pass before it begins listening. Ports are explicit decimal values in 1–65535; the receiver defaults only its address to `localhost`, not its port ([startup](02-flow.md#startup)). The destination parent directory must already exist and permit temporary-file creation and hard links; no environment variables are read by the application ([publication](../../internal/receiver/receiver.go#L85), [CLI configuration](../../main.go#L18)). Keep the source unchanged during checksum and transfer ([source constraint](../../README.md#L131)).

## The five-file tour

| # | File | Why this one | Then look at |
| --- | --- | --- | --- |
| 1 | [main.go](../../main.go#L18) | CLI forms and dispatch. | [Startup](02-flow.md#startup) |
| 2 | [sender.go](../../internal/sender/sender.go#L32) | Source preparation, listen, header, exact body writes. | [Transfer](02-flow.md#file-transfer) |
| 3 | [binary.go](../../internal/protocol/binary.go#L15) | The complete on-wire contract. | [Boundaries](01-architecture.md#boundaries-and-contracts) |
| 4 | [receiver.go](../../internal/receiver/receiver.go#L39) | Download, verification, safe publication. | [Publication](02-flow.md#verification-and-publication) |
| 5 | [sender_test.go](../../internal/sender/sender_test.go#L46) | A real localhost transfer tying the pieces together. | [Tests and build](02-flow.md#build-and-tests) |

## Reading order for this guide

1. [Architecture](01-architecture.md) — components, wire format, persistence and limits.
2. [Flow](02-flow.md) — startup, transfer, publication, cancellation and build.
3. [Structure](03-structure.md) — complete source and test inventory.
4. [Tech stack](04-tech-stack.md) — Go, standard library and local tooling.
5. [Decisions](05-decisions.md) — confirmed contracts, inferred rationale and gotchas.

## Open questions

- Should CLI signals trigger cooperative cleanup? `main` calls background-context APIs and installs no signal handler; the README's interrupted-transfer cleanup statement therefore does not cover process termination ([entry](../../main.go#L11), [claim](../../README.md#L118), [cleanup](../../internal/receiver/receiver.go#L93)).
- What cancellation guarantee is intended for checksum/filesystem work and the instant of publication? Cancellation closes sockets, but CRC is not context-aware and the final context check occurs before `os.Link` ([CRC](../../internal/utils/crc.go#L11), [publication](../../internal/receiver/receiver.go#L123)).
- Which filesystems/platforms must support publication, and is crash durability required? The implementation uses a same-directory temporary file, file sync and hard linking; it has no directory sync or crash recovery ([publication](../../internal/receiver/receiver.go#L85)).
- Should coverage expand to receiver cancellation, publication races and payload idle timeouts? Current lifecycle coverage directly tests cancelled sender accept, and the receiver timeout test covers only the header ([lifecycle tests](../../internal/sender/lifecycle_test.go#L12), [receiver tests](../../internal/receiver/receiver_test.go#L126)). No production operational validation is recorded ([history](../../HISTORY.md#L50)).
