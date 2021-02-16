package health

import (
	"context"
	"log"
	"time"
)

// Checker represents an health checker.
type Checker interface {
	Check(ctx context.Context) error
	Name() string
	TimeoutAfter() time.Duration
}

// HealthChecker holds the health checks.
type HealthChecker struct {
	checkers []Checker
}

// NewHealthChecker returns a new health checkers.
func NewHealthChecker(checkers ...Checker) HealthChecker {
	return HealthChecker{checkers: checkers}
}

// Check executes all the associated health checks till they complete successfully.
func (hc HealthChecker) Check(ctx context.Context) error {
	for _, c := range hc.checkers {
		ctx, cancel := context.WithTimeout(ctx, c.TimeoutAfter())
		defer cancel()

		for err := c.Check(ctx); err != nil; {
			log.Printf("failed health check %s: %v", c.Name(), err)
			ctx, cancel = context.WithTimeout(ctx, c.TimeoutAfter())
			defer cancel()
		}
	}

	return nil
}
