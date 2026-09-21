# Direct TCP file transfer history

> Scope: repository baseline through the v1 delivery branch.
> Last updated: 2026-09-21. Dates below retain Git author timezone offsets.

## At a glance

The v1 scope and acceptance contract is recorded in [V1.md](V1.md). Current behavior and limitations are documented in [README.md](README.md).

## Timeline

### 2025-09-06T17:57:29+05:30: feat: init

- **What happened:** The repository records `feat: init`.
- **Evidence:** [commit 2013f85fa9](https://github.com/rushikeshg25/p2p-file-sharing/commit/2013f85fa9db0e3785f91633d27a428e7adbacf2).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:21:37+05:30: docs: define p2p-file-sharing v1 contract

- **What happened:** The repository records `docs: define p2p-file-sharing v1 contract`.
- **Evidence:** [commit 93a89195ef](https://github.com/rushikeshg25/p2p-file-sharing/commit/93a89195efdbc8b69b3733b743f373cdc54e00d0).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:34:06+05:30: feat: add cancellable transfers and enforce advertised source length

- **What happened:** The repository records `feat: add cancellable transfers and enforce advertised source length`.
- **Evidence:** [commit 0ab9f19c00](https://github.com/rushikeshg25/p2p-file-sharing/commit/0ab9f19c00b3def50a36cc28a6ad6225d4c2fce2).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:34:06+05:30: test: verify cancellation source truncation and exact transfer bounds

- **What happened:** The repository records `test: verify cancellation source truncation and exact transfer bounds`.
- **Evidence:** [commit 26b23772b7](https://github.com/rushikeshg25/p2p-file-sharing/commit/26b23772b787fa7b88303bf5a48a8f9365908d98).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

### 2026-09-21T13:38:57+05:30: docs: document p2p-file-sharing v1 usage and limitations

- **What happened:** The repository records `docs: document p2p-file-sharing v1 usage and limitations`.
- **Evidence:** [commit cb8dcbf1e0](https://github.com/rushikeshg25/p2p-file-sharing/commit/cb8dcbf1e000a89b15db9bd1e3ca4b367372f36a).
- **Confidence:** Confirmed by Git history; the title alone does not establish runtime correctness.

## Delivery verification

Passed `go test -race ./...` including full localhost sender/receiver transfers and [lifecycle tests](internal/sender/lifecycle_test.go).

## Turning points

The v1 contract made failure behavior, lifecycle semantics and executable verification part of the delivery. The new tests and README describe the resulting boundaries.

## Open questions

No production deployment or long-running operational validation was performed as part of this delivery.
