# Performance Notes

Bubblefetch is built for speed with parallel collection, minimal external commands,
and targeted caching. This document captures benchmark results and key optimizations.

## Benchmarks

Quick comparison:

```bash
time ./bubblefetch --export text > /dev/null   # ~0.003s
time neofetch > /dev/null                      # ~0.350s
time fastfetch > /dev/null                     # ~0.012s
```

Example benchmark output:

```
Running 10 iterations...
Average: 1.34ms
Total: 13.42ms
```

JSON benchmark output:

```bash
bf --benchmark --format json
```

```json
{
  "runs": 10,
  "average_ms": 1.34,
  "total_ms": 13.42,
  "durations_ms": [1.1, 1.2, 1.4],
  "timestamp": "2025-01-18T12:00:00Z",
  "version": "0.3.0",
  "remote": false
}
```

## Why it is fast

- Parallel goroutines for independent metrics.
- /proc and /sys reads instead of shelling out.
- Cached OS detection and theme lookups.
- GPU detection prefers `/sys/class/drm`, with a short fallback timeout.

## Optimization highlights

- GPU vendor mapping avoids unnecessary `lspci` calls.
- OS detection cached with `sync.Once`.
- Theme loading avoids repeated disk reads.

## How to test locally

```bash
go build -ldflags="-s -w" -o bubblefetch ./cmd/bubblefetch
./bubblefetch --benchmark
time ./bubblefetch --export text > /dev/null
```

## Sample test environment

Recorded on Arch Linux (i9-13900KF, RTX 3060 Ti), optimized build size ~6.2MB.

## Benchmark methodology

The "1.2ms" headline is the average of 10 back-to-back runs on a warmed cache
using the optimized binary (`-ldflags="-s -w"`). To reproduce:

1. Close heavy background tasks.
2. Run `bf --benchmark` multiple times.
3. Use the average reported per run (not the `time` shell wrapper).

For apples-to-apples comparisons, ensure the same module list and that remote
collection is disabled.
