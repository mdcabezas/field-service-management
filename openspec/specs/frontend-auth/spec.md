## Purpose

Maneja la autenticación LDAP del usuario, el almacenamiento seguro de JWT (access + refresh tokens), la renovación automática de tokens y el cierre de sesión.

## ADDED Requirements

### Requirement: Login form
The system SHALL display a login form with fields for employee_number and password.

#### Scenario: Successful login
- **WHEN** user submits valid employee_number and password
- **THEN** system sends POST to /auth/login with credentials
- **AND** stores access_token and refresh_token in httpOnly cookies
- **AND** redirects user to /dashboard

#### Scenario: Failed login
- **WHEN** user submits invalid credentials
- **THEN** system displays error message "Credenciales inválidas"
- **AND** form remains editable for retry

### Requirement: JWT token storage
The system SHALL store JWT tokens in httpOnly cookies, not in localStorage or sessionStorage.

#### Scenario: Token persistence
- **WHEN** user successfully logs in
- **THEN** access_token cookie is set with 1 hour expiry
- **AND** refresh_token cookie is set with 7 day expiry
- **AND** cookies are httpOnly, secure, and sameSite=strict

### Requirement: Automatic token refresh
The system SHALL automatically refresh the access token before expiry.

#### Scenario: Token refresh triggered
- **WHEN** access_token is within 5 minutes of expiry
- **THEN** system sends POST to /auth/refresh with refresh_token
- **AND** stores new access_token and refresh_token
- **AND** retries failed request with new token

#### Scenario: Refresh token expired
- **WHEN** refresh_token is expired or invalid
- **THEN** system clears all cookies
- **AND** redirects user to /login

### Requirement: Logout
The system SHALL revoke both tokens and clear cookies on logout.

#### Scenario: User clicks logout
- **WHEN** user clicks logout button
- **THEN** system sends POST to /auth/logout with refresh_token
- **AND** clears access_token and refresh_token cookies
- **AND** redirects to /login

### Requirement: Route protection
The system SHALL redirect unauthenticated users to login for protected routes.

#### Scenario: Accessing protected route without token
- **WHEN** user navigates to any /dashboard/* route without valid access_token
- **THEN** system redirects to /login
- **AND** stores intended destination for redirect after login

#### Scenario: Accessing login page while authenticated
- **WHEN** user navigates to /login with valid access_token
- **THEN** system redirects to /dashboard

### Requirement: Current user info
The system SHALL display the current user's employee number and role.

#### Scenario: Fetch user info
- **WHEN** user is authenticated
- **THEN** system sends GET to /auth/me
- **AND** stores employee_number and role in auth store
