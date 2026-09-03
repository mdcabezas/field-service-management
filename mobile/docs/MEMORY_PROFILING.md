# Memory Profiling Guide

## Overview

This document provides guidance for profiling memory usage in the Localis Mobile app to ensure it stays under the 150MB target on low-end devices.

## Profiling Tools

### 1. React Native Performance Monitor

Enable in development builds:
```typescript
import { PerformanceMonitor } from 'react-native-performance';

// In your root component
if (__DEV__) {
  PerformanceMonitor.start();
}
```

### 2. Hermes Memory Profiling

Hermes provides built-in memory profiling:

```bash
# Enable heap profiling
adb shell am broadcast -a com.facebook.react.HMESampleHeapSampling -n com.localis.mobile/.MainActivity

# Dump heap
adb shell am dumpheap com.localis.mobile /data/local/tmp/heap.hprof

# Pull heap dump
adb pull /data/local/tmp/heap.hprof .
```

### 3. Flipper Integration

Add to `app.json`:
```json
{
  "expo": {
    "plugins": [
      [
        "expo-dev-client",
        {
          "flipper": true
        }
      ]
    ]
  }
}
```

## Memory Monitoring Points

### 1. Database Operations

Monitor memory during SQLite operations:
```typescript
// In lib/database.ts
export function getMemoryUsage(): number {
  // Estimate based on cached queries
  return process.memoryUsage().heapUsed / 1024 / 1024;
}
```

### 2. Photo Operations

Monitor during photo capture/compression:
```typescript
// In lib/photo-upload.ts
export async function monitorPhotoMemory(): Promise<void> {
  const before = process.memoryUsage().heapUsed;
  
  // Photo operation here
  
  const after = process.memoryUsage().heapUsed;
  const delta = (after - before) / 1024 / 1024;
  
  if (delta > 10) {
    console.warn(`Photo operation used ${delta.toFixed(2)}MB`);
  }
}
```

### 3. Sync Operations

Monitor during sync:
```typescript
// In lib/sync-engine.ts
export async function monitorSyncMemory(): Promise<void> {
  const startMem = process.memoryUsage().heapUsed;
  
  // Sync operation
  
  const endMem = process.memoryUsage().heapUsed;
  console.log(`Sync memory delta: ${(endMem - startMem) / 1024 / 1024}MB`);
}
```

## Memory Budget

| Component | Target | Notes |
|-----------|--------|-------|
| JavaScript Bundle | < 10MB | Hermes compiled |
| SQLite Cache | < 20MB | 2000 pages * 4KB |
| Photo Buffers | < 30MB | Max 3 photos in memory |
| Image Cache | < 20MB | LRU with 10MB limit |
| Sync Queue | < 10MB | Max 1000 items |
| UI Components | < 10MB | Virtualized lists |
| **Total** | **< 100MB** | Target: 150MB max |

## Memory Leak Detection

### 1. Interval/Timeout Cleanup

```typescript
// Always cleanup intervals
useEffect(() => {
  const interval = setInterval(() => {}, 5000);
  return () => clearInterval(interval);
}, []);
```

### 2. Event Listener Cleanup

```typescript
useEffect(() => {
  const subscription = EventEmitter.addListener('event', handler);
  return () => subscription.remove();
}, []);
```

### 3. SQLite Connection Cleanup

```typescript
// In App.tsx
useEffect(() => {
  return () => {
    closeDatabase();
  };
}, []);
```

## Performance Targets

| Metric | Target | Measurement |
|--------|--------|-------------|
| Startup Time | < 3s | Time to interactive |
| Memory Usage | < 150MB | Peak usage |
| JS Bundle Size | < 15MB | Compiled size |
| Photo Memory | < 30MB | Per operation |
| Sync Memory | < 10MB | Per sync cycle |

## Profiling Commands

### Android
```bash
# Get memory info
adb shell dumpsys meminfo com.localis.mobile

# Get heap info
adb shell dumpsys meminfo com.localis.mobile | grep "Heap"

# Monitor in real-time
adb shell top -d 1 | grep localis
```

### iOS
```bash
# Use Instruments
xcodebuild -workspace ios/LocalisMobile.xcworkspace -scheme LocalisMobile -sdk iphonesimulator

# Or use Memory Graph Debugger in Xcode
```

## Common Memory Issues

### 1. Large Photo Buffers
**Problem:** Loading full-resolution photos into memory.
**Solution:** Use expo-image-manipulator to resize before display.

### 2. Unbounded Cache
**Problem:** SQLite cache growing without limit.
**Solution:** Use LRU cache with max size.

### 3. Event Listener Leaks
**Problem:** Not removing listeners on unmount.
**Solution:** Always cleanup in useEffect return.

### 4. Large JSON Parsing
**Problem:** Parsing large sync responses.
**Solution:** Stream and process incrementally.

## Monitoring in Production

Add to `lib/memory-monitor.ts`:
```typescript
export class MemoryMonitor {
  private static instance: MemoryMonitor;
  private readings: number[] = [];
  
  static getInstance(): MemoryMonitor {
    if (!MemoryMonitor.instance) {
      MemoryMonitor.instance = new MemoryMonitor();
    }
    return MemoryMonitor.instance;
  }
  
  recordReading(): void {
    const mem = process.memoryUsage().heapUsed / 1024 / 1024;
    this.readings.push(mem);
    
    if (this.readings.length > 100) {
      this.readings.shift();
    }
    
    if (mem > 150) {
      console.error(`Memory warning: ${mem.toFixed(2)}MB`);
    }
  }
  
  getAverage(): number {
    return this.readings.reduce((a, b) => a + b, 0) / this.readings.length;
  }
  
  getMax(): number {
    return Math.max(...this.readings);
  }
}
```