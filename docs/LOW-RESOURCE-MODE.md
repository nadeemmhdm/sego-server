# Low Resource Mode

SEGO is designed to remain practical on older computers, including hosts with around 4 GB of total RAM. That is a design target, not a certified benchmark result.

## Current MVP behavior

Low-resource mode currently reduces upload concurrency to one and exposes the mode in status/diagnostic output.

Future low-resource controls should also cover:

- thumbnail/preview worker limits
- dashboard refresh frequency
- reduced visual effects
- download concurrency
- background indexing
- cache ceilings

## Enable

```bash
sego --config /path/to/sego.json mode set low-resource
```

## Disable / standard mode

```bash
sego --config /path/to/sego.json mode set standard
```

## Benchmark policy

Do not claim the 4 GB target is met until measurements have been captured on representative older hardware under idle, transfer, large-directory, preview, and update workloads.
