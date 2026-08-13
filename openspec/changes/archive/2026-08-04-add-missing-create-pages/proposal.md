## Why

23 modules have "Crear" buttons that navigate to `/new` routes, but only 4 have `new/page.tsx` files. The other 19 show "not found" errors when users click create.

## What Changes

Add `new/page.tsx` to 19 modules. Each page renders the corresponding form component in create mode (no props, defaults to `isEdit=false`).

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none — specs skipped)

## Impact

- 19 new files created under `frontend/src/app/(dashboard)/`
- No API changes, no backend changes
- Each file is ~15 lines, following the established pattern from `core/users/new/page.tsx`
