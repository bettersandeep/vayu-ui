# Vayu fork

Vayu-ui is a customized fork of [datazip-inc/olake-ui](https://github.com/datazip-inc/olake-ui).
This document tracks what the fork adds on top of upstream.

**Upstream base:** olake-ui `v0.4.9`

## Custom features

- **Replication-slot lifecycle in the UI** - drop the slot on job and source deletion, guarded so a
  slot shared by another **active** job blocks deletion, while inactive jobs proceed with a warning.
  Slot-deletion checkboxes and validation warnings are surfaced in the frontend.
- **Fusion compaction scoping** - a **Databases** selector on the maintenance catalog form maps to
  the optimizer's `database-filter` catalog property (round-tripped on edit and preserved across
  updates), so Fusion only explores/compacts the chosen databases in a shared catalog.
- **Sentry integration** - server (Gin middleware) and frontend (React init + error boundary).
- **Session resilience** - re-verify with the server before logging out on a transient 401.
- **Telemetry** - disable external analytics calls.
- **catalog_name** - allow digits in the pattern.
- **Destination keys** - persist the destination-key column selection across stream-config edits.

## Versioning

- Release tags: `vayu-<X.Y.Z>+olake<upstream-base>` (e.g. `vayu-0.1.0+olake0.4.9`).
- Image tags map `+` → `-`.
- The cross-component set that deploys together is pinned in `vayu-release.yaml` (vayu-helm repo),
  tagged `vayu-platform-<X.Y.Z>`.

## Branch model & contributing

- `upstream` / `master` - pristine upstream mirror; merge base for syncs. Never commit here.
- `vayu-main` - protected integration/release branch. **Direct pushes are blocked**; changes land
  via pull requests.
- `feature/*` - branch off `vayu-main`, PR back in, updating `VAYU.md` and `CHANGELOG.md`.
- Sync upstream by merging `upstream` into `vayu-main` (never rebase the release branch).
