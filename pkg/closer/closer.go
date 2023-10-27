package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

// Add func to closer
func Add(f ...func() error) {
	globalCloser.Add(f...)
}

// Wait when all functions will be done
func Wait() {
	globalCloser.Wait()
}

// CloseAll connections
func CloseAll() {
	globalCloser.CloseAll()
}

// Closer struct for closer
type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	done  chan struct{}
	funcs []func() error
}

// New instance of closer
func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}

	if len(sig) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sig...)

			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}

	return c
}

// Add func to slice
func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

// Wait when closer wil be done
func (c *Closer) Wait() {
	<-c.done
}

// CloseAll do all functions before closing
func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))

		for _, f := range funcs {
			go func(f func() error) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}
