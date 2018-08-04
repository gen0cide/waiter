package waiter

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/fatih/color"

	"gopkg.in/cheggaaa/pb.v2"
)

var (
	titleColor    = color.New(color.FgHiWhite, color.Bold)
	labelColor    = color.New(color.FgHiGreen)
	labelaltColor = color.New(color.FgHiWhite)
	template      = `{{ string . "title" }}{{ counters . "%s/%s" "%s/?" | yellow }} {{ bar . (white "[") (green "=") (green ">") (red "--") (white "]") }} {{ percent . | yellow }} {{ etime . | cyan }}`
)

// Waiter is a blend of a sync.WaitGroup and a terminal progress bar.
type Waiter struct {
	sync.RWMutex
	name    string
	wg      *sync.WaitGroup
	bar     *pb.ProgressBar
	wr      io.Writer
	ipcount int64
	started bool
}

// New returns a new Waiter with defaults
func New(name string, writer io.Writer) *Waiter {
	if name == "" {
		name = "default"
	}
	b := pb.New64(1)
	b.SetTemplate(pb.ProgressBarTemplate(template))
	b.SetWriter(writer)
	b.Set("prefix", name)
	title := fmt.Sprintf("%s%s%s %s ", titleColor.Sprintf("  STATUS"), labelaltColor.Sprintf(":"), labelColor.Sprintf(name), labelaltColor.Sprintf(">>"))
	b.Set("title", title)
	b.SetWidth(150)
	wg := new(sync.WaitGroup)
	return &Waiter{
		bar:  b,
		wg:   wg,
		wr:   writer,
		name: name,
	}
}

// Reset resets the waiter back to a default state with a new name
func (w *Waiter) Reset(name string) {
	if name == "" {
		name = "default"
	}
	b := pb.New64(1)
	b.SetTemplate(pb.ProgressBarTemplate(template))
	b.SetWriter(w.wr)
	b.Set("prefix", name)
	title := fmt.Sprintf("%s%s%s %s ", titleColor.Sprintf("  STATUS"), labelaltColor.Sprintf(":"), labelColor.Sprintf(name), labelaltColor.Sprintf(">>"))
	b.Set("title", title)
	b.SetWidth(150)
	wg := new(sync.WaitGroup)
	w.wg = wg
	w.bar = b
	w.name = name
	w.ipcount = 0
	w.started = false
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
func (w *Waiter) Wait() {
	w.Start()
	w.wg.Wait()
	w.Stop()
}

// Start begins to render the progress bar in the terminal
func (w *Waiter) Start() {
	w.started = true
	w.bar.Start()
}

// Stop ends the progress bar's rendering in the terminal
func (w *Waiter) Stop() {
	w.bar.Finish()
	w.started = false
}

func (w *Waiter) Write(p []byte) (n int, err error) {
	w.Lock()
	if w.started {
		w.bar.Write()
		w.wr.Write([]byte("\n"))
	}
	a, b := w.wr.Write(p)
	if w.started {
		w.bar.Write()
	}
	w.Unlock()
	return a, b
}
