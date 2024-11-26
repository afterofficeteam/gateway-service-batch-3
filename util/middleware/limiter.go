package middleware

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type UserLimiter struct {
	Limiter  *rate.Limiter // Rate limiter for the user
	LastSeen time.Time     // Timestamp of the last user activity
}

var (
	// Map that associates user IDs with their respective rate limiters.
	limiters = make(map[string]*UserLimiter)
	// Mutex used to protect concurrent access to the limiters map.
	mutex sync.Mutex
)

const cleanupInterval = time.Minute * 5 // Duration after which inactive limiters are cleaned up

func GetLimiter(userID string) *rate.Limiter {
	// Lock the mutex to ensure thread-safe access to the limiters map.
	mutex.Lock()
	defer mutex.Unlock() // Unlock after the function completes.

	// Check if the limiter already exists for the given user ID.
	if limiter, exists := limiters[userID]; exists {
		// If the limiter exists, update the last seen time and return the limiter.
		limiter.LastSeen = time.Now()
		return limiter.Limiter
	}

	// If no limiter exists for the user, create a new rate limiter.
	// The rate limiter allows 5 requests every 5 minutes.
	limiter := rate.NewLimiter(rate.Every(time.Minute*5), 5)

	// Store the new limiter in the map along with the current timestamp as LastSeen.
	limiters[userID] = &UserLimiter{
		Limiter:  limiter,
		LastSeen: time.Now(),
	}

	// Return the newly created limiter.
	return limiter
}

func CleanupOldLimiters() {
	for {
		// Sleep for the duration of the cleanup interval before performing the cleanup.
		time.Sleep(cleanupInterval)

		// Lock the mutex to safely modify the limiters map.
		mutex.Lock()

		// Iterate through all limiters to check for inactivity.
		for userID, limiter := range limiters {
			// If the user has been inactive for longer than the cleanup interval, remove the limiter.
			if time.Since(limiter.LastSeen) > cleanupInterval {
				delete(limiters, userID)
			}
		}

		// Unlock the mutex after cleanup is complete.
		mutex.Unlock()
	}
}
