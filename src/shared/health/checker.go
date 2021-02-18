package health

import (
	"context"
	"fmt"
	"log"
)

// Checker represents an health checker.
type Checker interface {
	Check(ctx context.Context) error
	Name() string
}

// Manager holds the health checks.
type Manager struct {
	checkers []Checker
}

// NewManager returns a new health checkers manager.
func NewManager(checkers ...Checker) Manager {
	return Manager{checkers: checkers}
}

// Check executes all the associated health checks until they complete successfully.
func (m Manager) Check(ctx context.Context) error {
	for _, c := range m.checkers {
		for err := c.Check(ctx); err != nil; {
			log.Printf("failed health check %s: %v", c.Name(), err)
			if _, ok := err.(Retrier); !ok {
				return fmt.Errorf("failed health check %s: %w", c.Name(), err)
			}
		}
	}
	return nil
}
