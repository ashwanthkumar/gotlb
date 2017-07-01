package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobinStrategy(t *testing.T) {
	s := RoundRobinStrategy()
	s.AddBackend("a")
	s.AddBackend("b")
	s.AddBackend("c")
	assert.Equal(t, "a", s.Next())
	assert.Equal(t, "b", s.Next())
	assert.Equal(t, "c", s.Next())
	// We should start over again
	assert.Equal(t, "a", s.Next())
	assert.Equal(t, "b", s.Next())
	assert.Equal(t, "c", s.Next())
}

func TestRoundRobinStrategyUponRemovingBackend(t *testing.T) {
	s := RoundRobinStrategy()
	s.AddBackend("a")
	s.AddBackend("b")
	s.AddBackend("c")
	assert.Equal(t, "a", s.Next())
	s.RemoveBackend("b")
	assert.Equal(t, "c", s.Next())
	assert.Equal(t, "a", s.Next())
	assert.Equal(t, "c", s.Next())
}
