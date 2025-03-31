package mirror

// Auth represents authentication credentials for a registry
type Auth struct {
	Username string
	Password string
}

// Registry represents a container registry mirror configuration
type Registry struct {
	Source  string   // The registry to be replaced (e.g., "docker.io")
	Mirrors []string // The mirror endpoints (e.g., ["reg.footprint-ai.com"])
	Auth    *Auth    // Optional authentication information
}

// Getter defines the interface to get registry mirror configuration
type Getter interface {
	GetRegistry() []Registry
}
