# Upstream synchronization

This branch was rebuilt from OpenList `v4.2.5` (`cc87e88f`) instead of
continuing the previous conflict-heavy merge. The official source is tracked
as the `upstream` remote:

```text
https://github.com/OpenListTeam/OpenList.git
```

## Maintained iNoi layer

The backend should stay as close to upstream as possible. iNoi-specific
changes are intentionally limited to:

- the iNoi executable, CLI text, frontend package and default branding;
- compatibility with iNoi's password static salt, plus recovery support for
  accounts created while the upstream OpenList/AList salt was in use;
- iNoi container paths, image names and release metadata;
- a per-storage `web`/`android` protocol selector for the 123Pan driver. Web
  remains the default and uses the official v4.2.5 endpoints and signing;
  Android only replaces the endpoint/header/login compatibility layer and
  continues to use the upstream list, upload, offline-task and response code.

The 123Pan extension lives in `drivers/123/protocol.go`. Never replace the
whole driver with the pre-sync implementation or restore the old
process-wide `PAN123_PROTOCOL` switch. Existing and invalid protocol values
must fall back to Web so an upgrade cannot silently change request behavior.

Small format-string compatibility changes across affected v4.2.5 drivers and
offline-download adapters keep the baseline compatible with newer Go vet
checks. They only make remote text an explicit string argument and do not
change request behavior; they can be dropped after upstream carries
equivalent fixes.

Do not restore old copies of `internal/op`, `internal/stream`,
`internal/cache`, `internal/net`, `server/handles`, `server/middlewares`, or
whole driver files from the pre-sync branch. Those differences were mostly
stale upstream snapshots and caused cache corruption, broken SSO protocol
fields, missing routes and upload regressions.

## Future updates

1. Fetch `upstream` and choose a released OpenList tag.
2. Merge or rebase that tag while keeping the maintained layer above small.
3. Review changes to `build.sh`, `internal/bootstrap/data/setting.go`,
   `internal/model/user.go`, `server/handles/auth.go`, `drivers/123/protocol.go`,
   and container files.
4. Run the upstream test/build matrix and iNoi compatibility tests before
   publishing an image.

## Acceptance on 2026-08-27

An independent `gpt-5.6-luna` acceptance run used the repository's declared
Go 1.26.4 toolchain on Windows. The following checks passed:

- `gofmt` and `git diff --check v4.2.5 --`;
- `go test ./drivers/123 ./internal/model ./server/handles`;
- `go build ./...`;
- shell syntax for every tracked shell script and YAML parsing for all ten
  workflows;
- the downloaded iNoi frontend archive layout and both extraction branches in
  `FetchINoiWebPackage`.

`go test ./... -timeout 10m` completed with failures only in unchanged
upstream/environment-dependent packages: `internal/multipart` (Windows file
locking), `internal/net` (Go 1.26 transport wrapper), `pkg/aria2/rpc` (no
local aria2 service), and `server/mcp` (session-limit assertion). No failure
was attributed to the iNoi customization. `go vet ./...` has the unchanged
Windows `internal/mem` unsafe-pointer diagnostic. Race testing requires CGO,
and no Docker-compatible runtime was available on the acceptance host.
