// Package event provides data types and tools for working with events in the
// lapis system.
package event

import (
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// A Worker is used to execute a lapis event.
type Worker interface {
	AddWork(work WorkFunction)
	AddCriticalWork(work WorkFunction)
	AddNotes(notes string)
	AddError(err error)
	Do() error
}

// WorkFunction encapsulates some work to be executed by a lapis task.
type WorkFunction func() error

// NewWorker returns an empty event worker.
func NewWorker(name string) Worker {

	w := &worker{Name: name}

	return w

}

// A worker is the default concrete implementation of a lapis Worker.
type worker struct {
	ID     uint
	Name   string
	Errors string
	Notes  string
	State  WorkerState
	work   []workItem
}

// WorkerState stores the current state of a lapis event.
type WorkerState int

const (
	// Pending indicates that the worker has not started executing.
	Pending WorkerState = iota
	// Running indicates that the worker is currently running.
	Running
	// Finished indicates that the worker has run successfully.
	Finished
	// Failed indicates that the worker encountered an error during execution.
	Failed
)

// workItem represents a function to be executed as part of processing an event.
type workItem struct {
	critical bool
	work     WorkFunction
}

// AddWork adds a work function to this lapis event; work functions are executed
// in the order they are added; if an error is encountered while executing a
// work function subsequent work functions will not be executed unless marked as
// critical.
func (w *worker) AddWork(work WorkFunction) {
	w.work = append(w.work, workItem{false, work})
}

// AddCriticalWork adds a critical work function to this lapis event; work
// functions are executed in the order they are added; if an error is
// encountered while executing a work function subsequent work functions will
// not be executed unless marked as critical.
func (w *worker) AddCriticalWork(work WorkFunction) {
	w.work = append(w.work, workItem{true, work})
}

// AddNotes appends the supplied string to the event worker notes field; notes
// are indended to provide insight into the execution of the worker but do not
// represent an error.
func (w *worker) AddNotes(notes string) {
	if w.Notes == "" {
		w.Notes = notes
		return
	}
	w.Notes = fmt.Sprintf("%s\n%s", w.Notes, notes)
}

// AddError appends the supplied error to the event worker errors field; adding
// an error puts the event into the failed state.
func (w *worker) AddError(err error) {
	w.State = Failed
	if w.Errors == "" {
		w.Errors = err.Error()
		return
	}
	w.Errors = fmt.Sprintf("%s\n%s", w.Errors, err.Error())
}

// Do executes an event worker returning any errors encountered.
func (w *worker) Do() error {

	// execute each work function sequentially
	for _, work := range w.work {

		// if we are in an error state only execute critical work functions
		if work.critical || w.Errors == "" {

			// execute the work function and add any errors to the worker
			if err := work.work(); err != nil {
				w.AddError(err)
			}

		}

	}

	// if errors were encountered return them
	if w.State == Failed {
		return errors.New(w.Errors)
	}

	return nil
}
