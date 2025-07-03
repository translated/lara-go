package lara

type Credentials struct {
	accessKeyID     string
	accessKeySecret string
}

func NewCredentials(accessKeyID, accessKeySecret string) *Credentials {
	return &Credentials{
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
	}
}
