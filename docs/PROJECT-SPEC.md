# PROJECT: SEGO SERVER

## Open-Source Self-Hosted Private Cloud Storage Platform

Act as a senior full-stack developer, cross-platform systems engineer, cybersecurity engineer, DevOps engineer, performance engineer, and UI/UX designer.

Build a complete, production-oriented, open-source storage server named **SEGO Server**.

SEGO Server must turn an old laptop or desktop PC into a private cloud storage server. Users must be able to upload, organize, preview, download and securely share files and folders. Administrators must be able to manage users, storage quotas, security, updates and remote access.

The platform must support Windows and Linux through one shared application architecture, with OS-specific installation and service management.

The project must include a working backend, database, frontend, CLI, installers, tests, documentation, update manager and security controls. Do not deliver only a visual prototype.

---

# 1. CORE REQUIREMENTS

Build the following modules:

1. Cross-platform installation and uninstallation.
2. Hardware and disk discovery.
3. Admin account creation and secure authentication.
4. User, group and permission management.
5. File and folder management.
6. Secure file sharing.
7. Secure folder sharing.
8. File previews and metadata.
9. Search, filtering and sorting.
10. Storage allocation and quota enforcement.
11. Cloudflare Tunnel integration for worldwide access.
12. macOS-inspired animated dashboard.
13. Resource-efficient background services.
14. Application updates and compatible package maintenance.
15. Backup, recovery and file versioning.
16. Audit logs and security monitoring.
17. CLI tools and diagnostics.
18. Automated testing and release automation.

Use the name SEGO Server consistently across the CLI, dashboard, installer, logs and documentation.

---

# 2. TECHNOLOGY ARCHITECTURE

Select a lightweight, maintainable, cross-platform technology stack.

Suggested stack:

Frontend:

* React
* TypeScript
* Vite
* Tailwind CSS
* A lightweight accessible component library
* CSS transitions and restrained animations

Backend:

* Go for the main API server and CLI, or justify an equivalent lightweight cross-platform stack.
* SQLite for the initial single-server deployment.
* A database migration system.
* A structured logging library.

Storage:

* Native local filesystem.
* Separate application data, database, configuration and user-data directories.
* Support multiple explicitly configured storage roots.

Deployment:

* Native Windows service and Linux systemd service.
* Optional Docker deployment for users who prefer containers.
* Cloudflare Tunnel integration through the official cloudflared client.

Avoid unnecessary services, heavy runtime dependencies, multiple always-running processes and excessive JavaScript bundles.

Keep the storage engine and permission checks independent of the frontend.

---

# 3. HARD PERFORMANCE LIMIT: 4 GB RAM

The platform must be designed to run smoothly on older computers with approximately 4 GB of physical RAM.

Treat 4 GB as a total system-memory budget, not a dedicated application allocation.

Performance targets:

* Aim for approximately 300–600 MB idle RAM usage for the SEGO application services on a lean supported system, where feasible.
* Avoid exceeding 1 GB for SEGO application services during ordinary file operations on the reference system.
* Leave memory for the operating system, filesystem cache, browser and other essential services.
* Do not claim these targets are met until measured on real hardware.
* Optimize for older dual-core CPUs and mechanical HDDs, not just modern SSDs.
* Prefer bounded concurrency and streaming over loading entire files into memory.

Required optimizations:

* Stream uploads and downloads.
* Never load entire large files into RAM for previews, downloads or checksum calculations.
* Use pagination and virtualized file lists for large directories.
* Use lazy loading for thumbnails and previews.
* Generate thumbnails asynchronously with strict CPU and memory limits.
* Limit concurrent uploads, downloads and preview generation.
* Use bounded worker pools and queues.
* Configure request timeouts and maximum upload sizes.
* Use database indexes and efficient queries.
* Keep database connection counts small.
* Avoid frequent polling when event-driven updates or slower polling are sufficient.
* Cache only necessary metadata and enforce cache limits.
* Debounce search and sorting controls.
* Avoid loading entire directory trees at once.
* Respect browser reduced-motion preferences.
* Disable nonessential animations on low-power devices.
* Use efficient, small frontend bundles.
* Avoid unnecessary telemetry and external API requests.
* Use a configurable background-job concurrency limit.

Provide a Low Resource Mode that:

* Reduces simultaneous background tasks.
* Limits thumbnail generation.
* Uses conservative upload/download concurrency.
* Reduces dashboard refresh frequency.
* Disables nonessential visual effects.
* Provides a simple, readable file-list view.

Create a resource diagnostics page showing actual CPU, memory, disk usage, and service status when available.

Test resource consumption under idle conditions, concurrent file transfers, large directories, preview generation and background updates.

---

# 4. CROSS-PLATFORM INSTALLER

Create an installer for supported Windows and Linux systems.

Linux:

* Initially support tested Ubuntu Server and Debian releases.
* Detect architecture and distribution.
* Create a dedicated service user.
* Use systemd where available.

Windows:

* Support explicitly tested Windows 10/11 configurations.
* Provide a PowerShell bootstrap and a signed native executable.
* Register a Windows service with a least-privileged identity.
* Use Windows-specific filesystem and service APIs where necessary.

Installer behavior:

1. Detect operating system and architecture.
2. Check minimum requirements and available disk space.
3. Detect existing installations.
4. Detect missing dependencies and incompatible versions.
5. Show planned installation and package changes.
6. Install missing dependencies and supported updates with user approval.
7. Download and verify a signed SEGO release.
8. Create secure directories and service configuration.
9. Collect administrator setup information.
10. Configure storage roots.
11. Configure the background service.
12. Optionally configure Cloudflare Tunnel.
13. Run database migrations.
14. Start the service and execute health checks.
15. Display the local dashboard address and recovery commands.

The installer must be idempotent.

Never format, repartition, erase or overwrite an existing disk automatically.

Never run an unverified remote script as administrator or root.

Provide an uninstall option that preserves user data by default.

---

# 5. ADMIN ACCOUNT AND AUTHENTICATION

During initial setup, ask for:

* Administrator username.
* Administrator password.
* Password confirmation.
* Initial storage directory.
* Access mode.
* Update policy.

Password requirements:

* Never echo passwords in terminal output.
* Never store plaintext passwords.
* Use a modern password hashing algorithm such as Argon2id.
* Never include credentials in command-line arguments or logs.
* Protect reset and recovery workflows.
* Do not ship default credentials.

Implement:

* Secure sessions.
* MFA/TOTP.
* Login rate limiting.
* Brute-force protection.
* Password reset.
* Session revocation.
* Optional account lockout with safe recovery.
* Public registration disabled by default.

Roles:

* Super Admin.
* Admin.
* Group Manager.
* User.

Enforce authorization in the backend for every protected endpoint.

---

# 6. FILE MANAGER

Build a full-featured file manager with:

* Upload and download.
* Drag-and-drop upload.
* Resumable large-file uploads.
* Upload progress and cancellation.
* Download progress where supported.
* Create folders.
* Rename files and folders.
* Copy and move.
* Bulk selection and actions.
* Delete and restore.
* Trash and permanent deletion.
* File version history.
* Favorites.
* Recent files.
* Breadcrumb navigation.
* Search by name.
* Search and filtering by metadata.
* Storage usage indicators.
* File and folder details.
* Keyboard shortcuts.
* Context menus.
* Multi-select.
* Grid and list views.
* Responsive mobile file manager.

Use safe filesystem APIs. Prevent path traversal, unsafe symbolic-link handling, unauthorized file access and directory escape.

Enforce quotas and permissions on the server, not merely in the UI.

Do not load complete file contents into memory for normal file operations.

---

# 7. SECURE FILE SHARING

Build a complete secure file-sharing system.

Users must be able to share a single file with:

* Specific registered users.
* Groups.
* An authenticated external recipient through a protected link.
* An optional public link, only when explicitly enabled by policy.

Share-link settings:

* Read-only access by default.
* Optional download permission.
* Optional preview-only permission.
* Optional password protection.
* Expiration date and time.
* Optional maximum download count.
* Optional maximum access count.
* Immediate revocation.
* Optional recipient authentication.
* Optional email verification.
* Optional download watermarking for supported formats.
* Optional restrictions on displaying sensitive metadata.

Security implementation:

* Generate cryptographically secure random tokens.
* Store only a secure hash of share tokens where practical.
* Make tokens unguessable.
* Never expose filesystem paths in URLs.
* Apply rate limiting to share endpoints.
* Do not leak filenames or metadata before authorization.
* Validate expiration and access limits server-side.
* Enforce permissions on every request.
* Use short-lived authorized download sessions.
* Support audit logs for successful and failed access attempts.
* Avoid logging full share tokens or passwords.
* Allow the owner or authorized admin to revoke access instantly.

The sharing interface must show:

* Share URL.
* Access type.
* Expiration.
* Password-protection status.
* Download limits.
* Created date.
* Last-access date, if available.
* Current status.
* Revoke and edit actions.

Copy links through an explicit user action.

Do not describe a link as secure merely because it is difficult to guess. Enforce actual authorization and revocation.

---

# 8. SECURE FOLDER SHARING

Implement folder sharing as a first-class feature.

Allow users to share a folder with:

* Specific users.
* Groups.
* Authenticated external recipients.
* An optional public link when allowed.

Folder permission options:

* Viewer: list and preview authorized contents.
* Downloader: view and download.
* Editor: upload, rename and modify permitted items.
* Uploader: upload without browsing existing contents, where technically supported.
* Custom permissions: configurable read, download, upload, edit and delete capabilities.

Requirements:

* Folder permissions must propagate to descendants according to explicit inheritance rules.
* Child-level restrictions must be supported without accidentally granting more access than the parent permits.
* Prevent unauthorized access through direct file URLs.
* Do not allow users to escape a shared folder using path manipulation.
* Support revoking access to the entire folder.
* Show the owner, recipients, permissions and inherited access.
* Provide an option to exclude selected subfolders.
* Provide optional read-only shares.
* Apply password, expiration and download policies to public folder shares.
* Ensure moving a file into or out of a shared folder cannot accidentally leak it.
* Audit permission changes and sensitive downloads.

Create tests for nested folders, permission inheritance, revoked access, direct URL access and conflicting permissions.

---

# 9. FILE PREVIEW ENGINE

Create a safe preview system.

Support previews where appropriate:

* Images: JPEG, PNG, WebP, GIF and other explicitly supported formats.
* PDF documents.
* Plain-text files.
* Markdown.
* JSON, CSV and log files.
* Audio playback.
* Video playback when the browser supports the codec.
* Code files with syntax highlighting.
* Office documents through an optional, isolated preview integration.

Preview security:

* Require authorization before reading or generating previews.
* Treat filenames, metadata and file contents as untrusted.
* Prevent stored XSS in HTML, Markdown, SVG and code previews.
* Do not execute uploaded HTML, JavaScript, macros or other active content.
* Render untrusted content safely, using sanitization or sandboxing as appropriate.
* Use content-security policies and safe response headers.
* Validate file signatures where useful.
* Enforce size, memory, CPU and processing-time limits.
* Generate thumbnails asynchronously.
* Cache generated previews with access-control checks.
* Provide a download option when previewing is unsupported.

Do not install a heavy document-conversion suite by default on 4 GB systems. Make resource-intensive document preview an optional component.

---

# 10. SORTING, FILTERING AND SEARCH

Implement professional file sorting and organization.

Sorting fields:

* Name ascending and descending.
* File size ascending and descending.
* Created date.
* Modified date.
* File type.
* Extension.
* Owner.
* Last accessed date, when available.
* Folder first.
* Custom manual order, if supported.

Filtering:

* File type.
* Extension.
* Size range.
* Date range.
* Owner.
* Shared status.
* Favorites.
* Recent files.

Search:

* Filename search.
* Extension and metadata filters.
* Search within authorized directories only.
* Optional content search as an opt-in feature.
* Keyboard shortcut to open search.
* Debounced search.
* Paginated results.

Requirements:

* Sorting and searching must respect permissions.
* Use server-side sorting and filtering for large collections.
* Index only necessary metadata.
* Avoid scanning all file contents on every search.
* Keep sorting preferences per user.
* Provide grid/list view and adjustable column visibility.
* Include clear filters, empty states and useful error messages.

---

# 11. STORAGE DISCOVERY AND QUOTAS

Discover existing storage safely.

Display:

* Device or volume.
* Total capacity.
* Used capacity.
* Available capacity.
* Filesystem.
* Mount point.
* Health information where supported.

Allow administrators to select existing storage roots.

Never automatically format a disk.

Implement:

* Global capacity monitoring.
* Per-user quotas.
* Group quotas if supported.
* Allocated versus actual used space.
* Warnings at configurable usage thresholds.
* Upload rejection when quotas are exceeded.
* Handling for full disks and unavailable storage roots.
* Multiple storage roots.
* Safe handling of removable drives.
* Storage reports.

Ensure quota checks are concurrency-safe. Two simultaneous uploads must not bypass the user's limit.

---

# 12. CLOUDFLARE TUNNEL

Integrate the official cloudflared client for worldwide access.

Features:

* Guided authentication.
* Domain and hostname configuration.
* Tunnel creation or selection.
* Tunnel service installation.
* Connection health.
* Diagnostics.
* Restart and reconnect controls.
* Clear error reporting.
* Optional Cloudflare Access integration.
* Local-only operation when no tunnel is configured.

Requirements:

* Do not require public IP addresses or router port forwarding.
* Keep local service ports private.
* Never expose SSH or database ports through the public file-sharing hostname.
* Protect credentials with OS-appropriate permissions.
* Never log tokens.
* Keep application-level authentication enabled.
* Apply authorization to every API and file endpoint.
* Document DNS, domain, account and network prerequisites.

---

# 13. macOS-INSPIRED UI AND MOTION

Create an original premium interface inspired by the polish and motion of macOS.

Design:

* Elegant sidebar.
* Compact top navigation.
* Refined typography.
* Light and dark themes.
* Selective translucent panels.
* Smooth transitions.
* Rounded cards.
* Animated modals and menus.
* Toast notifications.
* Storage usage visualizations.
* File grid and list views.
* Responsive mobile layouts.
* Accessible contrast and keyboard navigation.

Animation behavior:

* Animate file uploads with genuine progress.
* Animate sidebar expansion and page transitions.
* Animate dialog opening and closing.
* Use subtle loading indicators.
* Show real status changes for tunnel connectivity and updates.
* Avoid continuous decorative animations that waste CPU.
* Respect reduced-motion preferences.
* Disable heavy blur and complex transitions in Low Resource Mode.

Do not reproduce Apple logos, proprietary icons or copyrighted UI assets.

---

# 14. CLI AND ANIMATED INSTALLER

Create a CLI named `sego`.

Required commands:

* `sego install`
* `sego status`
* `sego doctor`
* `sego storage list`
* `sego tunnel setup`
* `sego tunnel status`
* `sego update check`
* `sego update`
* `sego backup`
* `sego restore`
* `sego logs`
* `sego uninstall`

Add a professional terminal experience:

* SEGO Server banner.
* Animated spinner.
* Dependency detection.
* Hardware discovery.
* Storage selection.
* Package installation.
* Release verification.
* Service configuration.
* Health checks.
* Final summary.

Animations must reflect real operations, never fake success.

Support quiet, verbose, non-interactive and no-animation modes.

---

# 15. AUTOMATIC UPDATES

Implement an update manager that:

* Checks the official release feed.
* Identifies newer versions.
* Shows release notes and security fixes.
* Verifies release signatures.
* Checks compatibility and disk space.
* Backs up the database and configuration.
* Applies migrations safely.
* Restarts the service.
* Performs health checks.
* Supports rollback.

Keep application updates separate from operating-system updates.

Offer manual updates, scheduled checks and configurable automatic application updates.

Never silently apply major OS upgrades or destructive package changes.

---

# 16. BACKUP AND RECOVERY

Implement:

* Scheduled backups.
* Database-consistent backup.
* Configuration backup.
* Optional encrypted backup archives.
* File versioning.
* Trash retention.
* Restore by file, folder or supported backup snapshot.
* Restore validation.
* Backup status and notifications.

Support external disks and optional off-device backup destinations.

Warn when backup storage is on the same physical disk as the primary data.

Test restoration rather than merely testing backup creation.

---

# 17. SECURITY TESTING

Test:

* Unauthorized file access.
* Cross-user IDOR vulnerabilities.
* Path traversal.
* Malicious filenames.
* Expired share links.
* Revoked share links.
* Share password brute-force protection.
* Folder permission inheritance.
* Direct URL bypass.
* Upload abuse.
* Stored XSS.
* CSRF.
* Session fixation and invalidation.
* Quota bypass through concurrent uploads.
* Unauthorized preview access.
* Secret leakage in logs.
* Unsafe update packages.
* Cloudflare Tunnel exposure and authentication configuration.

Run dependency security checks and document the threat model.

---

# 18. TESTING AND RELEASE QUALITY

Implement automated tests for:

* Fresh installation.
* Reinstallation.
* Windows and Linux service startup.
* Admin setup.
* Authentication.
* File CRUD operations.
* Secure file sharing.
* Folder sharing.
* Preview authorization.
* Sorting and search.
* Quota enforcement.
* Cloudflare Tunnel configuration.
* Updates and rollback.
* Backup and restoration.
* Interrupted transfers.
* Disk-full scenarios.

Benchmark:

* Idle RAM.
* CPU use during normal file operations.
* Concurrent uploads and downloads.
* Large-directory performance.
* Thumbnail generation.
* Startup time.
* Browser bundle size.

Publish measured results and hardware specifications. Do not claim guaranteed 4 GB RAM performance until it is validated.

---

# 19. GITHUB REPOSITORY

Create a complete repository with:

* README.md
* LICENSE
* SECURITY.md
* CONTRIBUTING.md
* INSTALL.md
* UPGRADE.md
* BACKUP-RESTORE.md
* ARCHITECTURE.md
* THREAT-MODEL.md
* CLI reference.
* API documentation.
* Installer scripts.
* Application source.
* Automated tests.
* CI workflows.
* Dependency update automation.
* Release signing and verification.

Use an OSI-approved open-source license and document third-party dependencies.

Do not publish fake release URLs or claim tests passed without executing them.

---

# 20. IMPLEMENTATION ORDER

Phase 1: Architecture, threat model, repository structure and supported OS matrix.

Phase 2: Cross-platform CLI, installation, admin bootstrap and secure service setup.

Phase 3: Storage discovery, file manager, database and quotas.

Phase 4: Secure file links, folder sharing and permissions.

Phase 5: Preview engine, sorting, filtering and search.

Phase 6: Cloudflare Tunnel integration.

Phase 7: macOS-inspired responsive UI and resource-aware animations.

Phase 8: Update manager, backup and recovery.

Phase 9: Security tests, resource benchmarks, documentation and signed releases.

First build a fully functional end-to-end MVP. Then expand the feature set incrementally.

Prioritize data safety, privacy, low resource consumption, verified updates and reliable recovery over decorative features.
