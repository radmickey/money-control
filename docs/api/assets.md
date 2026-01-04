# Assets API

All endpoints require authentication 🔒

## Asset Types

| Type | Description |
|------|-------------|
| `stock` | Stocks |
| `bond` | Bonds |
| `etf` | ETF |
| `crypto` | Cryptocurrency |
| `commodity` | Commodities |
| `real_estate` | Real Estate |
| `other` | Other |

---

## List Assets

Get all assets for the current user.

**Endpoint:** `GET /assets`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `sub_account_id` | string | Filter by sub-account |
| `type` | string | Filter by asset type |

**Response:**

```json
{
  "success": true,
  "data": {
    "assets": [
      {
        "id": "uuid",
        "user_id": "uuid",
        "sub_account_id": "uuid",
        "symbol": "AAPL",
        "name": "Apple Inc.",
        "type": "stock",
        "quantity": 10.0,
        "purchase_price": 150.00,
        "current_price": 175.00,
        "currency": "USD",
        "total_value": 1750.00,
        "profit_loss": 250.00,
        "profit_loss_percent": 16.67,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-15T00:00:00Z"
      }
    ],
    "total": 5,
    "total_value": 15000.00,
    "total_profit_loss": 2500.00
  }
}
```

---

## Create Asset

Create a new asset.

**Endpoint:** `POST /assets`

**Request Body:**

```json
{
  "sub_account_id": "uuid",
  "symbol": "AAPL",
  "name": "Apple Inc.",
  "type": "stock",
  "quantity": 10.0,
  "purchase_price": 150.00,
  "currency": "USD"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `sub_account_id` | string | No | Sub-account ID |
| `symbol` | string | **Yes** | Asset symbol (e.g., AAPL, BTC) |
| `name` | string | No | Asset name |
| `type` | string | **Yes** | Asset type from list above |
| `quantity` | float | **Yes** | Quantity owned |
| `purchase_price` | float | No | Purchase price per unit |
| `currency` | string | No | Currency code |

**Response:** Created asset object

---

## Get Asset

Get a single asset by ID.

**Endpoint:** `GET /assets/:id`

**Response:** Asset object

---

## Update Asset

Update an existing asset.

**Endpoint:** `PUT /assets/:id`

**Request Body:**

```json
{
  "name": "Updated Name",
  "quantity": 15.0,
  "purchase_price": 145.00,
  "current_price": 180.00
}
```

**Response:** Updated asset object

---

## Delete Asset

Delete an asset.

**Endpoint:** `DELETE /assets/:id`

**Response:** `204 No Content`

---

## Get Asset Price

Get current price for an asset symbol.

**Endpoint:** `GET /assets/price/:symbol`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `type` | string | Asset type (stock, crypto, etc.) |

**Response:**

```json
{
  "success": true,
  "data": {
    "symbol": "AAPL",
    "price": 175.50,
    "currency": "USD",
    "change_24h": 2.50,
    "change_percent_24h": 1.45,
    "updated_at": "2024-01-15T15:30:00Z"
  }
}
```

---

## Search Assets

Search for assets by query.

**Endpoint:** `GET /assets/search`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `q` | string | Search query |
| `type` | string | Filter by asset type |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "symbol": "AAPL",
      "name": "Apple Inc.",
      "type": "stock",
      "exchange": "NASDAQ"
    },
    {
      "symbol": "AMZN",
      "name": "Amazon.com Inc.",
      "type": "stock",
      "exchange": "NASDAQ"
    }
  ]
}
```

