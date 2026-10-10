# Release Process

SEGO should not publish a production release until the corresponding code, migrations, platform builds, and rollback path have been verified.

## Release checklist

1. update changelog
2. run backend tests
3. build the web application
4. run platform installation tests
5. run security checks
6. verify upgrade and rollback
7. generate platform artifacts
8. generate checksums
9. sign artifacts using the approved release-signing process
10. publish release notes with known limitations

Never publish fake download URLs, fabricated checksums, or unexecuted benchmark claims.
