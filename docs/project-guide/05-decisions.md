# Decisions

Observed behavior is distinguished below from inferred rationale. The explicit v1 contract is to preserve the binary protocol, exact streaming, whole-file CRC and no-clobber publication ([V1.md:3](../../V1.md#L3)).

## Keep one direct TCP connection and a compact binary header

- **What:** the sender accepts one receiver and sends fixed metadata, a bounded filename and raw body bytes.
- **Evidence:** [single accept](../../internal/sender/sender.go#L90), [codec](../../internal/protocol/binary.go#L24), [body](../../internal/sender/sender.go#L125).
- **Why:** preservation of this protocol is confirmed by [V1.md:5](../../V1.md#L5). Simplicity of deployment and avoiding per-chunk framing are inferred benefits, not recorded original motives.
- **Tradeoff:** no centralized service is required, but there is no discovery, relay, multi-client serving loop or resume. Extending incompatible metadata requires handling the decoder's strict version check ([binary.go:63](../../internal/protocol/binary.go#L63), [documented limits](../../README.md#L131)).
- **Confidence:** behavior and v1 compatibility requirement confirmed; original design rationale inferred.

## Precompute whole-file CRC and require a stable source

- **What:** stat, checksum the source, rewind, then send exactly the saved size. Receiver rereads all received bytes to check the same IEEE CRC.
- **Evidence:** [source prepass](../../internal/sender/sender.go#L46), [bounded sending](../../internal/sender/sender.go#L133), [receiver verification](../../internal/receiver/receiver.go#L104), [checksum algorithm](../../internal/utils/crc.go#L9).
- **Why:** whole-file CRC and exact streaming are confirmed requirements ([V1.md:5](../../V1.md#L5)); obtaining checksum before metadata is a direct consequence of its position in the header ([binary.go:38](../../internal/protocol/binary.go#L38)).
- **Tradeoff:** fixed-size buffers keep memory bounded, but both sides incur additional disk reads and the sender waits through hashing before listening. A changing source can invalidate the metadata; exact-length streaming prevents overrun, not mutation. CRC detects accidental corruption and cannot authenticate a peer ([README.md:131](../../README.md#L131)).
- **Confidence:** confirmed contract and documented source-stability requirement.

## Publish via a same-directory temporary file and hard link

- **What:** accept no existing destination, write and verify a temporary file, sync/close it, link to the requested destination, then remove the temporary name.
- **Evidence:** [early guard](../../internal/receiver/receiver.go#L49), [temp creation](../../internal/receiver/receiver.go#L85), [link](../../internal/receiver/receiver.go#L126), [existing-file test](../../internal/receiver/receiver_test.go#L73).
- **Why:** no-clobber publication is explicitly required ([V1.md:5](../../V1.md#L5)). Using a hard link rather than replacement rename provides a filesystem operation that fails when the destination already exists; the implementation reason is inferred from this behavior.
- **Tradeoff:** readers see a fully written verified destination and existing paths are preserved, but the filesystem must support hard links. Only the file is synced, there is no crash recovery, and failed temporary-name removal can report failure after the destination is already visible ([receiver.go:117](../../internal/receiver/receiver.go#L117)).
- **Confidence:** no-clobber goal confirmed; mechanism rationale inferred.

## Expose cooperative cancellation without changing old call sites

- **What:** `Send` and `Receive` wrap context-bearing methods with background contexts. Callbacks close network resources; receiver checks cancellation before publication.
- **Evidence:** [sender wrappers](../../internal/sender/sender.go#L30), [callbacks](../../internal/sender/sender.go#L82), [receiver wrappers](../../internal/receiver/receiver.go#L37), [publication check](../../internal/receiver/receiver.go#L123).
- **Why:** the compatible API forms and cancellation support are documented ([README.md:129](../../README.md#L129)). Retaining simple CLI call sites is an inferred benefit.
- **Tradeoff:** API callers can interrupt blocked network operations, but the CLI does not wire OS signals into contexts; CRC and filesystem operations are not cancellable. Closed-socket errors are not normalized to `ctx.Err`, and cancellation can race the final link after the explicit check ([main.go:30](../../main.go#L30), [receiver.go:123](../../internal/receiver/receiver.go#L123)).
- **Confidence:** API shape and network cancellation confirmed; full cancellation semantics remain an [open question](README.md#open-questions).

## Keep protocol and checksum independent of sockets

- **What:** decode and checksum accept `io.Reader`; transfer orchestration uses concrete files and `net.Conn`.
- **Evidence:** [Decode](../../internal/protocol/binary.go#L45), [CalculateCRC](../../internal/utils/crc.go#L11), [reader tests](../../internal/utils/crc_test.go#L31), [prepared sender seam](../../internal/sender/sender.go#L78).
- **Why, inferred:** small boundaries allow deterministic unit tests without requiring a network, while prepared-listener helpers enable real TCP integration on ephemeral ports ([header tests](../../internal/protocol/binary_test.go#L10), [integration](../../internal/sender/sender_test.go#L70)).
- **Tradeoff:** codec/CRC logic is easy to isolate, but the main lifecycle still directly invokes filesystem/network APIs. Existing tests do not simulate every publication failure or receiver cancellation path ([receiver test inventory](03-structure.md#receiver)).
- **Confidence:** interfaces and tests confirmed; design rationale inferred.

## Gotchas

- **Sender success is not receiver success.** There is no acknowledgement after checksum or publication; `file sent` can precede receiver failure ([sender.go:148](../../internal/sender/sender.go#L148), [receiver.go:114](../../internal/receiver/receiver.go#L114)).
- **No default port exists.** The README's “default ports can be customized” wording is loose: both command forms require a port, and only the receive address defaults ([README.md:123](../../README.md#L123), [main.go:25](../../main.go#L25)).
- **Cleanup is not a crash guarantee.** The README says interrupted transfers remove partial files automatically, but this is a deferred best-effort operation on normal unwinding. CLI process termination is not converted to cancellation, and removal errors inside the defer are ignored ([README.md:118](../../README.md#L118), [main.go:11](../../main.go#L11), [receiver.go:93](../../internal/receiver/receiver.go#L93)).
- **Output location is caller-controlled.** The README describes saving in the current directory, but a supplied path may name any existing writable directory; the peer filename is only checked and displayed ([README.md:76](../../README.md#L76), [receiver.go:79](../../internal/receiver/receiver.go#L79)). Source permission bits and timestamps are not preserved by the metadata format ([binary.go:15](../../internal/protocol/binary.go#L15)).
- **Header fields are partly derived.** Setting `FileHeader.Protocol`, `Version` or `NameLen` before `Encode` does not override the protocol constants or string byte length ([binary.go:35](../../internal/protocol/binary.go#L35)).
- **An “idle” timeout is not uniform.** Receiver refreshes its deadline before a whole chunk read, so a trickle of bytes does not extend that chunk's 30 seconds. Sender refreshes for each write attempt ([receiver.go:155](../../internal/receiver/receiver.go#L155), [sender.go:153](../../internal/sender/sender.go#L153)).
- **Progress completion is transfer completion only.** Receiving reaches its final progress render before CRC and publication; zero-byte transfers still render 0% because percentage divides only when total bytes is positive ([receiver.go:100](../../internal/receiver/receiver.go#L100), [progress.go:46](../../internal/utils/progress.go#L46), [progress.go:60](../../internal/utils/progress.go#L60)).

## Conventions

- Public role operations return errors; `main` owns logging and exit. Filesystem/network errors usually carry operation-specific `%w` context ([sender.go:40](../../internal/sender/sender.go#L40), [receiver.go:57](../../internal/receiver/receiver.go#L57), [main.go:11](../../main.go#L11)).
- Tests share the implementation package to reach helpers; local socket fixtures use loopback ephemeral ports, `net.Pipe` and `t.TempDir` ([sender tests](../../internal/sender/sender_test.go#L26), [receiver fixtures](../../internal/receiver/receiver_test.go#L15)).
- The receiver header-timeout test temporarily mutates a package global and restores it; making that test parallel would conflict with other tests using the setting ([receiver_test.go:126](../../internal/receiver/receiver_test.go#L126)).
