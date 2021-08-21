package pgstore

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestCleanup(t *testing.T) {
	dsn := os.Getenv("PGSTORE_TEST_CONN")
	if dsn == "" {
		t.Skip("This test requires a real database.")
	}

	ss, err := NewPGStore(dsn, []byte(secret))
	if err != nil {
		t.Fatal("failed to get store", err)
	}
	defer ss.Close()

	// Start the cleanup goroutine.
	defer ss.StopCleanup(ss.Cleanup(time.Millisecond * 500))

	if err := makeExpiredSession(ss, 1); err != nil {
		t.Fatal(err)
	}
	defer func() { deleteAllSessions(t, ss) }()

	// Give the ticker a moment to run.
	time.Sleep(time.Millisecond * 1500)

	// SELECT expired sessions. We should get a count of zero back.
	var count int
	err = ss.DbPool.QueryRow("SELECT count(*) FROM http_sessions WHERE expires_on < now()").Scan(&count)
	if err != nil {
		t.Fatalf("failed to select http_sessions; %v", err)
	}

	if count != 0 {
		t.Fatalf("wrong number of remaining sessions; got %d, expected %d", count, 0)
	}
}

func TestRunCleanup(t *testing.T) {
	dsn := os.Getenv("PGSTORE_TEST_CONN")
	if dsn == "" {
		t.Skip("This test requires a real database.")
	}

	ss, err := NewPGStore(dsn, []byte(secret))
	if err != nil {
		t.Fatal("failed to get store", err)
	}
	defer ss.Close()

	if err := makeExpiredSession(ss, 1); err != nil {
		t.Fatal(err)
	}
	defer func() { deleteAllSessions(t, ss) }()

	time.Sleep(time.Millisecond * 1500) // Allow sessions to expire.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var numOKs, numTimeouts int
	for err := range ss.RunCleanup(ctx, time.Millisecond*500) {
		switch err {
		case nil:
			numOKs++
			t.Log("cleanup ok")
		case context.DeadlineExceeded:
			numTimeouts++
			t.Log("timeout")
		default:
			t.Error(err)
		}
	}
	if numOKs < 1 {
		t.Errorf("expected at least 1 ok, but got %d", numOKs)
	}
	if numTimeouts != 1 {
		t.Errorf("%d timeout is OK, but got %d", 1, numTimeouts)
	}

	var count int
	err = ss.DbPool.QueryRow("SELECT count(*) FROM http_sessions WHERE expires_on < now()").Scan(&count)
	if err != nil {
		t.Fatalf("failed to select http_sessions; %v", err)
	}

	if count != 0 {
		t.Fatalf("wrong number of remaining sessions; got %d, expected %d", count, 0)
	}
}

func TestDeleteExpired(t *testing.T) {
	dsn := os.Getenv("PGSTORE_TEST_CONN")
	if dsn == "" {
		t.Skip("This test requires a real database.")
	}

	ss, err := NewPGStore(dsn, []byte(secret))
	if err != nil {
		t.Fatal("failed to get store", err)
	}
	defer ss.Close()

	for i := 0; i < 4; i++ {
		maxAge := 1
		if i%2 == 1 {
			maxAge = 9999
		}
		if err := makeExpiredSession(ss, maxAge); err != nil {
			t.Fatal(err)
		}
	}
	defer func() { deleteAllSessions(t, ss) }()
	time.Sleep(time.Millisecond * 1500) // Allow some sessions to expire.
	if err = ss.DeleteExpired(); err != nil {
		t.Fatal(err)
	}

	var countDeleted, countRemaining int
	err = ss.DbPool.QueryRow("SELECT count(*) FROM http_sessions WHERE expires_on < now()").Scan(&countDeleted)
	if err != nil {
		t.Fatalf("failed to select http_sessions; %v", err)
	}
	if countDeleted != 0 {
		t.Fatalf("wrong number of remaining sessions; got %d, expected %d", countDeleted, 0)
	}

	err = ss.DbPool.QueryRow("SELECT count(*) FROM http_sessions WHERE expires_on >= now()").Scan(&countRemaining)
	if err != nil {
		t.Fatalf("failed to select http_sessions; %v", err)
	}
	if countRemaining != 2 {
		t.Fatalf("wrong number of remaining sessions; got %d, expected %d", countRemaining, 2)
	}
}

func makeExpiredSession(store *PGStore, maxAge int) error {
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		return fmt.Errorf("failed to create request; %v", err)
	}

	session, err := store.Get(req, "newsess")
	if err != nil {
		return fmt.Errorf("failed to create session; %v", err)
	}

	// Setup session to expire.
	session.Options.MaxAge = maxAge

	m := make(http.Header)
	if err = store.Save(req, headerOnlyResponseWriter(m), session); err != nil {
		return fmt.Errorf("failed to save session; %v", err)
	}

	return nil
}

func deleteAllSessions(t *testing.T, store *PGStore) {
	t.Helper()
	if _, err := store.DbPool.Exec("DELETE FROM http_sessions"); err != nil {
		t.Fatal(err)
	}
	t.Log("deleted leftover http_sessions")
}
