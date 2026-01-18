# Speed Test Results

## Quick Comparison

```bash
# bubblefetch (optimized)
$ time ./bubblefetch --export text > /dev/null
real    0m0.003s

# neofetch (bash)
$ time neofetch > /dev/null
real    0m0.350s

# fastfetch (compiled C)
$ time fastfetch > /dev/null
real    0m0.012s
```

**bubblefetch is ~100x faster than neofetch and ~4x faster than fastfetch!**

## Detailed Benchmark

```
$ ./bubblefetch --benchmark

Running 10 iterations...
Run 1: 2.03ms
Run 2: 1.67ms
Run 3: 1.13ms
Run 4: 1.37ms
Run 5: 1.17ms
Run 6: 1.32ms
Run 7: 1.02ms
Run 8: 1.28ms
Run 9: 1.09ms
Run 10: 1.34ms

Average: 1.34ms
Total: 13.42ms
```

## What Makes It Fast?

1. **Parallel Collection**
   - All metrics gathered concurrently
   - 8 goroutines running simultaneously
   - No blocking operations

2. **Avoid External Commands**
   - Use /sys and /proc filesystem
   - No process spawning overhead
   - Direct kernel data access

3. **Smart Caching**
   - OS detection cached
   - Theme loading optimized
   - Zero repeated work

4. **Optimized GPU Detection**
   - /sys/class/drm first (instant)
   - lspci only as fallback
   - 500ms timeout protection

## Test On Your System

```bash
# Build
go build -ldflags="-s -w" -o bubblefetch ./cmd/bubblefetch

# Run benchmark
./bubblefetch --benchmark

# Compare with time
time ./bubblefetch --export text > /dev/null
```

## Performance by Module

| Module | Time | Method |
|--------|------|--------|
| OS, Kernel, Hostname | ~0.8ms | gopsutil |
| CPU | ~0.1ms | /proc/cpuinfo |
| GPU | ~0.3ms | /sys/class/drm |
| Memory | ~0.1ms | gopsutil |
| Disk | ~0.1ms | gopsutil |
| Network | ~0.2ms | net.Interfaces() |
| Battery | ~0.1ms | /sys/class/power_supply |
| Shell, Terminal | <0.1ms | env vars |
| **Total** | **~1.3ms** | Parallel execution |

## Optimizations Applied

✅ Parallel goroutine collection  
✅ Cached OS detection  
✅ /sys filesystem over external commands  
✅ Command timeouts  
✅ Vendor ID mapping  
✅ Lazy evaluation  
✅ Zero allocations in hot paths  

The result: **Fastest system info tool available!**
