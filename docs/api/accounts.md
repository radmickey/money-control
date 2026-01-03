# Accounts API

All endpoints require authentication 🔒

## Account Types

| Type | Description |
|------|-------------|
| `bank` | Bank account |
| `cash` | Cash |
| `investment` | Investment account |
| `crypto` | Cryptocurrency wallet |
| `real_estate` | Real estate |
| `other` | Other |

---

## List Accounts

Get all accounts for the current user.

**Endpoint:** `GET /accounts`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `type` | string | Filter by account type |
| `page` | int | Page number |
| `page_size` | int | Items per page |

**Response:**

```json
{
  "success": true,
  "data": {
    "accounts": [
      {
        "id": "uuid",
        "user_id": "uuid",
        "name": "Checking Account",
        "type": "bank",
        "currency": "USD",
        "total_balance": 5000.00,
        "description": "Main checking",
        "icon": "🏦",
        "sub_accounts": [
          {
            "id": "uuid",
            "name": "Savings",
            "balance": 2500.00,
            "currency": "USD"
          }
        ],
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20,
    "total_pages": 1
  }
}
```

---

## Create Account

Create a new account.

**Endpoint:** `POST /accounts`

**Request Body:**

```json
{
  "name": "Checking Account",
  "type": "bank",
  "currency": "USD",
  "description": "Main checking account",
  "icon": "🏦"
}
```

**Response:** Created account object

---

## Get Account

Get a single account by ID.

**Endpoint:** `GET /accounts/:id`

**Response:** Account object

---

## Update Account

Update an existing account.

**Endpoint:** `PUT /accounts/:id`

**Request Body:**

```json
{
  "name": "Updated Name",
  "description": "Updated description"
}
```

**Response:** Updated account object

---

## Delete Account

Delete an account.

**Endpoint:** `DELETE /accounts/:id`

**Response:**

```json
{
  "success": true,
  "data": {
    "message": "Account deleted successfully"
  }
}
```

---

## Create Sub-Account

Create a sub-account within an account.

**Endpoint:** `POST /accounts/:id/sub-accounts`

**Request Body:**

```json
{
  "name": "Savings",
  "currency": "USD",
  "balance": 1000.00,
  "asset_type": "cash"
}
```

**Response:** Created sub-account object

---

## List Sub-Accounts

Get all sub-accounts for an account.

**Endpoint:** `GET /accounts/:id/sub-accounts`

**Response:**

```json
{
  "success": true,
  "data": {
    "sub_accounts": [...]
  }
}
```

---

## Update Sub-Account

Update a sub-account.

**Endpoint:** `PUT /sub-accounts/:id`

**Request Body:**

```json
{
  "name": "Updated Name",
  "balance": 2000.00
}
```

---

## Delete Sub-Account

Delete a sub-account.

**Endpoint:** `DELETE /sub-accounts/:id`

---

## Update Sub-Account Balance

Update only the balance of a sub-account.

**Endpoint:** `PATCH /sub-accounts/:id/balance`

**Request Body:**

```json
{
  "balance": 3000.00
}
```

---

## Get Net Worth

Get total net worth across all accounts.

**Endpoint:** `GET /net-worth`

**Response:**

```json
{
  "success": true,
  "data": {
    "total": 150000.00,
    "currency": "USD",
    "by_type": {
      "bank": 50000.00,
      "investment": 80000.00,
      "crypto": 20000.00
    }
  }
}
```

