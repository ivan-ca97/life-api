package ports

type GoogleClaims struct {
	Subject string
	Email   string
}

type GoogleTokenVerifier interface {
	Verify(idToken, clientId string) (*GoogleClaims, error)
}
