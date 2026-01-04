# Insights API

All endpoints require authentication 🔒

---

## Get Net Worth

Get current net worth with currency conversion.

**Endpoint:** `GET /insights/net-worth`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `baseCurrency` | string | Base currency for conversion (default: USD) |
| `currency` | string | Alias for baseCurrency |

**Response:**

```json
{
  "success": true,
  "data": {
    "total": 150000.00,
    "currency": "USD",
    "change24h": 500.00,
    "changePercent24h": 0.33,
    "change7d": 2000.00,
    "changePercent7d": 1.35,
    "change30d": 5000.00,
    "changePercent30d": 3.45,
    "calculatedAt": "2024-01-15T12:00:00Z"
  }
}
```

---

## Get Trends

Get net worth trends over time.

**Endpoint:** `GET /insights/trends`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `currency` | string | Base currency |
| `period` | string | Time period: `7d`, `30d`, `90d`, `1y` (default: 30d) |

**Response:**

```json
{
  "success": true,
  "data": {
    "period": "30d",
    "currency": "USD",
    "data_points": [
      {
        "date": "2024-01-01",
        "value": 145000.00
      },
      {
        "date": "2024-01-15",
        "value": 150000.00
      }
    ]
  }
}
```

---

## Get Allocation

Get asset allocation breakdown.

**Endpoint:** `GET /insights/allocation`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `baseCurrency` | string | Base currency for conversion |
| `currency` | string | Alias for baseCurrency |
| `group_by` | string | Group by: `asset_type`, `account`, `currency` (default: asset_type) |

**Response:**

```json
{
  "success": true,
  "data": {
    "allocations": [
      {
        "category": "Stocks",
        "value": 80000.00,
        "percentage": 53.33,
        "currency": "USD"
      },
      {
        "category": "Crypto",
        "value": 30000.00,
        "percentage": 20.00,
        "currency": "USD"
      },
      {
        "category": "Cash",
        "value": 40000.00,
        "percentage": 26.67,
        "currency": "USD"
      }
    ],
    "total": 150000.00,
    "currency": "USD"
  }
}
```

---

## Get Dashboard Summary

Get complete dashboard summary.

**Endpoint:** `GET /insights/dashboard`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `currency` | string | Base currency |

**Response:**

```json
{
  "success": true,
  "data": {
    "net_worth": {
      "total": 150000.00,
      "change_24h": 500.00,
      "change_percent_24h": 0.33
    },
    "allocation": [...],
    "recent_transactions": [...],
    "top_performers": [...]
  }
}
```

---

## Get Cash Flow

Get cash flow analysis.

**Endpoint:** `GET /insights/cash-flow`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `currency` | string | Base currency |
| `period` | string | Period: `daily`, `weekly`, `monthly` (default: monthly) |
| `start_date` | string | Start date (YYYY-MM-DD) |
| `end_date` | string | End date (YYYY-MM-DD) |

**Response:**

```json
{
  "success": true,
  "data": {
    "period": "monthly",
    "currency": "USD",
    "data": [
      {
        "date": "2024-01",
        "income": 5000.00,
        "expense": 3000.00,
        "net": 2000.00
      }
    ],
    "total_income": 5000.00,
    "total_expense": 3000.00,
    "net_cash_flow": 2000.00
  }
}
```

---

## Get Net Worth History

Get historical net worth data.

**Endpoint:** `GET /insights/history`

**Query Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `currency` | string | Base currency |
| `period` | string | Period: `7d`, `30d`, `90d`, `1y` |
| `start_date` | string | Start date (YYYY-MM-DD) |
| `end_date` | string | End date (YYYY-MM-DD) |

**Response:**

```json
{
  "success": true,
  "data": {
    "history": [
      {
        "date": "2024-01-01T00:00:00Z",
        "value": 145000.00,
        "currency": "USD"
      }
    ]
  }
}
```

