# Dilin əsasları və tiplər

1. **slice vs array fərqi?**

* `array` sabit ölçülü, tipin bir hissəsidir (`[4]int`).
* `slice` dinamik pəncərədir: header (`ptr,len,cap`) + paylaşılan əsas massiv. `append` bəzən yeni massivə kopyalayır.

2. **nil slice vs boş slice (`[]T{}`)?**

* `nil` slice: `len=0, cap=0, ptr=nil`. JSON-da `null` ola bilər.
* boş slice: `len=0, cap=0`, ancaq `ptr!=nil`. JSON-da `[]`. API-lərdə çox vaxt boş slice üstündür.

3. **map-in sıralaması?**

* Təkrarlama sırası **deterministik deyil**. Sort lazımdırsa, açarları `make([]K,0,len(m))` toplayıb sort edin.

4. **map-ə paralel yazı?**

* `fatal error: concurrent map writes`. Mütləq `sync.Mutex` və ya `sync.Map`.

5. **value vs pointer receiver nə vaxt?**

* Kiçik, dəyişməz tiplərdə value; böyük strukturlar, mutasiya, interface method set tələbi olduqda pointer.

6. **interface nil tələsi**

```go
var e *MyErr = nil
var err error = e
fmt.Println(err == nil) // false
```

Çünki dinamik tip var, dəyər `nil`. Yoxlama: `if err != nil`.

7. **defer-in LIFO davranışı və cost-u**

* LIFO çağırış; isti yollarda çox defer performansa təsir edə bilər — bəzən manual `Close()` daha sürətlidir.

8. **copy necə işləyir?**

* `copy(dst, src)` minimum `len(dst), len(src)` qədər kopyalayır. Slice-lar eyni əsas massivi paylaşa bilər — ehtiyatlı olun.

9. **string → \[]byte paylaşımmı edir?**

* Yox, kopyalanır (Go 1.x-də). Zero-alloc üçün `unsafe`/`bytes.Builder`/`strings.Builder`.

10. **panic nə vaxt məntiqlidir?**

* Bərpa olunmaz proqramlaşdırma xətaları (invariant pozuntusu). Biznes xətaları üçün `error`.

# Concurrency & Go Memory Model

1. **goroutine sızıntısı nədir?**

* Heç vaxt dayanmayan (blocked/select-də ilişən) goroutine. Həll: `context` ilə ləğv, kanalı bağlamaq, `done` siqnalı.

2. **unbuffered vs buffered channel?**

* Unbuffered: göndərən və alan sinxronlaşır.
* Buffered: `cap>0`, müəyyən dərəcə kənarlaşdırma/buffer.

3. **channel-ı kim bağlamalıdır?**

* Adətən **göndərən** (producer). Consumer bağlamamalıdır.

4. **`select` və `default`**

* `default` busy-loop yarada bilər. Tətbiq edin: `select { case ...: case <-ctx.Done(): return }`.

5. **for-range pointer tələsi**

```go
for _, v := range arr {
    go func(x int){ fmt.Println(x) }(v) // düz
}
```

Loop dəyişəninin adresini ötürməyin.

6. **WaitGroup düzgün istifadəsi**

* `wg.Add(n)` spawn-dan əvvəl; hər goroutine-də `defer wg.Done()`.

7. **time.After leak**

* Tez-tez çağırılan select-lərdə hər dəfə yeni timer yaradır — `time.NewTimer` + `Stop()/Reset()` istifadə edin.

8. **worker pool dizaynı**

* İşləri `jobs` kanalına, nəticələri `results` kanalına; `close(jobs)` ilə bitiş siqnalı.

9. **atomic vs mutex**

* Sadə sayğac/flag üçün `sync/atomic`; kompleks obyekt üçün `Mutex/RWMutex`.

10. **Go memory model əsas ideya**

* Data race yoxdur → ardıcıllıq zəmanətləri `happens-before` (channel send/recv, Mutex Lock/Unlock).

11. **context nə zaman, necə?**

* Request-scoped lifecycle, deadline/timeout, cancel fan-out. Handler → downstream: `ctx context.Context`.

12. **`select` ilə prioritet**

* Prioritet yoxdur; pseudo-prioritet üçün `try` pattern (iki mərhələli select) və ya buffer ölçüsü.

# Error handling, Test, Tooling

1. **error wrapping və yoxlama**

```go
if err != nil { return fmt.Errorf("save user: %w", err) }
errors.Is(err, os.ErrNotExist)
errors.As(err, &pathErr)
```

2. **sentinel errors vs typed errors**

* Public paketlərdə sentinel (`var ErrX = ...`) stabil API üçündür; daxildə typed errors daha zəngindir.

3. **table-driven tests & subtests**

* `t.Run(name, func(t *testing.T){ ... })` ilə kombinasiya; bir çox input üçün cədvəl.

4. **benchmark & profiling**

* `go test -bench=. -benchmem`; `pprof`: `go tool pprof`, `net/http/pprof` prod-safe yalnız behind flag.

5. **race detector**

* `go test -race`, `-race` ilə bina; performance düşər, amma data race aşkarlar.

6. **gomod & versioning**

* Semantic import versioning (`module v2` → `my/mod/v2`). `replace` yalnız lokal dev üçün.

7. **lint & staticcheck**

* `golangci-lint`, `staticcheck`; CI-də bloklayıcı kimi işlədin.

8. **mocking**

* Kiçik interfeyslər (`Writer`, `Clock`, `Now()`), `httptest`, ya da `moq`/`testify/mock`.

# Generics & Reflection

1. **Generics nə dəyər qatır?**

* Boilerplate azalır (container, set, constraints ilə numeric ops), type-safety.

2. **Type constraints**

```go
type Number interface{ ~int | ~int64 | ~float64 }
func Sum[T Number](xs []T) T { var s T; for _, v := range xs { s += v }; return s }
```

3. **Comparable məhdudiyyəti**

* `map` açarı/`==` tələb edirsə `comparable`. `[]T` comparable deyil.

4. **Any vs interface{}**

* Ekvivalent, amma generics kontekstində `any` daha oxunaqlı.

5. **Reflection nə vaxt?**

* Serialization, ORM, validation; performans/komplike risklərinə görə yalnız məcbur olduqda.

6. **`reflect.Value` setlənməsi**

* Yalnız addressable və export olunan sahələrdə (`CanSet()`).

# HTTP, RPC, Microservices

1. **HTTP server sağlamlıq və shutdown**

```go
srv := &http.Server{ Addr: ":8080", Handler: mux }
go srv.ListenAndServe()
<-ctx.Done()
ctx2,_ := context.WithTimeout(context.Background(), 5*time.Second)
_ = srv.Shutdown(ctx2)
```

2. **Middleware pattern**

* `func(next http.Handler) http.Handler` ilə chain; logging, auth, recovery.

3. **gRPC vs REST**

* gRPC: ikili Protobuf, HTTP/2, bi-directional streaming, sərt müqavilə. REST: sadə, brauzer uyumlu.

4. **idempotency**

* POST əməliyyatları üçün `Idempotency-Key`; retri zamanı ikiqat icranı önləyir.

5. **circuit breaker & bulkhead**

* `gobreaker` + worker pool limitləri; kaskad sıradan çıxmanın qarşısı.

6. **observability**

* `otel` tracing, `promhttp` metrics, strukturlaşdırılmış log (`zerolog`, `zap`).

7. **config & 12-factor**

* Env-first, `viper`/`envconfig`, immutable build, stateless servislər.

# Performance & Memory

1. **Escape analysis**

* Heap/stack qərarı; `-gcflags="-m"` ilə yoxla. Lazımsız heap allocation → GC təzyiqi.

2. **Allocation azaldılması**

* `bytes.Buffer/Builder`, `make([]T,0,cap)`, re-use (sync.Pool), avoid `fmt` isti path-də.

3. **GC tuning**

* `GOGC` hədəf nisbəti (default 100). Telemetriya ilə ölç, erkən optimizasiya etmə.

4. **Large object reuse**

* `sync.Pool` read-heavy yollar üçün faydalıdır, lakin GC hər dövrdə boşalda bilər.

5. **I/O throughput**

* `bufio.Reader/Writer`, `http.Transport` tuningi (connection reuse, idle conns).

# Hiyləgər mini-tapşırıqlar (teoriya yoxlama)

1. **Nə çap edəcək? Niyə?**

```go
arr := []int{1,2,3}
for _, v := range arr {
    defer fmt.Print(v)
}
// Çıxış: 3 2 1  (LIFO defer)
```

2. **Channel ilə deadlock nümunəsi?**

```go
ch := make(chan int)
go func(){ ch <- 1 }()
fmt.Println(<-ch) // ok
fmt.Println(<-ch) // burada deadlock (ikinci dəyər gəlmir)
```

# Qısa real-case ssenarilər (cavab konturları ilə)

## 1) Ödəniş aggregatoru: 3 provayderə paralel müraciət, **overall timeout=5s**, idempotent emal

**Həll xülasəsi:**

* `ctx, cancel := context.WithTimeout(ctx, 5*time.Second)`; `defer cancel()`.
* Hər provayder üçün `goroutine` + `select { case res := <-ch; case <-ctx.Done(): }`.
* `Idempotency-Key` (request UID) + DB-də status tabeli (`PENDING/SUCCEEDED/FAILED`).
* İlk uğurlu cavabı qəbul et (race), digərlərini ləğv et. Bütövlük üçün outbox + txn.

**Sürətli skelet:**

```go
type Result struct{ Name string; Err error }
ch := make(chan Result, 3)
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second); defer cancel()
for _, p := range providers {
    go func(p Provider){
        ch <- callProvider(ctx, p) // ctx deadline/ cancel
    }(p)
}
var final error
for i := 0; i < len(providers); i++ {
    select {
    case r := <-ch:
        if r.Err == nil { final = nil; goto DONE }
    case <-ctx.Done():
        final = ctx.Err(); goto DONE
    }
}
DONE:
// idempotent persist, outbox publish...
return final
```

## 2) Yüksək yüklü read-heavy endpoint

**Həll:** `RWMutex` və ya `atomic.Value` ilə konfiqin periodik yenilənməsi; `http.ServeContent`/`ETag`; `gzip`/`brotli` reverse proxy; `Keep-Alive` tuning.

## 3) Fan-in/fan-out pipeline və backpressure

**Həll:** Producer → `jobs` (buffered) → N worker → `results`; `close(jobs)`; `context` ləğv; `select` ilə `ctx.Done()` bütün mərhələlərdə.

