package domain

// HealthStatus represents the health status of a service component
type HealthStatus struct {
	Status string
	Error  string
}

// HealthChecker defines the interface for checking system health
type HealthChecker interface {
	Check() (*HealthStatus, error)
}

// DatabaseHealthChecker defines the interface for checking database health
type DatabaseHealthChecker interface {
	HealthChecker
	PingDB() error
}
