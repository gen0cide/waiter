package waiter

import (
	"io"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/cheggaaa/pb.v2"
)

var (
	defaultBrakeTime = time.Duration(100) * time.Millisecond
	defaultLength    = 100
	template         = ` >> {{ counters . | yellow }} {{ bar . ("[" | white) ("=" | green) (">" | green) ("--" | red) ("]" | white) }} {{ percent . | cyan }}`
)

// Waiter is a blend of a sync.WaitGroup and a terminal progress bar.
type Waiter struct {
	sync.RWMutex
	wg      *sync.WaitGroup
	bar     *pb.ProgressBar
	ipcount int64
}

// New returns a new Waiter with defaults
func New(writer io.Writer) *Waiter {
	b := pb.New64(1)
	b.SetTemplate(pb.ProgressBarTemplate(template))
	b.SetWidth(50)
	b.SetWriter(writer)
	wg := new(sync.WaitGroup)
	return &Waiter{
		bar: b,
		wg:  wg,
	}
}

// Add functions just like sync.WaitGroup's Add function
func (w *Waiter) Add(delta int) {
	atomic.AddInt64(&w.ipcount, int64(delta))
	w.Lock()
	w.bar.SetTotal(w.ipcount)
	w.Unlock()
	w.wg.Add(delta)
}

// Done functions just like sync.WaitGroup's Done function
func (w *Waiter) Done() {
	w.bar.Increment()
	w.wg.Done()
}

// Wait functionsdfasdf just like sync.WaitGroup's Wait function
// with an option to automatically start and stop the progress bar
func (w *Waiter) Wait(autorun bool) {
	w.Start()
	w.wg.Wait()
	w.bar.Increment()
	w.Stop()
}

// Start begins to render the progress bar in the terminal
func (w *Waiter) Start() {
	w.bar.Start()
}

// Stop ends the progress bar's rendering in the terminal
func (w *Waiter) Stop() {
	w.bar.Finish()
}
