# Localis Mobile

Offline-first React Native + Expo app for field technicians to manage visits, checklists, photos, GPS checkpoints, measurements, and reports.

## Features

- **Offline-first**: All data stored locally in SQLite, syncs when online
- **Visits Management**: List, filter, and view visit details
- **Checklists**: Multi-stage checklists with required items
- **Photo Capture**: Camera + gallery, auto-compression (800px, 70% JPEG)
- **GPS Checkpoints**: Manual capture with haversine distance validation
- **Measurements**: Key/value/unit with optional GPS
- **Material Usage**: Track materials with quantities and notes
- **Sync Engine**: Full sync on login, delta sync on reconnect, outbox pattern
- **Status Transitions**: Validated transitions with checklist blocking
- **Partner Context**: Auto-resolve partner from property, filter templates

## Tech Stack

- **Expo SDK 57**, React Native 0.86
- **@op-engineering/op-sqlite** - Raw SQLite (~500KB, no WatermelonDB)
- **Zustand 5.x** - State management
- **Expo Router** - File-based navigation
- **expo-camera**, **expo-image-picker** - Photo capture
- **expo-location** - GPS
- **expo-image-manipulator** - Compression
- **expo-sharing**, **expo-secure-store** - Sharing & tokens

## Design System

Brutalist style:
- Monospace system font
- `borderWidth: 2`, `borderRadius: 0`
- Solid status colors, no shadows
- Spacing: 4/8/16/24/32

## Project Structure

```
mobile/
├── app/                    # Expo Router screens
│   ├── (tabs)/            # Tab navigation
│   │   ├── index.tsx      # Visit list
│   │   ├── sync.tsx       # Sync status
│   │   └── settings.tsx   # Config
│   ├── login.tsx          # Login screen
│   └── visit/[id]/        # Visit detail stack
│       ├── reopen.tsx
│       ├── measurements.tsx
│       ├── usage.tsx
│       └── audit.tsx
├── components/ui/         # Brutalist components
├── lib/                   # Core libraries
│   ├── database.ts        # SQLite + PRAGMAs
│   ├── auth.ts            # Login, tokens, apiRequest
│   ├── sync-engine.ts     # Full/delta sync, outbox
│   ├── photo-upload.ts    # Compress, store, upload
│   ├── gps.ts             # Capture, distance, waivers
├── stores/                # Zustand stores
├── utils/                 # Helpers
│   ├── distance.ts        # Haversine
│   ├── ordered-sync.ts    # Outbox ordering
│   └── status-transition.ts # Transition rules
```

## Getting Started

```bash
cd mobile
npm install
npm run start
```

## Configuration

Environment variables:
- `EXPO_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8081)

App config (`app.json`):
- Android package: `com.localis.mobile`
- iOS bundle ID: `com.localis.mobile`
- Permissions: CAMERA, LOCATION, STORAGE

## Build

```bash
# Development build
eas build --profile development

# Preview (internal testing)
eas build --profile preview

# Production (app store)
eas build --profile production
```

## Database Schema

Key tables:
- `visits` - Visit assignments
- `checklist_items` - Checklist items with stage
- `checkpoints` - GPS checkpoints
- `gps_waivers` - GPS waivers
- `photos` - Photo metadata
- `measurements` - Key/value measurements
- `usages` - Material usage
- `outbox` - Pending sync operations
- `config` - App configuration

## Sync Protocol

**Full Sync** (on login):
1. GET `/api/mobile/sync/pull` → reference data + visits
2. Store in SQLite

**Delta Sync** (on reconnect):
1. GET `/api/mobile/sync/pull?since=timestamp` → changes
2. Apply changes

**Push** (on any change):
1. POST `/api/mobile/sync/push` with outbox items
2. Ordered: checkpoints → photos → usages → measurements → reports → status
3. Server returns results with server IDs
4. Update local IDs

## Status Transitions

```
assigned → en_route → in_progress → completed
    ↓           ↓              ↓
cancelled   cancelled      cancelled
    ↑                              ↑
    └────────────── reopen ────────┘ (within 24h)
```

Requires checklist completion for `completed`.

## Testing

```bash
# Unit tests
npm test

# Watch mode
npm test:watch

# Coverage
npm test:coverage
```

## License

Internal use only.