# Testing Strategy Spec

## Delta: New capability — comprehensive test coverage across all layers

## Coverage Targets

| Layer | Target | Framework | Tests |
|-------|--------|-----------|-------|
| Backend unit | 80% | Go testing + testify/mock | ~110 |
| Backend integration | Critical flows | Go testing + test DB | ~10 |
| Frontend E2E | All new features | Playwright | ~30 |
| Mobile unit | 80% | Jest + React Native Testing Library | ~50 |
| Mobile integration | Key flows | Jest + MSW | ~10 |
| Mobile E2E | Critical paths | Detox or Maestro | ~10 |

## Functional Requirements

### FR-1: Backend Unit Tests
- Every handler method has at least one test
- Every service method has success + error path tests
- Every repository method has success + not-found + error tests
- Mock-based (testify/mock) for unit tests
- Table-driven tests for multiple scenarios per function

### FR-2: Backend Integration Tests
- Test against real PostgreSQL (test database)
- Photo upload → thumbnail → serve flow
- Sync push → pull → conflict detection flow
- Visit create → partner resolve → SLA auto-link flow
- Checklist populate → stage filter → status block flow
- Reopen → modify → complete → audit trail flow

### FR-3: Frontend E2E Tests (Playwright)
- Every new feature has at least one E2E spec
- Tests run against live backend + frontend
- Follow existing pattern: login → navigate → interact → assert
- Use existing helpers: login, gotoPage, fillInput, fillSelect, submitForm, waitForToast, confirmDialog

### FR-4: Mobile Unit Tests (Jest)
- Every lib/ module has tests
- Every utils/ module has tests
- Every store has tests
- Mock SQLite for database tests
- Mock API for sync tests

### FR-5: Mobile Integration Tests (Jest + MSW)
- Mock Service Worker for API mocking
- Test auth flow with mocked endpoints
- Test sync engine with mocked backend
- Test photo upload with mocked backend

### FR-6: Mobile E2E Tests (Detox or Maestro)
- Critical user flows only (5-10 scenarios)
- Login → visit list → detail → checklist
- Photo capture → compress → upload
- GPS checkpoint → waiver
- Offline → reconnect → sync
- Status transitions with checklist blocking

## Test Naming Convention

### Backend (Go)
```
TestEntityName_MethodName_Scenario
Example: TestVisitService_StartVisit_Success
Example: TestVisitService_StartVisit_InvalidTransition
Example: TestVisitHandler_Create_BodyTooLarge
```

### Frontend (Playwright)
```
test.describe('Feature Name', () => {
  test('action description', async ({ page }) => { ... })
})
Example: test('Create visit with partner auto-resolution')
Example: test('Checklist blocks status when required items incomplete')
```

### Mobile (Jest)
```
describe('moduleName', () => {
  it('should description', () => { ... })
})
Example: it('should compress image to 800px width')
Example: it('should return haversine distance between two points')
```
