## Why

The frontend TypeScript types in `types/api.ts` were written as approximations of the backend Go models, not derived from them. Out of 35 interfaces compared, 9 are structurally divergent (describe different schemas), 11 have significant mismatches, 10 have moderate mismatches, and 5 have minor mismatches. This causes:

- "Invalid date" on visits list (frontend expects `date`, backend returns `scheduled_at`)
- "Invalid date" on vehicle-assignments (frontend expects `start_date`/`end_date`, backend returns `departure_time`/`return_time`)
- Empty/wrong data display across all 9 modules
- Forms that submit fields the backend doesn't expect
- Missing fields that the backend actually returns

## What Changes

Rewrite all 35 TypeScript interfaces in `types/api.ts` to exactly match the backend Go struct JSON tags. Then update all list pages, detail pages, forms, and hooks to reference the correct field names.

## Capabilities

### New Capabilities

(none — this is a type correction, not new behavior)

### Modified Capabilities

(none — specs are skipped)

## Impact

- Files modified: `frontend/src/types/api.ts` (all 35 interfaces)
- Files modified: ~20 list/detail/form pages across 9 modules
- No API changes, no backend changes, no new dependencies
- Breaking change: all frontend data references shift to backend-correct field names
