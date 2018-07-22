package waiter

import (
	"sync"

	"github.com/gosuri/uiprogress"
)

// Waiter is a blend of a sync.WaitGroup and a terminal progress bar.
type Waiter struct {
	sync.RWMutex
	progress *uiprogress.Progress
	wg       *sync.WaitGroup
	bar      *uiprogress.Bar
}

// New returns a new Waiter with defaults
func New() *Waiter {
	p := uiprogress.New()
	b := p.AddBar(0)
	b.AppendCompleted()
	b.PrependElapsed()
	wg := new(sync.WaitGroup)
	return &Waiter{
		progress: p,
		bar:      b,
		wg:       wg,
	}
}

// Add functions just like sync.WaitGroup's Add function
func (w *Waiter) Add(delta int) {
	w.Lock()
	w.bar.Total += delta
	w.Unlock()
	w.wg.Add(delta)
}

// Done functions just like sync.WaitGroup's Done function
func (w *Waiter) Done() {
	w.bar.Incr()
	w.wg.Done()
}

// Wait functions just like sync.WaitGroup's Wait function
// with an option to automatically start and stop the progress bar
func (w *Waiter) Wait(autorun bool) {
	if autorun {
		w.Start()
	}
	w.wg.Wait()
	if autorun {
		w.Stop()
	}
}

// Start begins to render the progress bar in the terminal
func (w *Waiter) Start() {
	w.progress.Start()
}

// Stop ends the progress bar's rendering in the terminal
func (w *Waiter) Stop() {
	w.progress.Stop()
}
