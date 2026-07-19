# Changelog

All notable vayu-fork changes are documented here. Format based on
[Keep a Changelog](https://keepachangelog.com/); versions follow `vayu-<X.Y.Z>+olake<base>`.

## [vayu-0.1.0+olake0.4.9] - 2026-07-19

First tagged release of the vayu-ui fork (based on upstream olake-ui `v0.4.9`).

### Added
- Replication-slot drop on job/source deletion with a shared-slot guard (active jobs block; inactive warn).
- Slot-deletion checkboxes and validation-warning surfacing in the frontend.
- Fusion catalog **Databases** scoping via the `database-filter` property (with edit round-trip and retention).
- Sentry integration: server Gin middleware and frontend React init + error boundary.

### Changed
- Re-verify session with the server before logout on a transient 401.
- `catalog_name` pattern allows digits.
- Persist destination-key column selection across stream-config edits.

### Removed
- External analytics telemetry calls.
