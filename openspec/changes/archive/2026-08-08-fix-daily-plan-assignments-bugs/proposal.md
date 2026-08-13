## Why

The Planning → Asignaciones module is broken in two blocking ways: (1) the list page at `/planning/assignments` passes an empty `dailyPlanId` to `useDailyPlanAssignments("")`, which sets `enabled: false` and silently never fetches, so users see an always-empty table; (2) the create form at `/planning/assignments/new` invokes `useCreateDailyPlanAssignment("")` whose mutation spreads `{ ...data, daily_plan_id: "" }` *after* the user payload, overwriting the typed `daily_plan_id` with `""` and producing orphaned assignments the backend rejects (FK violation).

Beyond the bugs, the assignment UX is unusable: `tech_id` and `role_id` are raw UUID text inputs with no dropdown lookup, even though `/api/users` (technicians) and `/api/tech-roles` are available. There is also no edit flow — `useUpdateDailyPlanAssignment` doesn't exist, and `/planning/assignments/[id]/page.tsx` is read-only with a delete button — forcing users to delete + recreate to fix a typo.

## What Changes

- Fix `/planning/assignments/page.tsx` to load all assignments (drop the broken `useDailyPlanAssignments("")` empty-id pattern; use a list-all query or feed a valid `dailyPlanId` from context).
- Fix `/planning/assignments/new/page.tsx` so the user-entered `daily_plan_id` is preserved. Either drop the `dailyPlanId` param from `useCreateDailyPlanAssignment("")` (call with no override) or refactor the hook so the call site passes `daily_plan_id` inside the payload without being overwritten.
- Replace the raw UUID `<Input>` for `tech_id` with a `<select>` populated from `useUsers` filtered to `role === "technician"` (label: name, value: id).
- Replace the raw UUID `<Input>` for `role_id` with a `<select>` populated from `useTechRoles` filtered to `active === true` (label: name, value: id).
- Add `useUpdateDailyPlanAssignment(dailyPlanId)` hook in `frontend/src/hooks/use-daily-plan-assignments.ts` issuing `PUT /api/daily-plan-assignments/{id}`.
- Wire `/planning/assignments/[id]/page.tsx` to an editable form (the existing "Guardar" submit branch only fires for `isNew`; extend it to call `updateAssignment.mutateAsync` when editing an existing record).
- **BREAKING**: none. All changes are additive or bug fixes; no API contract changes, no removed routes.

## Capabilities

### New Capabilities
- `planning-assignments`: Editing flow for daily-plan technician assignments — list all assignments, create with valid `daily_plan_id`, edit existing assignment, delete; dropdowns for technician and TechRole selection.

### Modified Capabilities
<!-- None. No existing spec in openspec/specs/. -->

## Impact

- **Frontend pages:**
  - `frontend/src/app/(dashboard)/planning/assignments/page.tsx` (list fix)
  - `frontend/src/app/(dashboard)/planning/assignments/new/page.tsx` (create fix + dropdowns)
  - `frontend/src/app/(dashboard)/planning/assignments/[id]/page.tsx` (edit form wiring)
- **Frontend hooks:**
  - `frontend/src/hooks/use-daily-plan-assignments.ts` (add `useUpdateDailyPlanAssignment`; possibly refactor `useCreateDailyPlanAssignment` so payload `daily_plan_id` is not overwritten)
- **Frontend forms:** may extract `AssignmentForm` into `frontend/src/components/forms/assignment-form.tsx` to share between `/new` and `/[id]` (optional, follows the pattern of `user-form.tsx` and `tech-role-form.tsx`).
- **Existing hooks reused (no changes):** `useUsers` (`/api/users`), `useTechRoles` (`/api/tech-roles`).
- **Backend:** no changes required. `PUT /api/daily-plan-assignments/{id}` must already exist on the backend; if it doesn't, add it (verify during implementation).
- **RBAC:** unchanged — `/planning` access stays `["admin","manager"]`; actions `create`/`edit` stay `["admin","manager"]`; `delete` stays `["admin"]`.
- **Out of scope:** Technician stub on customer detail, Visits/Assignments UI, profile/self-edit page, password reset flow.
