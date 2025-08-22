package concurrency

// simple semaphore for concurrency limiting

var sem chan struct{}

func Init(maxConcurrent int) {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}
	sem = make(chan struct{}, maxConcurrent)
}

func Acquire() bool {
	// non-blocking acquire; returns true if slot obtained
	select {
	case sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func Release() {
	select {
	case <-sem:
	default:
		// nothing to release
	}
}
