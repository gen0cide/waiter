# Waiter

Simple library that marries `sync.WaitGroup` with a `uiprogress.Bar` for easily configurable progress of concurrent work.

## Installation

```sh
go get github.com/gen0cide/waiter
```

## Usage

```golang
// Create your Waiter
wg := waiter.New()

count := 10 // example num of workers

// populate with work
for i := 0; i < count; i++ {
  // increase your waiter
  wg.Add(1)

  // start the workers
  go func(id int) {
    time.Sleep(time.Duration(id) * time.Second)
    wg.Done()
  }(i)
}

// block on the work finishing
// and render the UI as we wait for it to complete
wg.Wait(true)

// That's it!
```
