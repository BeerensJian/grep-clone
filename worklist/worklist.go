package worklist

type Entry struct {
	Path string
}

type Worklist struct {
	jobs chan Entry
}

func (w *Worklist) Add(work Entry) {
	w.jobs <- work
}

func (w *Worklist) Next() Entry {
	e := <-w.jobs
	return e
}
