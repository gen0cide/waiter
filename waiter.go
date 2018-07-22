package waiter

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/gosuri/uiprogress"
)

var (
	defaultRefreshInterval = time.Duration(500) * time.Millisecond
	defaultLength          = 100
)

// Waiter is a blend of a sync.WaitGroup and a terminal progress bar.
type Waiter struct {
	sync.RWMutex
	progress *uiprogress.Progress
	wg       *sync.WaitGroup
	bar      *uiprogress.Bar
	tunits   uint64
	cunits   uint64
	currprog int
}

// New returns a new Waiter with defaults
func New() *Waiter {
	p := uiprogress.New()
	b := p.AddBar(100)
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
	atomic.AddUint64(&w.tunits, uint64(delta))
	w.Lock()
	w.bar.Total += delta
	w.Unlock()
	w.wg.Add(delta)
}

// Done functions just like sync.WaitGroup's Done function
func (w *Waiter) Done() {
	w.bar.Set(1)
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
