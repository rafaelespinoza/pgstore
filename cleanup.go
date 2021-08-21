package pgstore

import (
	"context"
	"fmt"
	"log"
	"time"
)

var defaultInterval = time.Minute * 5

// Cleanup runs a background goroutine every interval that deletes expired
// sessions from the database.
//
// The design is based on https://github.com/yosssi/boltstore
func (db *PGStore) Cleanup(interval time.Duration) (chan<- struct{}, <-chan struct{}) {
	quit, done := make(chan struct{}), make(chan struct{})

	go func() {
		if interval <= 0 {
			interval = defaultInterval
		}
		ticker := time.NewTicker(interval)

		defer ticker.Stop()

		for {
			select {
			case <-quit:
				// Handle the quit signal.
				done <- struct{}{}
				return
			case <-ticker.C:
				// Delete expired sessions on each tick.
				err := db.DeleteExpired()
				if err != nil {
					log.Printf("pgstore: unable to delete expired sessions: %v", err)
				}
			}
		}
	}()

	return quit, done
}

// StopCleanup stops the background cleanup from running.
func (db *PGStore) StopCleanup(quit chan<- struct{}, done <-chan struct{}) {
	quit <- struct{}{}
	<-done
}

// DeleteExpired deletes expired sessions from the database.
func (db *PGStore) DeleteExpired() (err error) {
	_, err = db.DbPool.Exec("DELETE FROM http_sessions WHERE expires_on < now()")
	return
}

// RunCleanup deletes expired sessions from the database every interval within a
// background goroutine. It's similar to the Cleanup method except its API
// provides a way to consume errors via the read-only channel. Use the input
// context to facilitate quit or timeout signals from the parent goroutine.
func (db *PGStore) RunCleanup(ctx context.Context, interval time.Duration) <-chan error {
	errs := make(chan error)

	go func() {
		if interval <= 0 {
			interval = defaultInterval
		}
		ticker := time.NewTicker(interval)

		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				close(errs)
				return
			case <-ticker.C:
				// Delete expired sessions on each tick.
				err := db.DeleteExpired()
				if err != nil {
					err = fmt.Errorf("pgstore: unable to delete expired sessions: %v", err)
				}
				errs <- err
			}
		}
	}()

	return errs
}
