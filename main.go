package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Location represents the geographical location data
type Location struct {
	IP      string
	Country string
	City    string
}

// Provider interface for IP location services
type Provider interface {
	Name() string
	GetLocation(ctx context.Context, ip string) (*Location, error)
	GetRequestsThisMinute() int
	GetMaxRequestsPerMinute() int
}

// ProviderStats tracks quality metrics for a provider
type ProviderStats struct {
	provider            Provider
	mutex               sync.RWMutex
	errorsInLast5Min    []time.Time
	responseTimes       []time.Duration
	responseTimesMutex  sync.RWMutex
	requestsThisMinute  int
	requestsMinuteReset time.Time
}

// Broker manages multiple providers and routes requests
type Broker struct {
	providers     []*ProviderStats
	providerMutex sync.RWMutex
}

// NewBroker creates a new broker with the given providers
func NewBroker(providers []Provider) *Broker {
	broker := &Broker{
		providers: make([]*ProviderStats, len(providers)),
	}

	for i, p := range providers {
		broker.providers[i] = &ProviderStats{
			provider:            p,
			errorsInLast5Min:    make([]time.Time, 0),
			responseTimes:       make([]time.Duration, 0),
			requestsThisMinute:  0,
			requestsMinuteReset: time.Now(),
		}
	}

	// Start a goroutine to clean up old stats
	go broker.cleanupStatsRoutine()

	return broker
}

// cleanupStatsRoutine periodically cleans up old stats
func (b *Broker) cleanupStatsRoutine() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		b.cleanupStats()
	}
}

// cleanupStats removes old error and response time entries
func (b *Broker) cleanupStats() {
	fiveMinAgo := time.Now().Add(-5 * time.Minute)

	b.providerMutex.RLock()
	defer b.providerMutex.RUnlock()

	for _, ps := range b.providers {
		ps.mutex.Lock()

		// Clean up errors older than 5 minutes
		newErrors := make([]time.Time, 0)
		for _, t := range ps.errorsInLast5Min {
			if t.After(fiveMinAgo) {
				newErrors = append(newErrors, t)
			}
		}
		ps.errorsInLast5Min = newErrors

		// Clean up response times older than 5 minutes
		ps.responseTimesMutex.Lock()
		if len(ps.responseTimes) > 100 {
			// Keep only the most recent 100 response times
			ps.responseTimes = ps.responseTimes[len(ps.responseTimes)-100:]
		}
		ps.responseTimesMutex.Unlock()

		// Reset requests counter if a minute has passed
		if time.Since(ps.requestsMinuteReset) > time.Minute {
			ps.requestsThisMinute = 0
			ps.requestsMinuteReset = time.Now()
		}

		ps.mutex.Unlock()
	}
}

// GetLocation returns the location for an IP using the best available provider
func (b *Broker) GetLocation(ctx context.Context, ip string) (*Location, error) {
	bestProvider := b.selectBestProvider()
	if bestProvider == nil {
		return nil, errors.New("no suitable provider available")
	}

	// Track request start time
	startTime := time.Now()

	// Update request count
	bestProvider.mutex.Lock()
	bestProvider.requestsThisMinute++
	bestProvider.mutex.Unlock()

	// Make the request to the provider
	location, err := bestProvider.provider.GetLocation(ctx, ip)

	// Record response time
	responseTime := time.Since(startTime)
	bestProvider.responseTimesMutex.Lock()
	bestProvider.responseTimes = append(bestProvider.responseTimes, responseTime)
	bestProvider.responseTimesMutex.Unlock()

	// Record error if any
	if err != nil {
		bestProvider.mutex.Lock()
		bestProvider.errorsInLast5Min = append(bestProvider.errorsInLast5Min, time.Now())
		bestProvider.mutex.Unlock()
		return nil, err
	}

	return location, nil
}

// selectBestProvider chooses the most reliable provider based on metrics
func (b *Broker) selectBestProvider() *ProviderStats {
	b.providerMutex.RLock()
	defer b.providerMutex.RUnlock()

	var bestProvider *ProviderStats
	var bestScore float64 = -1

	for _, ps := range b.providers {
		ps.mutex.RLock()

		// Skip if provider is at or over rate limit
		if ps.requestsThisMinute >= ps.provider.GetMaxRequestsPerMinute() {
			ps.mutex.RUnlock()
			continue
		}

		// Calculate error rate (lower is better)
		errorRate := float64(len(ps.errorsInLast5Min)) / 300.0 // errors per second in last 5 min

		// Calculate average response time
		ps.responseTimesMutex.RLock()
		var avgResponseTime float64
		if len(ps.responseTimes) > 0 {
			var total time.Duration
			for _, rt := range ps.responseTimes {
				total += rt
			}
			avgResponseTime = float64(total) / float64(len(ps.responseTimes))
		}
		ps.responseTimesMutex.RUnlock()

		// Calculate capacity left (higher is better)
		capacityLeft := 1.0 - (float64(ps.requestsThisMinute) / float64(ps.provider.GetMaxRequestsPerMinute()))

		// Calculate score (higher is better)
		// We prioritize providers with lower error rates and faster response times
		// while also considering available capacity
		score := (1.0 - errorRate) * (1000.0 / (avgResponseTime + 1.0)) * capacityLeft

		ps.mutex.RUnlock()

		if bestScore < 0 || score > bestScore {
			bestScore = score
			bestProvider = ps
		}
	}

	return bestProvider
}

func main() {
	// Example usage
	providers := []Provider{
		NewIPInfoProvider(100),  // 100 requests per minute
		NewIPAPIProvider(120),   // 120 requests per minute
		NewIPStackProvider(150), // 150 requests per minute
	}

	broker := NewBroker(providers)

	// Set up HTTP server
	http.HandleFunc("/location", func(w http.ResponseWriter, r *http.Request) {
		ip := r.URL.Query().Get("ip")
		if ip == "" {
			http.Error(w, "IP parameter is required", http.StatusBadRequest)
			return
		}

		location, err := broker.GetLocation(r.Context(), ip)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting location: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "IP: %s\nCountry: %s\nCity: %s\n",
			location.IP, location.Country, location.City)
	})

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
