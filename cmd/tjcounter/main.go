package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kbalashoff/tjcounter/internal/counter"
)

type broker struct {
	mu      sync.Mutex
	clients map[chan int]struct{}
}

func newBroker() *broker {
	return &broker{clients: make(map[chan int]struct{})}
}

func (b *broker) subscribe() chan int {
	ch := make(chan int, 8)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *broker) unsubscribe(ch chan int) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
	close(ch)
}

func (b *broker) publish(v int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		select {
		case ch <- v:
		default:
		}
	}
}

func main() {
	ctr := counter.New()
	events := newBroker()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("web/index.html")
		if err != nil {
			http.Error(w, "failed to load page", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(data)
	})

	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]int{"value": ctr.Get()})
	})

	mux.HandleFunc("/api/increment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		v := ctr.Increment()
		events.publish(v)
		writeJSON(w, map[string]int{"value": v})
	})

	mux.HandleFunc("/api/decrement", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		v := ctr.Decrement()
		events.publish(v)
		writeJSON(w, map[string]int{"value": v})
	})

	mux.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		v := ctr.Reset()
		events.publish(v)
		writeJSON(w, map[string]int{"value": v})
	})

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		ch := events.subscribe()
		defer events.unsubscribe(ch)

		fmt.Fprintf(w, "event: value\ndata: %d\n\n", ctr.Get())
		flusher.Flush()

		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-ticker.C:
				fmt.Fprint(w, ": keep-alive\n\n")
				flusher.Flush()
			case v := <-ch:
				fmt.Fprintf(w, "event: value\ndata: %d\n\n", v)
				flusher.Flush()
			}
		}
	})

	addr := ":8080"
	log.Printf("tjcounter listening on %s", addr)
	if err := http.ListenAndServe(addr, loggingMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
