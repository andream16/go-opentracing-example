package health

// Retrier says whether an error is retryable or not.
type Retrier interface {
	error
	IsRetriable() bool
}
