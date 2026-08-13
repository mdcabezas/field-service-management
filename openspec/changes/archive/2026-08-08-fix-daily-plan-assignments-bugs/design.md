## Context

The DailyPlanAssignment flow in the frontend is broken in two places (always-empty list page; `daily_plan_id` overwrite bug) and lacks an edit flow plus usable dropdowns. The backend (`/api/daily-plan-assignments`) currently exposes only `GET /:id`, `POST /`, and `DELETE /:id` — no `PUT /:id` and no list-all `GET /` route. The repo (`DailyPlanAssignmentRepo`) has `GetByID`, `ListByPlan`, `Create`, `Delete` — no `List` and no `Update`.

Three existing hooks already drive the dependent dropdowns: `useDailyPlans()`, `useUsers()`, `useTechRoles()`. Each returns `r.items` from `{ items, total }` envelopes. The existing `useCreateDailyPlanAssignment("")` pattern is the root of both frontend bugs: the hook takes a `dailyPlanId` parameter and silently embeds it into the body, which is wrong for the standalone create page where the user picks the plan themselves.

Two reference implementations to mirror: `use-vehicles.ts` and `use-tech-roles.ts` — both have a no-arg `useCreateX()`, a `useUpdateX()` taking `{id, data}`, and a `useDeleteX()` taking `id`. The same shape applies cleanly to assignments.

## Goals / Non-Goals

**Goals:**
- Restore functionality: list shows assignments, create persists the user-chosen `daily_plan_id`.
- Add a real edit flow (technician + TechRole editable; daily_plan immutable).
- Replace raw-UUID inputs with `<select>` dropdowns for `daily_plan_id`, `tech_id`, `role_id` on create; `tech_id` and `role_id` on edit.
- Add the two missing backend endpoints: `GET /api/daily-plan-assignments` (list all) and `PUT /api/daily-plan-assignments/:id`.
- Keep cache invalidation correct across both the `daily-plan-assignments` list and the parent daily-plan's assignment sub-list.

**Non-Goals:**
- No bulk assignment, no drag-and-drop scheduling UI, no calendar view.
- No move-assignment-to-another-plan flow (`daily_plan_id` immutable on edit).
- No Visits/Assignments UI (separate surface, separate change).
- No Technician stub on the customer detail page (separate surface, separate change).
- No audit log entries for assignment changes (can be added later as a follow-up; design choice below).
- No RBAC changes — `/planning` access stays `["admin","manager"]`, `delete` stays `["admin"]`.

## Decisions

### Decision 1: Refactor `useCreateDailyPlanAssignment` to no-arg (drop the `dailyPlanId` param)

**Choice:** Change `useCreateDailyPlanAssignment(dailyPlanId: string)` → `useCreateDailyPlanAssignment()` (no args). The hook's `mutationFn` takes `Omit<DailyPlanAssignment, "id" | "created_at">` (which includes `daily_plan_id` from the user payload) and POSTs the body as-is. Cache invalidation broadens to invalidate both `["daily-plan-assignments"]` (new top-level list key) and `["daily-plans", dailyPlanId, "assignments"]` per-plan sub-cache.

**Rationale:** The current signature forces the call site to pass a string that gets *spread after* the user's `data`, overwriting the user-supplied `daily_plan_id`. Making it no-arg is the minimal fix that preserves a clean call site for both the standalone `/new` page and future per-plan create buttons (which can call the same hook and supply `daily_plan_id` in the data).

**Alternative considered:** Keep the param and treat empty string as "use data.daily_plan_id" inside the mutation. Rejected — special-casing empty string is fragile and obscures intent.

**Alternative considered:** Split into `useCreateDailyPlanAssignment` (no arg) and `useCreateDailyPlanAssignmentForPlan(planId)`. Rejected — YAGNI; one no-arg hook fits all current call sites.

### Decision 2: Add `useAllDailyPlanAssignments()` hook + keep `useDailyPlanAssignments(planId)` for sub-lists

**Choice:** Add a new hook `useAllDailyPlanAssignments()` that GETs `/api/daily-plan-assignments` (the new list-all backend endpoint). Keep the existing `useDailyPlanAssignments(dailyPlanId)` for any sub-view nested under a daily plan (none exist today, but the daily-plan detail page may grow one). The top-level `/planning/assignments` list page switches to `useAllDailyPlanAssignments()`.

**Rationale:** Two distinct queries with distinct cache keys (`["daily-plan-assignments"]` vs `["daily-plans", planId, "assignments"]`) keeps invalidation precise. Mixing them would force over-invalidation or stale sub-lists.

**Alternative considered:** Rewrite `useDailyPlanAssignments(planId)` to fetch all and filter client-side. Rejected — wasteful and breaks the sub-list cache contract.

### Decision 3: Add `useUpdateDailyPlanAssignment()` hook with `{ id, data }` signature

**Choice:** `useUpdateDailyPlanAssignment()` returns a mutation whose `mutationFn` accepts `{ id: string; data: Partial<DailyPlanAssignment> }` and PUTs to `/api/daily-plan-assignments/${id}`. Cache invalidation on success: invalidate `["daily-plan-assignments"]` and `["daily-plans", "<affected plan id>", "assignments"]`. Since the hook doesn't know the affected plan id from the call site alone, the call site passes it in `meta` or invalidation invalidates the broad `["daily-plan-assignments"]` plus all `["daily-plans", "*", "assignments"]` keys (using `queryKey` prefix matching with `invalidateQueries({ queryKey: ["daily-plans"], refetchType: "active" })` is more aggressive but simpler).

**Rationale:** Mirrors the `useUpdateUser` / `useUpdateTechRole` shape for consistency. The unknown-affected-plan problem is solved by invalidating broadly and letting active sub-lists refetch.

**Alternative considered:** Pass `dailyPlanId` to the hook so it can target invalidation precisely. Rejected — couples the edit page to knowing the parent plan, which is fine but redundant since the broad invalidation is cheap for a small dataset.

### Decision 4: Single shared form component for new + edit

**Choice:** Extract an `AssignmentForm` component (`frontend/src/components/forms/assignment-form.tsx`) following the pattern of `user-form.tsx` and `tech-role-form.tsx`. The form takes a `mode: "new" | "edit"` prop and an optional `defaultValues`. Both pages render `<AssignmentForm mode="new" />` or `<AssignmentForm mode="edit" assignment={assignment} />`.

**Rationale:** Eliminates the existing duplication between `/new/page.tsx` and `/[id]/page.tsx` (which currently re-implements the same form). Keeps the `<select>` population logic in one place.

**Alternative considered:** Edit inline in the list. Rejected — inconsistent with every other entity in this codebase, which uses a dedicated `/[id]` edit page.

### Decision 5: Three dropdowns populated by their respective hooks, with empty-state options

**Choice:** `daily_plan_id` (create only), `tech_id`, `role_id` each render as `<select>` elements with a `"Seleccionar..."` placeholder option (empty value). When the underlying query returns an empty list, render a single disabled option ("No hay X disponibles") and disable the submit button. Use `useUsers()` filtered client-side to `role === "technician"`. Use `useTechRoles()` filtered client-side to `active === true`. Use `useDailyPlans()` as-is.

**Rationale:** Server-side filter-by-role endpoints don't exist for `/api/users`; client-side filter is acceptable since the technician population is small. Empty-state handling surfaces the configuration gap (e.g. "no technicians defined") instead of letting the form silently accept bad input.

**Alternative considered:** Add a backend `/api/users?role=technician` filter. Rejected for this change — scope creep. Mentioned here so a future change can pick it up.

### Decision 6: Backend adds `GET /api/daily-plan-assignments` and `PUT /api/daily-plan-assignments/:id`

**Choice:**

1. **Repo:** Add `List(ctx, limit, offset) ([]DailyPlanAssignment, error)` and `Update(ctx, id, techID, roleID) (*DailyPlanAssignment, error)` to `DailyPlanAssignmentRepo`. SQL for `List`: `SELECT id, daily_plan_id, tech_id, role_id, created_at FROM daily_plan_assignments ORDER BY created_at DESC LIMIT $1 OFFSET $2`. SQL for `Update`: `UPDATE ... SET tech_id = $2, role_id = $3 WHERE id = $1 RETURNING id, daily_plan_id, tech_id, role_id, created_at`.
2. **Handler:** Add `List` and `Update` methods to `DailyPlanAssignmentHandler` mirroring the existing `List`/`Update` shape on `VehicleHandler`. `List` reads `limit`/`offset` via `handler.ParsePagination(c)`. `Update` calls `repo.Update`.
3. **Router:** Add two lines: `api.GET("/daily-plan-assignments", dailyPlanAssignmentH.List)` and `api.PUT("/daily-plan-assignments/:id", adminManager, dailyPlanAssignmentH.Update)`. Both hook into the existing `adminManager` middleware for the write path; `List` inherits the `api` group's auth middleware.
4. **Handler always generates `ID`/`CreatedAt` on Create** — already done by a prior change (`fix-daily-plan-assignments-bugs` is a sibling to the just-fixed `inventory` handlers); if `Create` handler still doesn't generate `ID`/`CreatedAt`, the apply phase must include it as a sub-task (verify during implementation).

**Rationale:** Bash-the-bug-at-the-source — if the frontend calls these endpoints, the backend must exist. Putting both endpoints together keeps the round-trip cost low. `adminManager` mirrors the existing gates on POST/DELETE.

**Alternative considered:** Skip the backend `GET /` (list-all) and have the list page iterate daily plans, calling `useDailyPlanAssignments(planId)` per plan. Rejected — N+1 queries on every list view, fragile pagination, awkward UX (would need grouping by plan).

**Alternative considered:** PUT also accepts `daily_plan_id` to support move-between-plans. Rejected — out of scope (see Non-Goals); keeping it immutable reduces accidentally orphaning audit history.

### Decision 7: `daily_plan_id` shown as disabled-select on edit page

**Choice:** On the edit form, `daily_plan_id` renders as a disabled `<select>` populated with `useDailyPlans()` so the user can *see* what plan they're editing the assignment of, but cannot change it. Changing the parent plan requires deleting the assignment and creating a new one.

**Rationale:** Mirrors the `UserForm` `disabled` pattern on `email` during edit. Avoids the temptation to "just allow change" when the underlying semantics (which plan owns this assignment) are not trivially movable.

## Risks / Trade-offs

- **[Risk] Backend `Update` returns the full row but doesn't update `daily_plan_id`** → Mitigation: the repo's `Update` SQL only SETs `tech_id` and `role_id`; the RETURNING clause includes `daily_plan_id` so the response is consistent. Add a backend regression test in the apply phase if the backend has a test suite (verify during implementation).

- **[Risk] Cache invalidation after edit invalidates too much** → Mitigation: use targeted `invalidateQueries({ queryKey: ["daily-plan-assignments"] })` plus a broad `invalidateQueries({ queryKey: ["daily-plans"], refetchType: "active" })` only if a per-plan sub-list is currently mounted. Watch for infinite refetch loops in dev tools.

- **[Risk] The `useCreateDailyPlanAssignment` signature change is a breaking refactor** → Mitigation: the only two call sites are `/planning/assignments/new/page.tsx` and `/planning/assignments/[id]/page.tsx` (on `isNew` branch). Both will be updated in the same change. Lint + TypeScript catch any unfound callers.

- **[Risk] Client-side technician filtering misses users whose `role` field is null or unexpected** → Mitigation: the `User` type in frontend has `role: string` (not nullable) and the backend enforces the role enum. The `=== "technician"` predicate is sufficient.

- **[Risk] Empty-state dropdowns block the user from creating the first assignment** → Mitigation: the dropdown shows a hint linking to `/core/users/new` (for technicians) or `/core/tech-roles/new` (for roles) or `/planning/daily-plans/new` (for plans). This is a small affordance but important — without it the user is stuck.

- **[Trade-off] Two queries per `/planning/assignments/new` mount (daily plans + users + tech roles = 3 queries)** vs the single raw-uuid form. The cost is justified by the UX win and the small table sizes.

- **[Trade-off] Daily-plan dropdown shows just the `date` (and optional notes) — not a friendly"June 3rd" — because the design avoids introducing a date-formatting utility in this change. Acceptable; can be polished later.

## Migration Plan

1. **Backend first:** ship `GET /api/daily-plan-assignments` (List) and `PUT /api/daily-plan-assignments/:id` (Update) with their repo + handler changes. The frontend keeps working against the old endpoints during this step — no rollback needed.
2. **Frontend hooks second:** add `useAllDailyPlanAssignments`, `useUpdateDailyPlanAssignment`, refactor `useCreateDailyPlanAssignment` to no-arg. Compile but don't ship.
3. **Frontend pages third:** update `/planning/assignments/page.tsx`, `/planning/assignments/new/page.tsx`, `/planning/assignments/[id]/page.tsx` (and extract `AssignmentForm`). Compile, manual-test in dev, ship.
4. **Rollback:** revert the frontend bundle. The new backend endpoints are additive (no existing behavior removed), so no backend rollback is required.

## Open Questions

- Does the backend need an explicit audit log entry on assignment create/update/delete (matching the pattern used elsewhere)? The current `Delete` handler doesn't appear to write one. Defer to a follow-up unless the apply-phase review surfaces that other handlers do.
- Should the dropdown for `daily_plan_id` show the plan's technician-names (already assigned) as a secondary hint? Defer — out of scope.
