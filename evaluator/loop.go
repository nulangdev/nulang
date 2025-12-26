package evaluator

import "sync"

var (
	// eventLoopWg counts active async tasks (servers, timers, etc)
	eventLoopWg sync.WaitGroup
)

// RegisterAsyncTask tracks a new background task
func RegisterAsyncTask() {
	eventLoopWg.Add(1)
}

// UnregisterAsyncTask signals a background task has completed
func UnregisterAsyncTask() {
	eventLoopWg.Done()
}

// AwaitAsyncTasks blocks until all registered tasks complete
func AwaitAsyncTasks() {
	eventLoopWg.Wait()
}

// ResetAsyncTasks resets the async task counter for watch mode restarts
// This creates a new WaitGroup to start fresh
func ResetAsyncTasks() {
	eventLoopWg = sync.WaitGroup{}
}
