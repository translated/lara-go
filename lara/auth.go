package lara

// AccessKey represents API key authentication credentials
type AccessKey struct {
	ID     string
	Secret string
}

// NewAccessKey creates a new AccessKey
func NewAccessKey(id, secret string) *AccessKey {
	return &AccessKey{
		ID:     id,
		Secret: secret,
	}
}

// AuthToken represents JWT token authentication
type AuthToken struct {
	Token        string
	RefreshToken string
}

// NewAuthToken creates a new AuthToken
func NewAuthToken(token, refreshToken string) *AuthToken {
	return &AuthToken{
		Token:        token,
		RefreshToken: refreshToken,
	}
}
