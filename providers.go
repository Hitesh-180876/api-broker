package main

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// IPInfoProvider implements the Provider interface for ipinfo.io
type IPInfoProvider struct {
	maxRequestsPerMinute int
	requestsThisMinute   int
}

func NewIPInfoProvider(maxRequestsPerMinute int) *IPInfoProvider {
	return &IPInfoProvider{
		maxRequestsPerMinute: maxRequestsPerMinute,
	}
}

func (p *IPInfoProvider) Name() string {
	return "ipinfo.io"
}

func (p *IPInfoProvider) GetLocation(ctx context.Context, ip string) (*Location, error) {
	// In a real implementation, this would make an HTTP request to ipinfo.io
	// For this example, we'll simulate the request with random latency and errors

	// Simulate network latency (50-300ms)
	latency := 50 + rand.Intn(250)
	select {
	case <-time.After(time.Duration(latency) * time.Millisecond):
		// Continue processing
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Simulate occasional errors (5% chance)
	if rand.Intn(100) < 5 {
		return nil, errors.New("ipinfo.io service error")
	}

	// Return simulated data
	return &Location{
		IP:      ip,
		Country: "United States",
		City:    "New York",
	}, nil
}

func (p *IPInfoProvider) GetRequestsThisMinute() int {
	return p.requestsThisMinute
}

func (p *IPInfoProvider) GetMaxRequestsPerMinute() int {
	return p.maxRequestsPerMinute
}

// IPAPIProvider implements the Provider interface for ip-api.com
type IPAPIProvider struct {
	maxRequestsPerMinute int
	requestsThisMinute   int
}

func NewIPAPIProvider(maxRequestsPerMinute int) *IPAPIProvider {
	return &IPAPIProvider{
		maxRequestsPerMinute: maxRequestsPerMinute,
	}
}

func (p *IPAPIProvider) Name() string {
	return "ip-api.com"
}

func (p *IPAPIProvider) GetLocation(ctx context.Context, ip string) (*Location, error) {
	// Simulate network latency (75-350ms)
	latency := 75 + rand.Intn(275)
	select {
	case <-time.After(time.Duration(latency) * time.Millisecond):
		// Continue processing
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Simulate occasional errors (7% chance)
	if rand.Intn(100) < 7 {
		return nil, errors.New("ip-api.com service error")
	}

	// Return simulated data
	return &Location{
		IP:      ip,
		Country: "Germany",
		City:    "Berlin",
	}, nil
}

func (p *IPAPIProvider) GetRequestsThisMinute() int {
	return p.requestsThisMinute
}

func (p *IPAPIProvider) GetMaxRequestsPerMinute() int {
	return p.maxRequestsPerMinute
}

// IPStackProvider implements the Provider interface for ipstack.com
type IPStackProvider struct {
	maxRequestsPerMinute int
	requestsThisMinute   int
}

func NewIPStackProvider(maxRequestsPerMinute int) *IPStackProvider {
	return &IPStackProvider{
		maxRequestsPerMinute: maxRequestsPerMinute,
	}
}

func (p *IPStackProvider) Name() string {
	return "ipstack.com"
}

func (p *IPStackProvider) GetLocation(ctx context.Context, ip string) (*Location, error) {
	// Simulate network latency (100-400ms)
	latency := 100 + rand.Intn(300)
	select {
	case <-time.After(time.Duration(latency) * time.Millisecond):
		// Continue processing
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Simulate occasional errors (10% chance)
	if rand.Intn(100) < 10 {
		return nil, errors.New("ipstack.com service error")
	}

	// Return simulated data
	return &Location{
		IP:      ip,
		Country: "Japan",
		City:    "Tokyo",
	}, nil
}

func (p *IPStackProvider) GetRequestsThisMinute() int {
	return p.requestsThisMinute
}

func (p *IPStackProvider) GetMaxRequestsPerMinute() int {
	return p.maxRequestsPerMinute
}

// In a real implementation, you would add actual HTTP client code to call the APIs
// Here's an example of what that might look like for a real provider:

/*
func (p *IPInfoProvider) GetLocation(ctx context.Context, ip string) (*Location, error) {
    url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("User-Agent", "IPLocationBroker/1.0")

    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
    }

    var result struct {
        IP      string `json:"ip"`
        Country string `json:"country"`
        City    string `json:"city"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &Location{
        IP:      result.IP,
        Country: result.Country,
        City:    result.City,
    }, nil
}
*/
