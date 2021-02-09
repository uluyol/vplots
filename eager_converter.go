package main

import (
	"runtime"
	"sync"

	"github.com/uluyol/vplots/conversions"
)

type resultErr struct {
	data []byte
	err  error
}

type workItem struct {
	i      int
	width  int
	notify chan<- resultErr
}

type cacheEntry struct {
	id    int
	width int
}

type eagerPNGConverter struct {
	pdfPaths          []string
	numAllowedWorkers int

	mu               sync.Mutex
	pngCache         map[cacheEntry][]byte
	c                *conversions.Converter
	hipri            []workItem
	lopri            []workItem
	numActiveWorkers int
}

func newEagerPNGConverter(pdfPaths []string, c *conversions.Converter) *eagerPNGConverter {
	return &eagerPNGConverter{
		pdfPaths:          pdfPaths,
		numAllowedWorkers: runtime.GOMAXPROCS(0) / 2,
		pngCache:          make(map[cacheEntry][]byte),
		c:                 c,
	}
}

func (c *eagerPNGConverter) kickWorkersLocked() {
	numCanStart := c.numAllowedWorkers - c.numActiveWorkers
	maxNeeded := len(c.hipri) + len(c.lopri)
	toStart := numCanStart
	if maxNeeded < numCanStart {
		toStart = maxNeeded
	}
	for i := 0; i < toStart; i++ {
		go c.runWorker()
	}
	c.numActiveWorkers += toStart
}

func (c *eagerPNGConverter) runWorker() {
	c.mu.Lock()
	for {
		var wi workItem
		logAttr := ""
		switch {
		case len(c.hipri) > 0:
			wi = c.hipri[len(c.hipri)-1]
			c.hipri = c.hipri[:len(c.hipri)-1]
			logAttr = "hipri"
		case len(c.lopri) > 0:
			wi = c.lopri[len(c.lopri)-1]
			c.lopri = c.lopri[:len(c.lopri)-1]
			logAttr = "lopri"
		default:
			// no more work
			c.numActiveWorkers--
			c.mu.Unlock()
			return
		}

		{
			result, ok := c.pngCache[cacheEntry{wi.i, wi.width}]
			c.mu.Unlock()
			if ok {
				if wi.notify != nil {
					wi.notify <- resultErr{result, nil}
				}
				c.mu.Lock()
				continue
			}
		}

		data, err := c.c.ToPNG(c.pdfPaths[wi.i], wi.width, logAttr)
		if err == nil {
			c.mu.Lock()
			c.pngCache[cacheEntry{wi.i, wi.width}] = data
			c.mu.Unlock()
		}

		if wi.notify != nil {
			wi.notify <- resultErr{data, err}
		}
		c.mu.Lock()
	}
}

func queuedIn(i, width int, q []workItem) bool {
	for _, wi := range q {
		if wi.i == i && wi.width == width {
			return true
		}
	}
	return false
}

func (c *eagerPNGConverter) Get(i int, width int) ([]byte, error) {
	c.mu.Lock()

	result, ok := c.pngCache[cacheEntry{i, width}]
	if ok {
		c.mu.Unlock()
		return result, nil
	}

	rc := make(chan resultErr)
	c.hipri = append(c.hipri, workItem{i, width, rc})

	// Enqueue conversion of other PDFs at the same resolution at lower priority
	for j := range c.pdfPaths {
		// Don't enqueue i at lopri
		if j == i {
			continue
		}
		// Don't enqueue if done
		if _, ok := c.pngCache[cacheEntry{j, width}]; ok {
			continue
		}
		// Don't enqueue if already queued
		if queuedIn(j, width, c.hipri) {
			continue
		}
		if queuedIn(j, width, c.lopri) {
			continue
		}
		c.lopri = append(c.lopri, workItem{j, width, nil})
	}
	c.kickWorkersLocked()
	c.mu.Unlock()

	ret := <-rc
	return ret.data, ret.err
}
