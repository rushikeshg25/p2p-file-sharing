# Structure

## What lives where

The inspected revision has 21 tracked files. One root Go entry point dispatches into `internal/sender` or `internal/receiver`; `internal/protocol` is their shared wire contract and `internal/utils` supplies checksum, validation and progress. Tests sit beside their implementations, including cross-package sender/receiver integration ([entry](../../main.go#L7), [integration](../../internal/sender/sender_test.go#L46)). No CI, container or deployment configuration is tracked in this revision.

```text
main.go, main_test.go        CLI dispatch and argument checks
go.mod, Makefile             toolchain declaration and local build
README.md, V1.md, HISTORY.md usage, acceptance contract, delivery history
internal/
  sender/                   source preparation and one-client transfer
  receiver/                 download, CRC and safe output publication
  protocol/                 binary header encode/decode
  utils/                    CRC, port validation and terminal progress
bin/.gitkeep                output-directory placeholder
```

## Root

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [main.go](../../main.go#L11) | Select send/receive role, parse positional forms and log failures. | `main`, `run`, `printUsage` (unexported) | Executable runtime; CLI tests |
| [main_test.go](../../main_test.go#L8) | Check both receive forms and validation before source open. | `TestReceiveCLIFormsRemainSupported`, `TestSendRejectsInvalidPortBeforeOpeningFile` | `go test` |
| [go.mod](../../go.mod#L1) | Module name and Go version; no dependency requirements. | Module `p2p-file-sharing` | Go tools |
| [Makefile](../../Makefile#L1) | Build binary, remove binary, set executable bit. | `build`, `clean`, `executable` | Developer/Make |
| [README.md](../../README.md#L1) | Public CLI examples, network requirements and v1 limits. | Documentation | Users and maintainers |
| [V1.md](../../V1.md#L1) | Wire/streaming/publication preservation and acceptance criteria. | Documentation contract | Delivery and review |
| [HISTORY.md](../../HISTORY.md#L1) | Dated delivery evidence and operational validation limits. | Documentation record | Maintainers |

## Sender

Owns source validation and the TCP listening side. `send`/`sendContext` accept a prepared listener and source, allowing tests to bind an ephemeral port instead of racing to discover when `Send` starts listening ([test seam](../../internal/sender/sender.go#L78), [usage](../../internal/sender/sender_test.go#L70)).

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [sender.go](../../internal/sender/sender.go#L15) | Validate/open/stat/CRC/rewind source, listen once, send metadata and exact body, close resources on cancellation. | `Sender`, `NewSender`, `Send`, `SendContext`, `SendFile`, `BUFFER_SIZE`; internal `send`, `sendContext`, `writeAll` | CLI and sender tests |
| [sender_test.go](../../internal/sender/sender_test.go#L15) | Fake partial writes using `net.Pipe`; exercise full localhost transfer and compare output bytes. | `TestWriteAllHandlesShortWrites`, `TestSenderReceiverIntegration`; internal `shortWriteConn` | `go test` |
| [lifecycle_test.go](../../internal/sender/lifecycle_test.go#L12) | Cancel a waiting accept; verify no bytes beyond size and error on shortened source. | `TestCancelledAccept`, `TestSourceLengthBound`, `TestSourceTruncation` | `go test` |

## Receiver

Owns connection initiation, metadata validation and the entire temporary-file-to-destination lifecycle ([receiver.go:39](../../internal/receiver/receiver.go#L39)).

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [receiver.go](../../internal/receiver/receiver.go#L18) | Validate destination, dial, decode, receive exact bytes, CRC check, sync and no-clobber link. | `Receiver`, `NewReceiver`, `Receive`, `ReceiveContext`, `BUFFER_SIZE`; internal `receiveFileData` | CLI, receiver tests, sender integration test |
| [receiver_test.go](../../internal/receiver/receiver_test.go#L15) | Local TCP fixtures and verified publication, no overwrite, failure cleanup and header timeout cases. | `TestReceivePublishesVerifiedFile`, `TestReceiveRefusesExistingDestination`, `TestReceiveCleansUpAfterCRCFailure`, `TestReceiveCleansUpAfterInterruptedTransfer`, `TestReceiveHeaderTimeout`; helper `startTestSender` | `go test` |

## Protocol

Contains a standalone codec taking `io.Reader` rather than a concrete socket ([binary.go:45](../../internal/protocol/binary.go#L45)).

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [binary.go](../../internal/protocol/binary.go#L9) | Define and serialize the version-1, big-endian file header; reject invalid/truncated input. | `FileHeader`, `Encode`, `Decode`, `P2PF_PROTOCOL`, `VERSION`, `HEADER_SIZE` | Sender, receiver and tests |
| [binary_test.go](../../internal/protocol/binary_test.go#L10) | Round trip, derived name length, invalid name, magic/version and truncation cases. | `TestHeaderRoundTrip`, `TestHeaderEncodeRejectsInvalidName`, `TestDecodeRejectsInvalidProtocolAndVersion`, `TestDecodeRejectsTruncatedHeaderAndName` | `go test` |

## Utilities

Shared helpers are synchronous; they do not own transfer goroutines or cancellation ([crc.go:11](../../internal/utils/crc.go#L11), [progress.go:35](../../internal/utils/progress.go#L35)).

| File | Responsibility | Key exports or entry points | Called by |
| --- | --- | --- | --- |
| [crc.go](../../internal/utils/crc.go#L9) | Stream IEEE CRC32 using an 8 KiB buffer and propagate read errors. | `CalculateCRC` | Sender prepass, receiver verification, tests |
| [crc_test.go](../../internal/utils/crc_test.go#L10) | Compare empty/nonempty checksums with standard CRC32 and inject reader error. | `TestCalculateCRC`, `TestCalculateCRCPropagatesReadError`; internal `failingReader` | `go test` |
| [network.go](../../internal/utils/network.go#L8) | Validate decimal port range 1–65535. | `ValidatePort` | Sender/receiver startup, tests |
| [network_test.go](../../internal/utils/network_test.go#L5) | Valid boundaries and malformed/out-of-range ports. | `TestValidatePort` | `go test` |
| [progress.go](../../internal/utils/progress.go#L9) | Terminal bar, byte counts, average rate and ETA with 100 ms throttling. | `ProgressBar`, `NewProgressBar`, `Add`, `Finish`; internal formatting helpers | Sender and receiver chunk loops |

## Excluded

The tracked [bin/.gitkeep](../../bin/.gitkeep) and [.gitignore](../../.gitignore) are an output-directory placeholder and routine ignore configuration. Generated binaries, this guide itself, and the pre-existing untracked `logs/` work artifact are outside the source inventory. There is no tracked vendor tree or dependency lockfile; [go.mod](../../go.mod#L1) lists no external requirements.
