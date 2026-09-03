# Performance Optimization

Target: Low-end Android/iOS (2GB RAM), 2G/3G networks

---

## Current Optimizations

### Bundle
- Hermes JS engine enabled
- Expo Router with typed routes (tree-shaking)
- Minimal dependencies (~500KB SQLite)

### Database
- PRAGMAs optimized for mobile:
  - `journal_mode = MEMORY` (no fsync on commit)
  - `synchronous = NORMAL` (balance safety/speed)
  - `cache_size = 2000` (8MB page cache)
  - `temp_store = MEMORY` (temp tables in RAM)
  - `mmap_size = 64MB` (memory-mapped I/O)

### Images
- Compression: 800px, 70% JPEG (~200KB)
- Thumbnail: 200px, 60% JPEG (~20KB)
- Max original: 1280x720

### Sync
- Full sync: ~50KB reference data + paginated visits
- Delta sync: <10KB typical
- Outbox ordering reduces round-trips
- Photo upload separate from JSON sync

### UI
- Brutalist: no shadows, no gradients, solid colors
- FlatList for lists (windowing)
- Lazy-loaded screens (Expo Router)
- Zustand for minimal re-renders

---

## Areas for Improvement

### 1. Bundle Size
```bash
# Analyze bundle
npx expo-export --dump-asset-map
# or
npx react-native-bundle-visualizer
```

Target: <15MB JS bundle

### 2. Startup Time
- Lazy-load heavy modules (camera, image-manipulator)
- Defer non-critical initialization
- Use `expo-splash-screen` to hide until ready

### 3. Memory
- Photo cleanup on app start
- Limit FlatList window size
- Release image references after upload

### 4. Network
- Request deduplication
- Response caching (React Query / SWR)
- Offline queue with persistence

### 5. Database
- Add indexes for common queries
- Consider WAL mode for concurrent reads
- Vacuum periodically

---

## Recommended Changes

### Add Indexes
```sql
-- In initializeDatabase()
CREATE INDEX idx_visits_status ON visits(status);
CREATE INDEX idx_visits_scheduled ON visits(scheduled_at);
CREATE INDEX idx_checklist_items_visit ON checklist_items(visit_id);
CREATE INDEX idx_checkpoints_visit ON checkpoints(visit_id);
CREATE INDEX idx_photos_visit ON photos(visit_id);
CREATE INDEX idx_outbox_status ON outbox(status);
```

### Lazy Load Camera
```tsx
// In photo-capture.tsx
const [ImagePicker, setImagePicker] = useState<any>(null);

useEffect(() => {
  import('expo-image-picker').then(setImagePicker);
}, []);
```

### Virtualized Lists
```tsx
// In visit list
<FlatList
  windowSize={5}
  maxToRenderPerBatch={10}
  updateCellsBatchingPeriod={50}
  initialNumToRender={10}
/>
```

### Memoize Components
```tsx
// In StatusTransitionComponent
const TransitionButton = React.memo(({...}) => (
  <Button {...} />
));

// In Checklist item
const ChecklistItem = React.memo(({item, onPress}) => (
  // ...
));
```

### Request Deduplication
```typescript
// In apiRequest
const pendingRequests = new Map<string, Promise<any>>();

export async function apiRequest(path, options) {
  const key = `${options.method || 'GET'}:${path}`;
  if (pendingRequests.has(key)) {
    return pendingRequests.get(key);
  }
  const promise = fetch(...);
  pendingRequests.set(key, promise);
  promise.finally(() => pendingRequests.delete(key));
  return promise;
}
```

---

## Profiling Commands

```bash
# JS profiling
npx react-native-performance-monitor

# Native profiling (Android)
adb shell am start -n com.localis.mobile/.MainActivity --es "reactNativeProfiling" "true"

# Bundle analyzer
npx expo-export --dump-asset-map | node scripts/analyze-bundle.js

# Memory
adb shell dumpsys meminfo com.localis.mobile
```

---

## CI Performance Gates

```yaml
# .github/workflows/performance.yml
- name: Bundle size
  run: |
    SIZE=$(npx expo-export --dump-asset-map | jq '.bundleSize')
    if [ $SIZE -gt 15000000 ]; then
      echo "Bundle too large: $SIZE bytes"
      exit 1
    fi
```