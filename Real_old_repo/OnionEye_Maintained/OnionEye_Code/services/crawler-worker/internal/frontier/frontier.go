// Package frontier implements a per-job, depth-bounded URL frontier with a
// seen-set for deduplication. Per the architecture, the seen-set is allocated
// at job start and discarded at job end — it is never shared across jobs.
package frontier

// Entry is one URL to fetch and its crawl depth.
type Entry struct {
	URL   string
	Depth int
}

// Frontier is a BFS queue bounded by maxDepth, deduped by URL.
type Frontier struct {
	maxDepth int
	seen     map[string]bool
	queue    []Entry
}

// New creates an empty frontier.
func New(maxDepth int) *Frontier {
	if maxDepth < 1 {
		maxDepth = 1
	}
	return &Frontier{maxDepth: maxDepth, seen: make(map[string]bool)}
}

// Add enqueues url at depth if it is within the depth bound and not seen.
// Returns false if the URL was ignored (too deep or already visited).
func (f *Frontier) Add(url string, depth int) bool {
	if depth > f.maxDepth {
		return false
	}
	if f.seen[url] {
		return false
	}
	f.seen[url] = true
	f.queue = append(f.queue, Entry{URL: url, Depth: depth})
	return true
}

// Next pops the next URL in BFS order. Returns false when the frontier is empty.
func (f *Frontier) Next() (Entry, bool) {
	if len(f.queue) == 0 {
		return Entry{}, false
	}
	e := f.queue[0]
	f.queue = f.queue[1:]
	return e, true
}

// Len reports the number of URLs still queued.
func (f *Frontier) Len() int { return len(f.queue) }
