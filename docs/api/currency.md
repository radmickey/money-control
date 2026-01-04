# Currency API

Public endpoints (no authentication required).

---

## List Currencies

Get list of supported currencies.

**Endpoint:** `GET /currencies`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `include_crypto` | bool | Include cryptocurrencies (default: false) |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "code": "USD",
      "name": "US Dollar",
      "symbol": "$",
      "type": "fiat"
    },
    {
      "code": "EUR",
      "name": "Euro",
      "symbol": "€",
      "type": "fiat"
    },
    {
      "code": "BTC",
      "name": "Bitcoin",
      "symbol": "₿",
      "type": "crypto"
    }
  ]
}
```

---

## Get Exchange Rates

Get all exchange rates for a base currency.

**Endpoint:** `GET /currencies/rates/:base`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `base` | string | Base currency code (e.g., USD, EUR) |

**Response:**

```json
{
  "success": true,
  "data": {
    "base": "USD",
    "rates": {
      "EUR": 0.92,
      "GBP": 0.79,
      "JPY": 149.50,
      "RUB": 89.50
    },
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

---

## Get Single Rate

Get exchange rate between two currencies.

**Endpoint:** `GET /currencies/rates/:from/:to`

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `from` | string | Source currency code |
| `to` | string | Target currency code |

**Response:**

```json
{
  "success": true,
  "data": {
    "rate": 0.92,
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

---

## Convert Amount

Convert an amount between currencies.

**Endpoint:** `POST /currencies/convert`

**Request Body:**

```json
{
  "amount": 100.00,
  "from": "USD",
  "to": "EUR"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `amount` | float | **Yes** | Amount to convert |
| `from` | string | **Yes** | Source currency code |
| `to` | string | **Yes** | Target currency code |

**Response:**

```json
{
  "success": true,
  "data": {
    "original_amount": 100.00,
    "from_currency": "USD",
    "converted_amount": 92.00,
    "to_currency": "EUR",
    "rate_used": 0.92
  }
}
```

