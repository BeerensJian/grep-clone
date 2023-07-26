package worklist

import "context"

type Entry struct {
	Path string
}

type Worklist struct {
	jobs chan Entry
	Ctx  context.Context
}

func (w *Worklist) Add(work Entry) {
	w.jobs <- work
}

func (w *Worklist) Next() Entry {
	e := <-w.jobs
	return e
}

func New(bufSize int) Worklist {
	return Worklist{make(chan Entry, bufSize), context.Background()}
}

func NewJob(path string) Entry {
	return Entry{path}
}

func (w *Worklist) Finalize(numWorker int) {
	for i := 0; i < numWorker; i++ {
		w.Add(Entry{""})
	}
}
