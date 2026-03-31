
# 🔗 Go URL Shortener — Production-Ready Backend

A high-performance, production-style URL shortener built in **Go**, designed with real-world backend principles like **concurrency control, caching, batching, and rate limiting**.

🚀 **Live Demo:**
👉 https://url-shortner-go-mwmq.onrender.com

---

# ✨ Features

### ⚡ Core Functionality

* Shorten long URLs into compact codes
* Fast redirection using cache-first lookup
* Analytics endpoint for click tracking

---

### 🧠 Advanced Backend Design

* **Cache-aside architecture**
* **Batched DB writes** for counter updates
* **Per-IP rate limiting** with fine-grained locking
* **Concurrent-safe data structures (RWMutex + per-key locks)**

---

### 🌐 Frontend

* Clean UI using **HTML + Tailwind (no frameworks)**
* Copy-to-clipboard support
* Stats page with live analytics
* Toast notifications for feedback

---

### ☁️ Production Deployment

* Deployed on Render
* PostgreSQL (Neon)
* Environment-based configuration
* Publicly accessible API

---

# 🏗️ Architecture Overview

```
Client
   ↓
Middleware (Rate Limiting)
   ↓
Handlers (API Logic)
   ↓
Cache Layer (Read/Write)
   ↓
Database (PostgreSQL via sqlc)
```

---

# ⚙️ Tech Stack

* **Language:** Go
* **Database:** PostgreSQL (Neon)
* **Query Layer:** sqlc
* **Frontend:** HTML + Tailwind
* **Deployment:** Render
* **Concurrency:** Goroutines + Mutexes

---

# 📂 Project Structure

```
url-shortener/
├── cmd/server/main.go      # Entry point
├── internal/
│   ├── cache/             # Cache + batching logic
│   ├── handlers/          # API handlers
│   ├── middleware/        # Rate limiting
│   └── db/
│       ├── queries/       # sqlc generated queries
│       └── migrations/    # Goose migration files
├── static/
│   ├── index.html
│   └── stats.html
├── go.mod
```

---

# 🚀 API Endpoints

## 🔗 Shorten URL

```http
POST /api/shorten
```

### Request

```json
{
  "url": "https://example.com"
}
```

### Response

```json
{
  "short_url": "https://your-app/b"
}
```

---

## 🔁 Redirect

```http
GET /{code}
```

* Redirects to original URL
* Uses cache-first lookup
* Increments click counter

---

## 📊 Stats

```http
GET /api/stats/{code}
```

### Response

```json
{
  "original_url": "...",
  "counter": 174,
  "created_at": "..."
}
```

* Combines DB + in-memory counters
* Works even when DB is slightly behind

---

# 🧠 Key Design Decisions

## 1. Cache-Aside Pattern

* Reads:

  * Cache → DB fallback
* Writes:

  * Cache immediately
  * DB eventually (batched)

---

## 2. Batched Counter Updates

Instead of updating DB on every request:

* Counters stored in memory
* Flushed every **30 seconds**
* Failed writes retried (merge-back logic)

```go
c.CounterM[k] += v
```

---

## 3. Rate Limiting (Per IP)

* Global map: `IP → limiter`
* Per-IP mutex (avoids global contention)
* Sliding window (1 minute)

```go
if value.counter >= limit {
    return false
}
```

---

## 4. Concurrency Safety

* `RWMutex` for cache reads/writes
* Per-IP locking for rate limiter
* No DB calls inside locks

---

# ⚡ Performance Optimizations

| Feature           | Benefit           |
| ----------------- | ----------------- |
| Cache             | Faster reads      |
| Batching          | Reduced DB load   |
| Mutex granularity | Low contention    |
| Rate limiting     | System protection |

---

# 🧪 Load Testing

Tested using `hey`:

```bash
hey -n 1000 -c 100 http://localhost:8080/b
```

### Observed:

* ~150 allowed requests (limit)
* ~850 rejected (429)
* System remained stable under load

---

# 🌱 Environment Variables

```env
DB_URL=postgresql://... (Neon pooler)
APP_URL=https://url-shortner-go-mwmq.onrender.com
PORT=provided by Render
```
### Cold Starts (Render Free Tier)

* App sleeps after inactivity
* First request may take ~10–30 seconds

---

# 🚀 How to Run Locally

```bash
git clone https://github.com/Arush71/url-shortener
cd url-shortener

go run ./cmd/server
```

---

# 🧭 Future Improvements

* Custom domains
* Redis-backed distributed cache
* Persistent rate limiter
* Analytics dashboard
* Authentication layer

---

## ⭐ If you found this useful, consider starring the repo!

