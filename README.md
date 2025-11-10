# Go Goroutines - Implementasi Lengkap dengan Kasus Nyata

Repository ini berisi implementasi lengkap dari berbagai pola goroutines di Go dengan contoh kasus nyata (real-world use cases).

## 📋 Daftar Isi

1. [Basic Goroutines](#1-basic-goroutines)
2. [Goroutines with Channels](#2-goroutines-with-channels)
3. [WaitGroup](#3-waitgroup)
4. [Worker Pool Pattern](#4-worker-pool-pattern)
5. [Pipeline Pattern](#5-pipeline-pattern)
6. [Fan-out/Fan-in Pattern](#6-fan-outfan-in-pattern)
7. [Select Statement](#7-select-statement)
8. [Context-based Cancellation](#8-context-based-cancellation)
9. [Rate Limiting](#9-rate-limiting)
10. [Error Handling](#10-error-handling)

## 🚀 Cara Menjalankan

### Prerequisites
- Go 1.16 atau lebih baru

### Menjalankan Contoh

Setiap contoh dapat dijalankan secara terpisah:

```bash
# Basic Goroutines
go run examples/basic/main.go

# Channels
go run examples/channels/main.go

# WaitGroup
go run examples/waitgroup/main.go

# Worker Pool
go run examples/workerpool/main.go

# Pipeline
go run examples/pipeline/main.go

# Fan-out/Fan-in
go run examples/fanout/main.go

# Select Statement
go run examples/select/main.go

# Context
go run examples/context/main.go

# Rate Limiting
go run examples/ratelimit/main.go

# Error Handling
go run examples/errorhandling/main.go
```

## 📚 Penjelasan Detail

### 1. Basic Goroutines

**File**: `examples/basic/main.go`

**Konsep**: Eksekusi konkuren sederhana menggunakan keyword `go`.

**Kasus Nyata**: Memproses beberapa file secara bersamaan untuk menghitung jumlah baris.

**Key Points**:
- Goroutine diluncurkan dengan keyword `go`
- Goroutine berjalan secara independen
- Program utama perlu menunggu goroutine selesai

```go
go processFile(filename)
```

### 2. Goroutines with Channels

**File**: `examples/channels/main.go`

**Konsep**: Komunikasi antar goroutines menggunakan channels.

**Kasus Nyata**: Sistem message queue dengan producer-consumer pattern.

**Key Points**:
- Channels untuk komunikasi yang aman antar goroutine
- Buffered vs unbuffered channels
- Close channels setelah selesai mengirim data

```go
jobs := make(chan int, 10) // Buffered channel
jobs <- value              // Send
value := <-jobs            // Receive
close(jobs)                // Close channel
```

### 3. WaitGroup

**File**: `examples/waitgroup/main.go`

**Konsep**: Sinkronisasi multiple goroutines.

**Kasus Nyata**: Memanggil beberapa API secara bersamaan dan menunggu semua selesai.

**Key Points**:
- `WaitGroup` untuk menunggu grup goroutines
- `Add()` sebelum meluncurkan goroutine
- `Done()` ketika goroutine selesai
- `Wait()` untuk menunggu semua selesai

```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work here
}()
wg.Wait()
```

### 4. Worker Pool Pattern

**File**: `examples/workerpool/main.go`

**Konsep**: Membatasi jumlah goroutine yang berjalan bersamaan.

**Kasus Nyata**: Memproses gambar dengan jumlah worker yang tetap untuk menghindari overload.

**Key Points**:
- Jumlah worker yang tetap
- Job queue dengan channel
- Efisien untuk resource management

```go
for i := 1; i <= numWorkers; i++ {
    go worker(i, jobs, results)
}
```

### 5. Pipeline Pattern

**File**: `examples/pipeline/main.go`

**Konsep**: Mengalirkan data melalui beberapa tahap pemrosesan.

**Kasus Nyata**: Memproses log entries melalui pipeline cleaning → parsing → enrichment.

**Key Points**:
- Setiap stage adalah goroutine terpisah
- Data mengalir melalui channels
- Komposisi function yang modular

```go
stage1 := cleanLogs(input)
stage2 := parseLogs(stage1)
stage3 := enrichLogs(stage2)
```

### 6. Fan-out/Fan-in Pattern

**File**: `examples/fanout/main.go`

**Konsep**: Mendistribusikan pekerjaan ke multiple goroutines (fan-out) dan mengumpulkan hasil (fan-in).

**Kasus Nyata**: Menghitung bilangan prima secara terdistribusi dengan multiple workers.

**Key Points**:
- Fan-out: Distribusi kerja ke banyak worker
- Fan-in: Agregasi hasil dari semua worker
- Paralelisme maksimal untuk CPU-intensive tasks

```go
// Fan-out
for i := 1; i <= numWorkers; i++ {
    go worker(i, jobs, results)
}

// Fan-in
go aggregator(results, done)
```

### 7. Select Statement

**File**: `examples/select/main.go`

**Konsep**: Multiplexing beberapa channel operations.

**Kasus Nyata**: Menangani timeout pada operasi I/O dan multiplexing beberapa sumber data.

**Key Points**:
- Select untuk multiple channel operations
- Timeout dengan `time.After()`
- Non-blocking operations

```go
select {
case msg := <-channel1:
    // handle channel1
case msg := <-channel2:
    // handle channel2
case <-time.After(timeout):
    // handle timeout
}
```

### 8. Context-based Cancellation

**File**: `examples/context/main.go`

**Konsep**: Graceful shutdown dan pembatalan goroutines.

**Kasus Nyata**: Menghentikan long-running workers secara graceful saat shutdown.

**Key Points**:
- Context untuk propagasi cancellation
- `WithTimeout` untuk operasi dengan timeout
- `WithCancel` untuk pembatalan manual
- Cleanup resources saat cancelled

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case <-ctx.Done():
    // handle cancellation
case result := <-resultChan:
    // handle result
}
```

### 9. Rate Limiting

**File**: `examples/ratelimit/main.go`

**Konsep**: Membatasi rate eksekusi goroutines.

**Kasus Nyata**: Membatasi jumlah API requests per detik untuk menghindari rate limit.

**Key Points**:
- Token bucket algorithm
- Ticker untuk refill tokens
- Backpressure handling

```go
limiter := NewRateLimiter(requestsPerSecond)
limiter.Wait() // Block until token available
```

### 10. Error Handling

**File**: `examples/errorhandling/main.go`

**Konsep**: Mengumpulkan dan menangani error dari multiple goroutines.

**Kasus Nyata**: Menjalankan multiple tasks dan mengumpulkan semua error yang terjadi.

**Key Points**:
- Error aggregation dengan mutex
- Structured error handling
- Result channel dengan error field

```go
type TaskResult struct {
    Data  string
    Error error
}

// Collect errors safely
errorCollector.Add(err)
```

## 🎯 Best Practices

1. **Selalu tutup channels** yang sudah tidak digunakan untuk menghindari goroutine leak
2. **Gunakan WaitGroup** atau channels untuk sinkronisasi
3. **Hindari share memory**, komunikasi dengan channels (CSP principle)
4. **Gunakan context** untuk cancellation dan timeout
5. **Limit jumlah goroutines** dengan worker pool untuk resource-intensive tasks
6. **Handle errors properly** dalam concurrent code
7. **Avoid goroutine leaks** dengan memastikan semua goroutine bisa terminate

## 🔧 Testing

Untuk menjalankan semua contoh secara berurutan:

```bash
for dir in examples/*/; do
    echo "Running $(basename $dir)..."
    go run "$dir/main.go"
    echo "---"
done
```

## 📖 Referensi

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)
- [Go by Example - Goroutines](https://gobyexample.com/goroutines)

## 🤝 Kontribusi

Kontribusi selalu diterima! Silakan buat pull request untuk menambahkan contoh baru atau memperbaiki yang sudah ada.

## 📝 Lisensi

MIT License - Bebas digunakan untuk pembelajaran dan proyek pribadi.