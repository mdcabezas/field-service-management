## Context

The frontend has 9 CRUD modules. Each module follows the pattern: list page → create page → detail/edit page. The Partners and Customers modules are missing their create pages. The form components (`PartnerForm`, `CustomerForm`) already support create mode (`isEdit=false`).

## Goals / Non-Goals

**Goals:**
- Add missing `new/page.tsx` files for Partners and Customers
- Follow the exact pattern from `core/users/new/page.tsx`

**Non-Goals:**
- No form changes, no API changes, no new components

## Decisions

Follow the established pattern exactly: minimal page component that renders the form with no props (defaults to create mode).

## Risks / Trade-offs

None — this is a trivial fix completing existing UI.
