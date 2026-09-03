# Mobile Visits & Reports App

## Problem
Field technicians currently have no mobile tool to manage their daily work. All visit tracking, checklist management, photo documentation, and reporting is done on paper or via the web dashboard (desktop-only). This leads to:
- Lost or incomplete paperwork
- Delayed data entry (end of day)
- No real-time visibility for dispatch
- No GPS verification of site attendance
- No standardized photo documentation
- Reports must be manually transcribed from paper

## Proposed Solution
Build a React Native + Expo mobile app (offline-first) for field technicians to manage visits, checklists, photos, GPS checkpoints, measurements, and reports.

## Key Features
1. **Offline-first visit management** — full visit lifecycle with optimistic sync
2. **Configurable checklists by stages** — templates with required items that block status transitions
3. **Photo capture pipeline** — compress → store → upload → thumbnail generation
4. **Manual GPS capture** — arrival/departure with distance validation and waiver system
5. **Partner-specific customization** — templates, SLAs, and data fields scoped to partner
6. **Property history** — previous visits to same property visible in visit detail
7. **Audit trail** — all changes visible in visit detail
8. **Reopen visits** — within configurable window with reason tracking
9. **Material usage tracking** — planned vs actually used
10. **Configurable photo retention** — days and storage cap

## Target Users
- Field technicians (iOS + Android, low-end devices 2GB RAM)
- Dispatchers (real-time visibility via web dashboard)
- Operations managers (reports and audit via web dashboard)

## Constraints
- Low-end Android/iOS devices (2GB RAM minimum)
- Low-bandwidth network conditions (2G/3G)
- Must work fully offline (no connection required)
- Battery optimization critical (8+ hour workday)
- Data cost optimization in developing countries
- Monospace brutalist design system matching web frontend

## Stack
- Expo SDK 57, React Native 0.79
- `@op-engineering/op-sqlite` (~500KB, raw SQLite)
- Zustand 5.x (state management)
- Expo Router (navigation)
- expo-camera, expo-location, expo-image-picker
- expo-image-manipulator (compression)
- expo-sharing (photo sharing)
- No WatermelonDB (too heavy for target devices)

## Success Criteria
- [ ] Technician can manage full visit lifecycle offline
- [ ] Checklists enforce required items before status transitions
- [ ] Photos compressed to <200KB, thumbnails to <20KB
- [ ] GPS distance validation against property coordinates
- [ ] Partner-specific templates and data auto-resolved
- [ ] Full audit trail visible in visit detail
- [ ] Sync completes in <5 seconds on 3G
- [ ] App stays under 50MB RAM usage
- [ ] 80% test coverage (backend, mobile, frontend E2E)

## Out of Scope
- Video capture
- Voice notes
- Offline maps
- Push notifications (future phase)
- Bluetooth scanner integration (future phase)
- Multi-language (future phase)
