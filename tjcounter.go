// TJ Counter
//
// Run this code like:
//  > go run tjcounter.go
//
// Then open up your browser to http://localhost:8181

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Broker struct {
	clients        map[chan string]bool
	newClients     chan chan string
	defunctClients chan chan string
	messages       chan string
}

func (b *Broker) Start() {
	go func() {
		for {
			select {

			case s := <-b.newClients:
				b.clients[s] = true
				// Added new client

			case s := <-b.defunctClients:

				delete(b.clients, s)
				close(s)
				// Removed client

			case msg := <-b.messages:

				for s := range b.clients {
					s <- msg
				}
			}
		}
	}()
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	messageChan := make(chan string)

	b.newClients <- messageChan

	notify := w.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		b.defunctClients <- messageChan
		log.Println("HTTP connection just closed.")
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	for {

		msg, open := <-messageChan

		if !open {
			break
		}
		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}

	log.Println("Finished HTTP request at ", r.URL.Path)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {

		log.Fatal(fmt.Sprintf("Error parsing template: %v", err))

	}

	t.Execute(w, nil)

	log.Println("Finished HTTP request at", r.URL.Path)
}

func main() {

	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	// Start processing events
	b.Start()

	http.Handle("/events/", b)

	go func() {
		freedom := time.Date(2022, time.February, 11, 16, 0, 0, 0, time.Now().Location())

		for i := 0; ; i++ {

			var diff time.Duration = freedom.Sub(time.Now())
			var days = diff / (60 * 60 * 24 * 1000000000)
			diff = diff - days*(60*60*24*1000000000)
			var hours = diff / time.Hour
			diff = diff - hours*time.Hour
			var minutes = diff / time.Minute
			diff = diff - minutes*time.Minute
			var seconds = diff / time.Second
			diff = diff - seconds*time.Second
			var tens = diff / 100000000
			b.messages <- fmt.Sprintf("%d days %02d:%02d:%02d,%d", days, hours, minutes, seconds, tens)
			time.Sleep(time.Millisecond * 100)

		}
	}()

	http.Handle("/", http.HandlerFunc(handler))

	http.ListenAndServe(":8181", nil)
}
