## Context

The backend Go models define the canonical API contract via JSON struct tags. The frontend TypeScript types were hand-written approximations that diverge in field names, missing fields, and phantom fields. Every module is affected.

## Goals / Non-Goals

**Goals:**
- Rewrite all 35 TypeScript interfaces to exactly match backend JSON tags
- Update all pages/forms/hooks to reference correct field names
- Fix "Invalid date" and all data display bugs

**Non-Goals:**
- No backend changes
- No new features or behavior changes
- Minimal form logic changes — field renames and adding backend-required fields that were missing

## Decisions

1. **Source of truth is the Go struct JSON tags** — every TS interface field must map 1:1 to a backend JSON key
2. **Optional fields use `?:` in TypeScript** — matching Go `omitempty` pointer fields
3. **Remove phantom `updated_at`** from interfaces where Go struct has no such field
4. **Rename fields to match Go**: `rut` → `tax_id`, `date` → `scheduled_at`, `start_date` → `departure_time`, etc.
5. **Add missing fields** that Go emits but TS doesn't declare
6. **Remove fields** that TS declares but Go doesn't emit (e.g., `description` on TechRole)

## Risks / Trade-offs

- Forms may need zod schema updates to match renamed fields
- Some pages reference phantom fields that will become `undefined` until column definitions are updated
- Build will break temporarily until all references are updated in sync
