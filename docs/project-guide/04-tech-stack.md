# Tech Stack

## Languages and runtimes

| Language or runtime | Version | Pinned at |
| --- | --- | --- |
| Go | `1.24.6` language/minimum toolchain directive; no separate toolchain pin | [go.mod:3](../../go.mod#L3) |
| Go module | `p2p-file-sharing`; no module release version declared here | [go.mod:1](../../go.mod#L1) |

The module contains no `require` directives. All application imports are Go standard-library packages or this module's internal packages ([go.mod](../../go.mod#L1), [sender imports](../../internal/sender/sender.go#L3), [receiver imports](../../internal/receiver/receiver.go#L3)). Standard-library versions follow the actual Go toolchain used; no independent library pins exist.

## Frameworks and major libraries

| Library | Version | Used for | Used in |
| --- | --- | --- | --- |
| `net`, `context`, `time` | Bundled with Go; [version directive](../../go.mod#L3) | TCP listen/dial, socket deadlines and cancellation callbacks. | [sender.go:69](../../internal/sender/sender.go#L69), [receiver.go:57](../../internal/receiver/receiver.go#L57) |
| `encoding/binary`, `io` | Bundled with Go; [version directive](../../go.mod#L3) | Big-endian framing and exact header/body reads. | [binary.go:24](../../internal/protocol/binary.go#L24), [receiver.go:158](../../internal/receiver/receiver.go#L158) |
| `hash/crc32`, `bufio` | Bundled with Go; [version directive](../../go.mod#L3) | Whole-file IEEE CRC32 in bounded memory. | [crc.go:9](../../internal/utils/crc.go#L9) |
| `os`, `path/filepath`, `io/fs` | Bundled with Go; [version directive](../../go.mod#L3) | File validation, temporary creation, sync, no-clobber hard link and cleanup. | [receiver.go:49](../../internal/receiver/receiver.go#L49), [receiver.go:85](../../internal/receiver/receiver.go#L85) |
| `fmt`, `log`, `strings`, `time` | Bundled with Go; [version directive](../../go.mod#L3) | CLI messages, error log and hand-built ANSI progress display. | [main.go:11](../../main.go#L11), [progress.go:56](../../internal/utils/progress.go#L56) |
| `testing` | Bundled with Go; [version directive](../../go.mod#L3) | Unit cases, temporary directories, TCP and pipe-based behavioral tests. | [sender_test.go:26](../../internal/sender/sender_test.go#L26), [receiver_test.go:48](../../internal/receiver/receiver_test.go#L48) |

There is no CLI framework: argument parsing is a direct switch over positional arguments ([main.go:18](../../main.go#L18)).

## Data and infrastructure

| Service or resource | Role | Configured at |
| --- | --- | --- |
| Direct TCP socket | One sender-to-receiver transfer on a required user-selected port; no TLS layer. | [sender.go:69](../../internal/sender/sender.go#L69), [receiver.go:55](../../internal/receiver/receiver.go#L55) |
| Sender filesystem | Existing regular source file; full CRC pass before listening. | [sender.go:40](../../internal/sender/sender.go#L40) |
| Receiver filesystem | Same-directory temporary file and hard-linked final destination after CRC. | [receiver.go:85](../../internal/receiver/receiver.go#L85), [receiver.go:117](../../internal/receiver/receiver.go#L117) |
| No datastore/service dependency | All transfer state is local process memory and filesystem state. | [sender lifecycle](../../internal/sender/sender.go#L32), [receiver lifecycle](../../internal/receiver/receiver.go#L39) |

No container, orchestration, cloud configuration or CI workflow is present in the [source inventory](03-structure.md#what-lives-where). External network/firewall reachability remains an operator responsibility ([README.md:93](../../README.md#L93)).

## Tooling

| Tool | Role | Configured at |
| --- | --- | --- |
| Make | `build` creates `bin/p2p-share`; `clean` removes it; `executable` applies `chmod +x`. Make version is not pinned. | [Makefile:1](../../Makefile#L1) |
| `go build` | Compile the root executable for the active platform; no cross-build matrix. | [Makefile:5](../../Makefile#L5) |
| `go test -race ./...` | Documented v1 validation command, using standard tests and the race detector. | [README.md:133](../../README.md#L133), [HISTORY.md:42](../../HISTORY.md#L42), [go.mod:3](../../go.mod#L3) |
| Git ignore rules | Keep binaries, test/coverage output and local workspace files out of tracking. | [.gitignore:4](../../.gitignore#L4) |

## Notes

- Constants define protocol version 1, 2 KiB streaming buffers and 8 KiB checksum buffer; they are application parameters, not dependency versions ([binary.go:9](../../internal/protocol/binary.go#L9), [sender.go:20](../../internal/sender/sender.go#L20), [receiver.go:32](../../internal/receiver/receiver.go#L32), [crc.go:14](../../internal/utils/crc.go#L14)).
- Deadline settings are code-level values, not CLI flags: sender write 30 seconds, receiver header/body 30 seconds, receiver dial 10 seconds ([sender.go:21](../../internal/sender/sender.go#L21), [receiver.go:34](../../internal/receiver/receiver.go#L34), [receiver.go:57](../../internal/receiver/receiver.go#L57)).
- Progress uses ANSI line clearing unconditionally; redirected output will contain control sequences, and no terminal library or TTY check is used ([progress.go:85](../../internal/utils/progress.go#L85)).
