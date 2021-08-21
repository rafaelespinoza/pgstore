package pgstore_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rafaelespinoza/pgstore"
)

// This example shows how to use pgstore in an HTTP handler.
func ExamplePGStore() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// Fetch new store.
		store, err := pgstore.NewPGStore(os.Getenv("DB_DSN"), []byte(os.Getenv("SECRET_KEY")))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to initialize session store:", err)
			return
		}
		defer store.Close()

		// Get a session.
		session, err := store.Get(r, "session-key")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to get session:", err)
			return
		}

		// Add a value.
		session.Values["foo"] = "bar"

		// Save.
		if err = session.Save(r, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to save session:", err)
			return
		}

		// Delete session.
		session.Options.MaxAge = -1
		if err = session.Save(r, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to delete session:", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}

	err := http.ListenAndServe("localhost:6543", http.HandlerFunc(handler))
	if err == http.ErrServerClosed {
		log.Println(err)
	} else {
		log.Fatal(err)
	}
}

func ExamplePGStore_RunCleanup() {
	store, err := pgstore.NewPGStore(os.Getenv("DB_DSN"), []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		panic(err)
	}
	defer store.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Run indefinitely. Shut it down with ^C (Control-C).
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		cancel()
	}()

	const interval = 5 * time.Minute

	// Consume error values as they're available.
	for err := range store.RunCleanup(ctx, interval) {
		switch err {
		case nil:
			log.Println("ran cleanup")
		case context.Canceled:
			log.Println("done", err)
		default:
			log.Println("error:", err)
		}
	}
}
