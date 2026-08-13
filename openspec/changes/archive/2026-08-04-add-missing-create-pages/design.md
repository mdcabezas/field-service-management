## Context

Each module follows a list → create → detail pattern. The form components already support create mode (`isEdit=false`). Only the page wrappers are missing.

## Goals / Non-Goals

**Goals:**
- Add `new/page.tsx` to all 19 modules
- Follow the exact pattern from existing create pages

**Non-Goals:**
- No form changes, no API changes

## Decisions

Each page is a minimal wrapper: imports the form component, renders it with a title. No logic, no state.

## Risks / Trade-offs

None — trivial file additions completing existing UI.
