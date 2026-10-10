# Security Hardening Guide

## Host

- keep Windows/Linux security updates current
- use a dedicated service identity
- limit interactive access to the host
- encrypt disks where appropriate
- maintain tested off-device backups

## Network

- keep SEGO bound to loopback by default
- use TLS for remote access
- do not expose SQLite, SSH, SMB, or administrative services through the public sharing hostname
- keep application authentication enabled behind tunnels/proxies
- restrict reverse-proxy headers to trusted proxy sources

## Accounts

- use unique strong administrator passwords
- create only the permissions users need
- revoke abandoned sessions/accounts
- enable MFA when the implementation becomes available and verified

## Sharing

- prefer authenticated/revocable links
- set expirations for external shares
- avoid public links for sensitive material
- review audit records when the full audit layer is available

## Files

Treat filenames, MIME types, metadata, Markdown, SVG, HTML, archives, and Office documents as untrusted. Active content must never be executed merely to generate a preview.

## Updates

Do not execute an update artifact unless its source and integrity can be verified. The signed update manager remains a roadmap item.
