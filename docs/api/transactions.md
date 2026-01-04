# Transactions API

All endpoints require authentication 🔒

## Transaction Types

| Type | Description |
|------|-------------|
| `income` | Income |
| `expense` | Expense |
| `transfer` | Transfer between accounts |

## Transaction Categories

| Category | Description |
|----------|-------------|
| `salary` | Salary |
| `food` | Food & Dining |
| `transport` | Transportation |
| `utilities` | Utilities |
| `entertainment` | Entertainment |
| `shopping` | Shopping |
| `health` | Healthcare |
| `education` | Education |
| `investment` | Investment |
| `transfer` | Transfer |
| `other` | Other |

---

## List Transactions

Get all transactions for the current user.

**Endpoint:** `GET /transactions`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `sub_account_id` | string | Filter by sub-account |

**Response:**

```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "id": "uuid",
        "user_id": "uuid",
        "sub_account_id": "uuid",
        "amount": 1500.00,
        "currency": "USD",
        "type": "income",
        "category": "salary",
        "custom_category": "",
        "description": "Monthly salary",
        "merchant": "Employer Inc",
        "date": "2024-01-15T00:00:00Z",
        "created_at": "2024-01-15T10:00:00Z",
        "updated_at": "2024-01-15T10:00:00Z"
      }
    ],
    "total": 25
  }
}
```

---

## Create Transaction

Create a new transaction.

**Endpoint:** `POST /transactions`

**Request Body:**

```json
{
  "sub_account_id": "uuid",
  "amount": 50.00,
  "currency": "USD",
  "type": "expense",
  "category": "food",
  "description": "Lunch",
  "merchant": "Restaurant",
  "date": "2024-01-15"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `sub_account_id` | string | No | Sub-account ID |
| `amount` | float | **Yes** | Transaction amount |
| `currency` | string | No | Currency code (default: USD) |
| `type` | string | **Yes** | `income`, `expense`, or `transfer` |
| `category` | string | No | Category from list above |
| `custom_category` | string | No | Custom category name |
| `description` | string | No | Description |
| `merchant` | string | No | Merchant name |
| `date` | string | No | Date (YYYY-MM-DD) |

**Response:** Created transaction object

---

## Get Transaction

Get a single transaction by ID.

**Endpoint:** `GET /transactions/:id`

**Response:** Transaction object

---

## Update Transaction

Update an existing transaction.

**Endpoint:** `PUT /transactions/:id`

**Request Body:**

```json
{
  "amount": 75.00,
  "type": "expense",
  "category": "entertainment",
  "description": "Updated description",
  "date": "2024-01-16",
  "currency": "USD"
}
```

**Response:** Updated transaction object

---

## Delete Transaction

Delete a transaction.

**Endpoint:** `DELETE /transactions/:id`

**Response:** `204 No Content`

---

## Get Summary

Get transaction summary for a period.

**Endpoint:** `GET /transactions/summary`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `start_date` | string | Start date (YYYY-MM-DD) |
| `end_date` | string | End date (YYYY-MM-DD) |
| `currency` | string | Base currency |

**Response:**

```json
{
  "success": true,
  "data": {
    "total_income": 5000.00,
    "total_expense": 2500.00,
    "net": 2500.00,
    "currency": "USD",
    "by_category": {
      "salary": 5000.00,
      "food": 500.00,
      "transport": 200.00
    }
  }
}
```

