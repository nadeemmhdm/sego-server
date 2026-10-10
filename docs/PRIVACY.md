# Privacy Model

SEGO is designed as a self-hosted system. File contents and application metadata remain on infrastructure selected by the operator unless the operator explicitly configures external services such as remote access, email, off-device backup, or future integrations.

## Data categories

SEGO may process:

- account identifiers and password hashes
- session metadata
- file and folder names
- file metadata and contents
- sharing policy and access state
- diagnostic and audit information

## External services

The MVP does not require telemetry for normal local operation. Future integrations must be opt-in where they send user data or metadata to third parties, and documentation should identify what leaves the server.

## Logs

Logs should avoid passwords, full session tokens, full share tokens, and unnecessary file content. Operators should protect logs according to the sensitivity of their deployment.
