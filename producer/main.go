package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpAddr   = flag.String("http", ":8080", "Address to listen for requests on")
	redisAddr  = flag.String("redis-server", "redis:6379", "Redis server address")
	redisQueue = flag.String("redis-queue", "esilv", "Redis queue name")
)

func main() {
	flag.Parse()

	client := redis.NewClient(&redis.Options{
		Addr: *redisAddr,
	})

	p := newProducer(client, *redisQueue)

	log.Printf("Listening on %s...", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, p.router); err != nil {
		log.Fatalf("Unable to listen for requests: %v", err)
	}
}

type producer struct {
	router    *mux.Router
	client    *redis.Client
	queue     string
	published prometheus.Counter
}

func newProducer(client *redis.Client, queue string) *producer {
	p := &producer{
		router: mux.NewRouter(),
		client: client,
		queue:  queue,
		published: promauto.NewCounter(prometheus.CounterOpts{
			Name: "producer_messages_published_total",
			Help: "The total number of published messages",
		}),
	}

	p.router.Handle("/publish", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handlePublish))).Methods("POST")
	p.router.Handle("/publish/{count:[0-9]+}", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(p.handlePublishMany))).Methods("POST")
	p.router.HandleFunc("/healthz", p.handleHealthcheck).Methods("GET")
	p.router.Handle("/metrics", promhttp.Handler())

	return p
}

func (p *producer) handlePublish(w http.ResponseWriter, r *http.Request) {
	if err := p.publish(r.Context()); err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish message: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "Message published successfully\n")
}

func (p *producer) handlePublishMany(w http.ResponseWriter, r *http.Request) {
	countStr := mux.Vars(r)["count"]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, "Invalid count parameter", http.StatusBadRequest)
		return
	}
	if err := p.publishMany(r.Context(), count); err != nil {
		http.Error(w, fmt.Sprintf("Failed to publish messages: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d messages published successfully\n", count)
}

func (p *producer) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "The server is healthy\n")
}

func (p *producer) publish(ctx context.Context) error {
    log.Println("Attempting to publish message to Redis...")
    if err := p.client.Ping(ctx).Err(); err != nil {
        log.Printf("Error connecting to Redis: %v", err)
        return fmt.Errorf("unable to connect to Redis: %w", err)
    }
    if err := p.client.RPush(ctx, p.queue, "Hello, World!").Err(); err != nil {
        log.Printf("Error publishing message to Redis: %v", err)
        return fmt.Errorf("failed to publish message: %w", err)
    }
    p.published.Inc()
    log.Println("Message published successfully")
    return nil
}

func (p *producer) publishMany(ctx context.Context, count int) error {
    log.Printf("Attempting to publish %d messages to Redis...", count)
    for i := 0; i < count; i++ {
        if err := p.publish(ctx); err != nil {
            log.Printf("Error publishing message %d: %v", i+1, err)
            return fmt.Errorf("only %d of %d messages published: %w", i, count, err)
        }
    }
    log.Printf("%d messages published successfully", count)
    return nil
}
