# Mobile Testing Specification

## Purpose

The testing strategy ensures comprehensive test coverage across all layers: backend unit tests (80% target), backend integration tests (critical flows), frontend E2E tests (all new features), mobile unit tests (80% target), mobile integration tests (key flows), and mobile E2E tests (critical paths).

## Requirements

### Requirement: Backend Unit Tests
Every handler method SHALL have at least one test. Every service method SHALL have success + error path tests. Every repository method SHALL have success + not-found + error tests. Use mock-based (testify/mock) for unit tests. Use table-driven tests for multiple scenarios per function.

#### Scenario: Backend unit tests exist
- **WHEN** backend code is written
- **THEN** every handler method has at least one test
- **AND** every service method has success + error path tests
- **AND** every repository method has success + not-found + error tests

### Requirement: Backend Integration Tests
Test against real PostgreSQL (test database). Cover: Photo upload → thumbnail → serve flow, Sync push → pull → conflict detection flow, Visit create → partner resolve → SLA auto-link flow, Checklist populate → stage filter → status block flow, Reopen → modify → complete → audit trail flow.

#### Scenario: Backend integration tests run
- **WHEN** integration tests are executed
- **THEN** photo upload → thumbnail → serve flow is tested
- **AND** sync push → pull → conflict detection flow is tested
- **AND** visit create → partner resolve → SLA auto-link flow is tested

### Requirement: Frontend E2E Tests (Playwright)
Every new feature SHALL have at least one E2E spec. Tests run against live backend + frontend. Follow existing pattern: login → navigate → interact → assert.

#### Scenario: Frontend E2E tests exist
- **WHEN** frontend feature is implemented
- **THEN** at least one E2E spec exists
- **AND** test follows pattern: login → navigate → interact → assert

### Requirement: Mobile Unit Tests (Jest)
Every lib/ module SHALL have tests. Every utils/ module SHALL have tests. Every store SHALL have tests. Mock SQLite for database tests. Mock API for sync tests.

#### Scenario: Mobile unit tests exist
- **WHEN** mobile code is written
- **THEN** every lib/ module has tests
- **AND** every utils/ module has tests
- **AND** every store has tests

### Requirement: Mobile Integration Tests (Jest + MSW)
Mock Service Worker for API mocking. Test auth flow with mocked endpoints. Test sync engine with mocked backend. Test photo upload with mocked backend.

#### Scenario: Mobile integration tests exist
- **WHEN** mobile integration tests are written
- **THEN** auth flow is tested with mocked endpoints
- **AND** sync engine is tested with mocked backend
- **AND** photo upload is tested with mocked backend

### Requirement: Mobile E2E Tests (Detox or Maestro)
Critical user flows only (5-10 scenarios). Cover: Login → visit list → detail → checklist, Photo capture → compress → upload, GPS checkpoint → waiver, Offline → reconnect → sync, Status transitions with checklist blocking.

#### Scenario: Mobile E2E tests exist
- **WHEN** mobile E2E tests are written
- **THEN** login → visit list → detail → checklist is tested
- **AND** photo capture → compress → upload is tested
- **AND** GPS checkpoint → waiver is tested
