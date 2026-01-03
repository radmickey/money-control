# Authentication API

## Register

Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "base_currency": "USD"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "base_currency": "USD",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

---

## Login

Authenticate with email and password.

**Endpoint:** `POST /auth/login`

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:** Same as Register

---

## Google OAuth

Initiate Google OAuth flow.

**Endpoint:** `GET /auth/google`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `redirect` | string | Set to `false` to get URL as JSON instead of redirect |

**Response (redirect=false):**

```json
{
  "success": true,
  "data": {
    "url": "https://accounts.google.com/o/oauth2/auth?..."
  }
}
```

---

## Google OAuth Callback

Handle Google OAuth callback.

**Endpoint:** `GET /auth/google/callback`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `code` | string | Authorization code from Google |
| `state` | string | State parameter for CSRF protection |

**Response:** Redirects to frontend with tokens

---

## Telegram Auth

Authenticate via Telegram WebApp.

**Endpoint:** `POST /auth/telegram`

**Request Body:**

```json
{
  "init_data": "query_id=...&user=%7B%22id%22%3A...&auth_date=...&hash=..."
}
```

**Response:** Same as Register

---

## Refresh Token

Get new access token using refresh token.

**Endpoint:** `POST /auth/refresh`

**Request Body:**

```json
{
  "refresh_token": "eyJ..."
}
```

**Response:** Same as Register

---

## Get Profile 🔒

Get current user profile.

**Endpoint:** `GET /auth/profile`

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "base_currency": "USD",
    "google_id": "",
    "telegram_id": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

---

## Update Profile 🔒

Update current user profile.

**Endpoint:** `PUT /auth/profile`

**Request Body:**

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "base_currency": "EUR"
}
```

**Response:** Updated user object

---

## Logout 🔒

Revoke refresh token.

**Endpoint:** `POST /auth/logout`

**Request Body:**

```json
{
  "refresh_token": "eyJ..."
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "message": "Successfully logged out"
  }
}
```

