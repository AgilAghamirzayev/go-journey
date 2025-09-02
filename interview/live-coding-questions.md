# 1) Top-K tezlik (map + heap/partial sort) — 10 dəq

**Məqsəd:** Verilən `[]string` içində ən çox təkrarlanan **K** sözü tap.
**İmza:** `func TopK(words []string, k int) []string`
**Şərtlər:** Stabil qaytarma lazım deyil; eyni tezlikdə olanların sırası fərq etmir.
**Edge:** `k<=0`, boş input.
**Baxılacaq:** `map` + min-heap və ya `sort.Slice`.

# 2) Order-preserving dedupe — 7 dəq

**Məqsəd:** Dəyərlərin ilk göründüyü sıranı saxlayaraq təkrarlıları sil.
**İmza:** `func Unique[T comparable](in []T) []T`
**Şərtlər:** Generics. O(1) əlavə yaddaş olmasa da olar.
**Edge:** `nil` vs `[]T{}`.

# 3) Two-Sum (linear) — 5 dəq

**Məqsəd:** Cəmi `target` olan iki indeks qaytar.
**İmza:** `func TwoSum(nums []int, target int) (i, j int, ok bool)`
**Baxılacaq:** HashMap, O(n).

# 4) Fan-in: K kanalı birləşdir (context-lə ləğv) — 12 dəq

**Məqsəd:** Bir neçə `(<-chan T)` axınını **təhlükəsiz** şəkildə birləşdir.
**İmza:** `func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T`
**Şərtlər:** `ctx.Done()` bağlananda çıxış kanalı **bağlanmalıdır**. Goroutine sızması olmamalıdır.
**Baxılacaq:** `sync.WaitGroup`, `defer`, close(…).

# 5) Bounded Worker Pool — 15 dəq

**Məqsəd:** N ədəd işi M worker ilə icra et, nəticələri qaytar.
**İmza:**

```go
type Task func(context.Context) (any, error)
func RunPool(ctx context.Context, tasks []Task, workers int) ([]any, error)
```

**Şərtlər:** İlk **error** gələndə hamısını **ləğv et** (ctx cancel).
**Baxılacaq:** backpressure (buffered jobs/results), error aggregation, cancel fan-out.

# 6) “First-success wins” (race) — 10 dəq

**Məqsəd:** 3 provayderə paralel çağır, ilk `nil error` alan kimi qalanları ləğv et.
**İmza:** `func FirstOK(ctx context.Context, calls ...func(context.Context) error) error`
**Baxılacaq:** `select` ilə nəticə/`ctx.Done()`, `errgroup` və ya öz kanal pattern.

# 7) Sliding-window rate limiter (Goroutine-safe) — 15 dəq

**Məqsəd:** Son 1 saniyədə **N** icazə.
**İmza:** `type Limiter struct { ... }` + `func NewLimiter(n int, window time.Duration) *Limiter` + `Allow() bool`
**Şərtlər:** Concurrency-safe.
**Baxılacaq:** `sync.Mutex`/`atomic.Value`, köhnə timestamp-lərin təmizlənməsi.

# 8) JSON stream processing (O(1) peak memory) — 12 dəq

**Məqsəd:** NDJSON (hər sətirdə JSON obyekt) faylından `amount` cəmini hesabla.
**İmza:** `func SumAmounts(r io.Reader) (int64, error)`
**Şərtlər:** Bütün faylı RAM-a yükləmə! `json.Decoder` + `UseNumber`.
**Edge:** Pis satırları skip et və say (opsiyonel).

# 9) Graceful HTTP server + shutdown — 12 dəq

**Məqsəd:** `/healthz` (200), `/sum` (POST JSON: `nums: []int`) qaytarır cəm. Graceful shutdown.
**Skeleton:**

```go
func main() {
  mux := http.NewServeMux()
  // handlers...
  srv := &http.Server{Addr: ":8080", Handler: mux}
  go func(){ _ = srv.ListenAndServe() }()
  // os.Signal catch → ctx timeout → srv.Shutdown(ctx)
}
```

**Baxılacaq:** `Shutdown` vs `Close`, request-scoped `ctx`.

# 10) LRU Cache (Generics) — 20 dəq

**Məqsəd:** `Get/Put`, O(1) əməliyyatlar, ən az istifadə olunanı sil.
**İmza:**

```go
type LRU[K comparable, V any] struct { /* map + list */ }
func NewLRU[K comparable, V any](cap int) *LRU[K,V]
func (c *LRU[K,V]) Get(k K) (V, bool)
func (c *LRU[K,V]) Put(k K, v V)
```

**Baxılacaq:** `container/list`, map + doubly-linked list sinxronu.

# 11) Parallel Map with error aggregation — 12 dəq

**Məqsəd:** Elementləri paralel emal et, səhvləri yığ və qaytar.
**İmza:**

```go
func PMap[T any, R any](ctx context.Context, in []T, fn func(context.Context, T) (R, error), workers int) ([]R, error)
```

**Şərtlər:** Sıra qorunub saxlanılsın (çıxış indekslənmiş). İlk error → cancel.

# 12) Deadlock/race tapan mini-fix — 7 dəq

**Kod veriləcək:** Map-ə paralel yazı və ya kanal bağlanmasını consumer edir.
**Tapşırıq:** Səbəbi izah et və **yerdəcə** düzəlt (`Mutex`/`RWMutex`, “sender closes channel” qaydası).

# 13) Timeout-safe retrier (exponential backoff) — 10 dəq

**Məqsəd:** `op()` funksiyasını `maxAttempts`, `baseDelay` ilə retry et.
**İmza:**

```go
func Retry(ctx context.Context, max int, base time.Duration, op func(context.Context) error) error
```

**Şərtlər:** `ctx` bitirsə — dayandır. Jitter (opsiyonel).

# 14) CSV → Grouped stats (stream, no alloc spikes) — 12 dəq

**Məqsəd:** `user_id,amount` csv-dən istifadəçi üzrə cəmlər.
**İmza:** `func GroupSum(r io.Reader) (map[string]int64, error)`
**Baxılacaq:** `csv.NewReader`, səhv satırları idarə, `ReadAll` **yox**.

# 15) Safe cancellable timers — 8 dəq

**Məqsəd:** `time.After` leak-lərini aradan qaldır.
**İmza:**

```go
func WaitWithTimeout(ctx context.Context, d time.Duration, fn func() error) error
```

**Şərtlər:** `time.NewTimer` + `Stop/Drain` düzgün istifadəsi.

# 16) Pipeline: parse → validate → persist — 18 dəq

**Məqsəd:** 3 mərhələli kanal boru-xətti.
**İmza:** `func RunPipeline(ctx context.Context, in <-chan string, db Saver) error`
**Şərtlər:** Hər mərhələ `ctx.Done()` dinləsin; kanallar vaxtında `close`.
**Baxılacaq:** backpressure, sızma yox.

# 17) HTTP client pool & timeouts — 12 dəq

**Məqsəd:** Eyni hosta çoxlu GET, `Transport` tuningi.
**İmza:** `func FetchAll(ctx context.Context, urls []string) ([][]byte, error)`
**Şərtlər:** `http.Client{Timeout}`, `Transport{MaxIdleConnsPerHost}`, context-aware request.

# 18) MinBy / MaxBy (Generics utility) — 7 dəq

**Məqsəd:** Komparatorla min/max element.
**İmza:**

```go
func MinBy[T any](xs []T, less func(a,b T) bool) (T, bool)
```

**Edge:** Boş dilim — `ok=false`.

# 19) Safe file rotate reader (opsiyonel çətin) — 20 dəq

**Məqsəd:** Log faylını “tail” et, rotate olunsa belə oxumağa davam et.
**İmza:** `func Tail(ctx context.Context, path string, out chan<- []byte) error`
**Baxılacaq:** inode dəyişimini aşkarlama, `os.Open` yeniləmə.

# 20) Table-driven tests + race detector — 10 dəq

**Tapşırıq:** Yuxarıdakı 1-2 funksiyaya `*_test.go` yaz, `t.Run` subtests, `go test -race` keçsin.
**Baxılacaq:** Müsbət/neqativ halları əhatə, deterministik test.

---

## Mini Skeletlər (başlamaq üçün)

**Fan-in**

```go
func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(chans))
	for _, ch := range chans {
		ch := ch
		go func() {
			defer wg.Done()
			for {
				select {
				case v, ok := <-ch:
					if !ok { return }
					select {
					case out <- v:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
```

**Bounded pool (skelet)**

```go
func RunPool(ctx context.Context, tasks []Task, workers int) ([]any, error) {
	type item struct{ idx int; res any; err error }
	jobs := make(chan int)
	out  := make(chan item)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for i := range jobs {
			r, err := tasks[i](ctx)
			select {
			case out <- item{i, r, err}:
			case <-ctx.Done():
				return
			}
		}
	}
	wg.Add(workers)
	for w := 0; w < workers; w++ { go worker() }

	go func() { // feed jobs
		for i := range tasks {
			select {
			case jobs <- i:
			case <-ctx.Done():
				close(jobs); return
			}
		}
		close(jobs)
	}()

	results := make([]any, len(tasks))
	var firstErr error
	done := 0
	for done < len(tasks) {
		select {
		case it := <-out:
			if it.err != nil && firstErr == nil {
				firstErr = it.err
				cancel()
			}
			results[it.idx] = it.res
			done++
		case <-ctx.Done():
			// drain optionally
			wg.Wait()
			return results, firstErr
		}
	}
	wg.Wait()
	return results, firstErr
}
```