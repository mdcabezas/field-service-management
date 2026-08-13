## 1. Backend — Add missing endpoints

- [x] 1.1 Add `List(ctx, limit, offset)` and `Update(ctx, id, a)` methods to `DailyPlanAssignmentRepo`
- [x] 1.2 Add `List` and `Update` handler methods in `daily_plan_assignment_handler.go`
- [x] 1.3 Register `GET /api/daily-plan-assignments` and `PUT /api/daily-plan-assignments/:id` in `router.go`

## 2. Frontend — Hooks refactor

- [x] 2.1 Refactor `useCreateDailyPlanAssignment` to no-arg (drop `dailyPlanId` param, send body as-is)
- [x] 2.2 Add `useAllDailyPlanAssignments()` hook (GET `/api/daily-plan-assignments`)
- [x] 2.3 Add `useUpdateDailyPlanAssignment()` hook (PUT `/api/daily-plan-assignments/:id`)

## 3. Frontend — Assignment form component

- [x] 3.1 Create `AssignmentForm` component with mode `new | edit` and three dropdown selects (daily_plan_id, tech_id, role_id)
- [x] 3.2 Create new/edit form handle empty-state scenarios (no technicians, no roles, no daily plans)

## 4. Frontend — Pages

- [x] 4.1 Fix `/planning/assignments/page.tsx` to use `useAllDailyPlanAssignments()` and display all assignments
- [x] 4.2 Fix `/planning/assignments/new/page.tsx` to use refactored `useCreateDailyPlanAssignment()` and `AssignmentForm`
- [x] 4.3 Fix `/planning/assignments/[id]/page.tsx` to wire edit mode with `useUpdateDailyPlanAssignment()` and `AssignmentForm`

## 5. Verify

- [x] 5.1 Rebuild backend, restart container, and verify frontend compilation
- [x] 5.2 Re-test all assignment flows with playwright-cli in headed mode