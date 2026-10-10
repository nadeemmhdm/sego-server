# Changelog

SEGO Server is pre-1.0. This changelog records meaningful development milestones; it does not imply a signed production release unless a GitHub release explicitly says so.

## Unreleased

### Added

- direct one-command GitHub installers for Windows and Linux
- automatic performance-mode selection
- simplified `sego open/start/stop/restart/status/help` commands
- automatic PATH/launcher setup
- verified Go toolchain fallback on Linux when the system Go is too old
- authorized image/PDF/text/audio/video previews
- browser video/audio playback with Range request support
- temporary scoped preview URLs
- admin Emergency SOS global lockdown
- session, share-link and preview-token revocation during SOS
- local-only emergency unlock flow
- administrator audit log list/detail UI and API
- audit events for authentication, file actions, sharing, previews and emergency controls

### Changed

- aligned CI/toolchain minimums with dependency requirements
- improved beginner setup guidance and installation documentation
- improved preview/download error handling in the web UI

### Security

- tokenized preview/share routes are logged by route pattern rather than full opaque-token URL
- active HTML/SVG is not allowed as inline preview content
- Emergency SOS does not restore previously revoked share links after recovery
