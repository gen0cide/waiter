package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/x-cray/logrus-prefixed-formatter"

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
	l      *logrus.Entry
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
			n.logRandQuote()
		}
	}
}

func (n *noisemaker) logRandQuote() {
	roll := rand.Intn(10)
	if roll > 7 {
		n.l.Errorf("%s", getquote.GetQuote())
		return
	}
	n.l.Infof("%s", getquote.GetQuote())
}

func main() {
	logger := logrus.New().WithField("prefix", "test")
	wtr := waiter.New("test", logger.Logger.Out)
	formatter := new(prefixed.TextFormatter)
	formatter.ForceColors = true
	formatter.ForceFormatting = true
	logger.Logger.Formatter = formatter
	logger.Logger.SetLevel(logrus.DebugLevel)
	nm := &noisemaker{
		l: logger,
	}
	logger.Logger.Out = wtr
	go nm.makenoise()
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
		wtr.Reset(fmt.Sprintf("test_%d", i))
		go addTo(wtr, t.count, t.delay)
		nm.start()
		wtr.Start()
		time.Sleep(10 * time.Second)
		nm.stop()
		wtr.Wait()
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
