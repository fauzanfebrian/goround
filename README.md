# Go-Round (`/bin/goround`)

* **Challenger:** Gemini 3.0
* **Target:** Go 1.22+ (Standard Library Only)
* **Difficulty:** Systems Programming Entry

## 🎯 Objective

Construct a Layer 7 (HTTP) Reverse Proxy and Load Balancer from scratch. The application must listen on a single port, accept incoming HTTP requests, and distribute them across a pool of backend servers using a **Round-Robin** algorithm.

Crucially, the system must be **Resilient**: it must detect when a backend server goes down (Health Check) and automatically stop routing traffic to it until it recovers.

This enforces the "Holy Trinity" of Go Systems Programming:

1. **Concurrency Patterns** (Goroutines & Channels vs. Mutexes).
2. **Shared State Management** (`sync/atomic` for counters, `sync.RWMutex` for pool status).
3. **Network Primitives** (`net/http`, `httputil`).

## 🛠️ Specifications

| Feature | Requirement |
| --- | --- |
| **Binary Name** | `goround` |
| **Listener Port** | `:8000` (The Entrypoint) |
| **Backends** | 3+ local servers (e.g., `:8081`, `:8082`, `:8083`) |
| **Algorithm** | Round-Robin (Cyclic distribution) |
| **Health Check** | Passive (Background Pings every 10s) |
| **Dependencies** | **Standard Library ONLY** (No Gin, Chi, or Viper) |

### 1. The `Proxy` Operation (The Spine)

* **Input:** Any HTTP Request to `localhost:8000`.
* **Behavior:**
1. Select the next available Backend from the pool.
2. Update the request Header (add `X-Forwarded-For`).
3. Forward the request to the selected Backend using `httputil.NewSingleHostReverseProxy`.
4. Return the Backend's response to the User.


* **Constraint:** Must handle high concurrency without race conditions on the backend selection index.

### 2. The `Round-Robin` Logic (The Brain)

* **Mechanism:** An incrementing counter (0, 1, 2, 0, 1...).
* **Behavior:**
1. On every request, increment a global counter.
2. Calculate `index = counter % length_of_pool`.
3. **Critical Check:** If the Backend at `index` is marked "Dead" (`Alive == false`), recursively try the next index until a live one is found or loop fails.


* **Safety:** Use `sync/atomic` for the counter to prevent race conditions during high load.

### 3. The `Health Check` Routine (The Heart)

* **Trigger:** Runs concurrently in a separate Goroutine (`go func()`).
* **Behavior:**
1. Every 10 seconds, loop through all Backends.
2. Attempt a TCP Dial or HTTP GET to the Backend.
3. **Transition:**
* If Success + Previous Status was Dead -> Mark **Alive**.
* If Fail + Previous Status was Alive -> Mark **Dead**.


4. **Locking:** Use `sync.RWMutex` when updating the `Alive` status to ensure the Proxy (Reader) doesn't read garbage data while the Health Checker (Writer) is updating it.



## 🤖 AI Usage Protocol (The Honor Code)

| Action | Status | Why? |
| --- | --- | --- |
| **Boilerplate Generation** | ✅ **ALLOWED** | You don't need to memorize how to type `http.ListenAndServe`. Ask AI to set up the skeleton. |
| **Syntax Lookup** | ✅ **ALLOWED** | "How do I parse a URL string in Go?" or "What is the signature for `ReverseProxy`?" |
| **Logic Generation** | ❌ **FORBIDDEN** | **Do not** prompt: *"Write a round-robin algorithm in Go."* You must write the `index % length` logic yourself. |
| **Concurrency Debugging** | ❌ **FORBIDDEN** | **Do not** paste your race condition error and ask for a fix. You must use `go run -race` and read the stack trace yourself. |
| **Architectural Decisions** | ❌ **FORBIDDEN** | **Do not** ask: *"How should I structure the struct for the backend?"* Figure out where the `Mutex` lives on your own. |

## 🧪 Verification Protocol (The Test)

The system is successful only if it passes the "Chaos Monkey" test:

```bash
# 1. Start the Dumb Backends (Terminal 1)
# (You'll need a tiny helper script to spin these up)
go run dummy_servers.go 
# Output: Listening on 8081, 8082, 8083...

# 2. Start the Load Balancer (Terminal 2)
go run main.go
# Output: Load Balancer started at :8000

# 3. The Happy Path (Terminal 3)
for i in {1..4}; do curl localhost:8000; echo; done
# Output:
# Hello from Server 8081
# Hello from Server 8082
# Hello from Server 8083
# Hello from Server 8081 (Cycle repeats)

# 4. The Kill (Terminal 1)
# Manually kill/stop the server on port 8082.

# 5. The Resilience Check (Terminal 3)
for i in {1..4}; do curl localhost:8000; echo; done
# Output:
# Hello from Server 8081
# Hello from Server 8083 (8082 is skipped!)
# Hello from Server 8081
# Hello from Server 8083

```

## 🚫 Constraints & Rules

1. **No Frameworks:** You cannot use Gin, Echo, or specialized LB libraries. You must wrap `net/http` yourself.
2. **No Database:** All state (Backend status) must live in memory (`structs`).
3. **The "Atomic" Rule:** You are forbidden from using a standard `int` for the request counter. You must use `atomic.AddUint64` to understand why simple increments are not thread-safe.

## 📚 Required Resources

* **Reverse Proxy:** `net/http/httputil` -> `NewSingleHostReverseProxy`
* **URLs:** `net/url` -> `Parse`
* **Concurrency:** `sync/atomic` -> `AddUint64`, `sync` -> `RWMutex`
* **Networking:** `net` -> `DialTimeout` (for health checks)
