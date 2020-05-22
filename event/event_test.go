// Package event_test contains unit tests for the lapis event system.
package event_test

import (
	"errors"
	"testing"

	"github.com/bsladewski/lapis/event"
)

// TestEventPositive tests the execution of an event worker that completes
// successfully.
func TestEventPositive(t *testing.T) {

	// create a new event worker
	w, err := event.NewWorker("test")

	if err != nil {
		t.Fatal(err)
	}

	// assert that the worker has been assigned an id greater than zero
	id := w.GetID()

	if id <= 0 {
		t.Fatalf("expected id greater than zero, got %d", id)
	}

	t.Logf("created event worker id %d", id)

	var notes = "lorem ipsum dolor sit amet"

	// add a work function that does nothing and returns no error
	w.AddWork(func() error {

		return nil

	})

	// add a work function that adds task notes
	w.AddWork(func() error {

		w.AddNotes(notes)
		return nil

	})

	// assert that the worker notes field is empty
	if workerNotes := w.GetNotes(); workerNotes != "" {
		t.Fatalf("expected empty notes, got: %s", workerNotes)
	}

	// execute the worker and assert that no error was returned
	if err := w.Do(); err != nil {
		t.Fatal(err)
	}

	// assert that the worker is in the finished state
	if state := w.GetState(); state != event.Finished {
		t.Fatalf("expected finished state (%d), got %d", event.Finished, state)
	}

	// assert that the notes field was filled in successfully
	if workerNotes := w.GetNotes(); workerNotes != notes {
		t.Fatalf("expected notes '%s', got '%s'", notes, workerNotes)
	}

}

// TestEventNegative tests the execution of an event worker that fails with an
// error.
func TestEventNegative(t *testing.T) {

	// create a new event worker
	w, err := event.NewWorker("test_error")

	if err != nil {
		t.Fatal(err)
	}

	// assert that the worker has been assigned an id greater than zero
	id := w.GetID()

	if id <= 0 {
		t.Fatalf("expected id greater than zero, got %d", id)
	}

	t.Logf("created event worker id %d", id)

	var workerError = "'Twas brillig, and the slithy toves did gyre and gimble in the wabe"

	var notes = "The vorpal blade went snicker-snack!"

	// add a work function that does nothing and returns no error
	w.AddWork(func() error {

		return nil

	})

	// add a work function that returns an error
	w.AddWork(func() error {

		return errors.New(workerError)

	})

	// assert that the worker errors field is empty
	if workerErrors := w.GetErrors(); workerErrors != nil {
		t.Fatalf("expected no errors, got: %v", workerErrors)
	}

	// add a work function that appends notes but should not be executed due to
	// to the previous error
	w.AddWork(func() error {

		w.AddNotes(notes)

		return nil

	})

	// assert that the worker notes field is empty
	if workerNotes := w.GetNotes(); workerNotes != "" {
		t.Fatalf("expected no notes, got: %s", workerNotes)
	}

	// execute the worker and assert that an error was returned
	if err := w.Do(); err == nil {
		t.Fatal("expected an error, got nil")
	}

	// assert that the notes field is still empty
	if workerNotes := w.GetNotes(); workerNotes != "" {
		t.Fatalf("expected no notes, got: %s", workerNotes)
	}

	// assert that the worker is in the failed state
	if state := w.GetState(); state != event.Failed {
		t.Fatalf("expected failed state (%d), got %d", event.Failed, state)
	}

	// assert that the errors field was filled in successfully
	if workerErrors := w.GetErrors(); workerErrors.Error() != workerError {
		t.Fatalf("expected notes '%s', got '%v'", workerError, workerErrors)
	}

}
