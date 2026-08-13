## Why

The Partners and Customers modules have list pages with "+ Crear Socio" and "+ Crear Cliente" buttons that navigate to `/partners/new` and `/customers/new`, but those routes have no `page.tsx` files. Users see "Socio no encontrado" / "Cliente no encontrado" errors when clicking create.

## What Changes

- Add `partners/new/page.tsx` to render `PartnerForm` in create mode
- Add `customers/new/page.tsx` to render `CustomerForm` in create mode

No API changes, no behavior changes — completing the existing UI pattern already established in `core/users/new` and `core/tech-roles/new`.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none)

## Impact

- Files created: `frontend/src/app/(dashboard)/partners/new/page.tsx`, `frontend/src/app/(dashboard)/customers/new/page.tsx`
- No dependencies added, no breaking changes
