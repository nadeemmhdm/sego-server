# Testing Strategy

## Unit tests

Priority areas include:

- password hashing/verification
- token generation/hashing
- storage confinement
- namespace translation
- runtime-mode configuration
- quota reservation logic

## Integration tests to expand

- admin bootstrap
- login/session expiration/revocation
- file CRUD
- cross-user IDOR attempts
- direct URL access
- share expiry/revocation
- quota races
- disk-full handling
- interrupted transfer cleanup

## Platform validation

CI should build/test on Windows and Linux. Claims of platform support should be tied to actual tested versions and architectures.

## Performance validation

Measure idle RAM, CPU under transfers, concurrent transfer behavior, large directory listing, startup time, and frontend bundle size. Publish hardware specifications with benchmark results.
