// Package event provides data types and tools for working with events in the
// lapis system.
package event

import (
	"errors"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3" // support for SQLite database
)

// db provides access to persistant storage.
var db *gorm.DB

// init initializes the SQLite database.
func init() {

	var err error

	// create tmp directory for persistant storage
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		if err := os.Mkdir("tmp", 0700); err != nil {
			panic(err)
		}
	}

	db, err = gorm.Open("sqlite3", "tmp/lapis.db")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(worker{})

}

// A Worker is used to execute a lapis event.
type Worker interface {
	GetID() uint
	AddWork(work WorkFunction)
	AddCriticalWork(work WorkFunction)
	AddTeardown(work WorkFunction)
	AddCriticalTeardown(work WorkFunction)
	AddNotes(notes string)
	AddError(err error)
	GetNotes() string
	GetErrors() error
	GetState() WorkerState
	Do() error
}

// WorkFunction encapsulates some work to be executed by a lapis task.
type WorkFunction func() error

// NewWorker returns an empty event worker.
func NewWorker(name string) (Worker, error) {

	w := &worker{Name: name}

	// create the initial event worker record
	if err := save(w); err != nil {
		return nil, err
	}

	// update the worker status when execution begins
	w.AddCriticalWork(func() error {

		w.State = Running

		return save(w)

	})

	// update the worker status upon completion of the event worker
	w.AddCriticalTeardown(func() error {

		if w.State != Failed {
			w.State = Finished
		}

		return save(w)

	})

	return w, nil

}

// A worker is the default concrete implementation of a lapis Worker.
type worker struct {
	gorm.Model
	Name     string      `gorm:"index:name" json:"name"`
	Errors   string      `json:"errors"`
	Notes    string      `json:"notes"`
	State    WorkerState `gorm:"index:state" json:"state"`
	work     []workItem  `gorm:"-"`
	teardown []workItem  `gorm:"-"`
}

// save inserts or updates an event worker.
func save(w *worker) error {

	return db.Save(w).Error

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

// GetID retrieves the id of this event worker.
func (w *worker) GetID() uint {

	return w.ID

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

// AddTeardown adds a teardown function to this lapis event; work functions are
// executed after the work functions and in the order they are added; if an
// error is encountered while executing a work function subsequent work
// functions will not be executed unless marked as critical.
func (w *worker) AddTeardown(teardown WorkFunction) {

	w.teardown = append(w.teardown, workItem{false, teardown})

}

// AddCriticalTeardown adds a critical teardown function to this lapis event;
// work functions are executed in the order they are added; if an error is
// encountered while executing a work function subsequent work functions will
// not be executed unless marked as critical.
func (w *worker) AddCriticalTeardown(teardown WorkFunction) {

	w.teardown = append(w.teardown, workItem{true, teardown})

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

// GetNotes returns the contents of the event worker notes column.
func (w *worker) GetNotes() string {

	return w.Notes

}

// GetErrors returns the contents of the event worker errors column.
func (w *worker) GetErrors() error {

	if w.Errors == "" {
		return nil
	}

	return errors.New(w.Errors)

}

// GetState returns the current state of the event worker.
func (w *worker) GetState() WorkerState {

	return w.State

}

// Do executes an event worker returning any errors encountered.
func (w *worker) Do() error {

	// execute each work function sequentially
	for _, work := range append(w.work, w.teardown...) {

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
