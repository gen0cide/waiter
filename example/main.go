package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mndrix/rand"
	"github.com/schollz/quotation-explorer/getquote"

	"github.com/pkg/errors"

	"github.com/gen0cide/waiter"
	"github.com/sirupsen/logrus"
)

var (
	noiseticker = time.NewTicker(1 * time.Second)
)

type testcase struct {
	count int
	delay time.Duration
}

func parseTest(s string) testcase {
	vals := strings.Split(s, ":")
	if len(vals) != 2 {
		panic(fmt.Errorf("argument %s is not in format count:duration", s))
	}

	cval, err := strconv.Atoi(vals[0])
	if err != nil {
		panic(errors.Wrap(errors.Wrap(err, "count argument is not a valid number"), "argument is not in format count:duration"))
	}

	dval, err := time.ParseDuration(vals[1])
	if err != nil {
		panic(errors.Wrap(errors.Wrap(err, "duration argument is not a valid duration string"), "argument is not in format count:duration"))
	}

	return testcase{
		count: cval,
		delay: dval,
	}
}

type noisemaker struct {
	l      *logrus.Logger
	active bool
}

func (n *noisemaker) start() {
	n.active = true
}

func (n *noisemaker) stop() {
	n.active = false
}

func (n *noisemaker) makenoise() {
	for range noiseticker.C {
		if n.active {
			logRandQuote(n.l)
		}
	}
}

func logRandQuote(l *logrus.Logger) {
	roll := rand.Intn(10)
	if roll > 7 {
		l.Errorf("%s", getquote.GetQuote())
		return
	}
	l.Infof("%s", getquote.GetQuote())
}

func main() {
	logger := logrus.New()
	tcs := []testcase{}

	for i, a := range os.Args {
		if i == 0 {
			continue
		}
		tcs = append(tcs, parseTest(a))
	}

	for i, t := range tcs {
		logger.Infof("##### BEGINNING TEST %d #####", i)
		logger.Infof("\tPARAMS: count=%d duration=%s (seconds)", t.count, t.delay.Seconds())
		nl := logrus.New()
		wtr := waiter.New(fmt.Sprintf("test-%d", i), os.Stderr)
		nl.Out = wtr
		nm := &noisemaker{
			l: nl,
		}
		go nm.makenoise()
		go addTo(wtr, t.count, t.delay)
		nm.start()
		time.Sleep(10 * time.Second)
		nm.stop()
		wtr.Wait(true)
		logger.Infof("Done!")
	}

}

func addTo(wtr *waiter.Waiter, count int, slp time.Duration) {
	for i := 0; i < count; i++ {
		wtr.Add(1)
		delay := slp * 2
		go waitDone(wtr, delay)
		time.Sleep(slp)
	}
}

func waitDone(wtr *waiter.Waiter, slp time.Duration) {
	time.Sleep(slp)
	wtr.Done()
}
