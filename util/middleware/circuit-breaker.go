package middleware

import (
	"time"

	"github.com/sony/gobreaker"
)

var CircuitBreaker *gobreaker.CircuitBreaker

func init() {
	var settings gobreaker.Settings
	settings.Name = "HTTP DELETE"
	settings.Timeout = time.Millisecond
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		if to == gobreaker.StateOpen {
			println("State Open!")
		}
		if from == gobreaker.StateOpen && to == gobreaker.StateHalfOpen {
			println("Going from Open to Half Open")
		}
		if from == gobreaker.StateHalfOpen && to == gobreaker.StateClosed {
			println("Going from Half Open to Closed!")
		}
	}
	CircuitBreaker = gobreaker.NewCircuitBreaker(settings)
}

func CircuitBreakerExecute(cb *gobreaker.CircuitBreaker, execFunc func() (interface{}, error)) (interface{}, error) {
	return cb.Execute(execFunc)
}
