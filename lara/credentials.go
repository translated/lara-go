package lara

// Credentials is deprecated. Use AccessKey instead.
// Deprecated: Use AccessKey instead.
type Credentials struct {
	*AccessKey
}

// NewCredentials creates a new Credentials (deprecated).
// Deprecated: Use NewAccessKey instead.
func NewCredentials(accessKeyID, accessKeySecret string) *Credentials {
	return &Credentials{
		AccessKey: NewAccessKey(accessKeyID, accessKeySecret),
	}
}
