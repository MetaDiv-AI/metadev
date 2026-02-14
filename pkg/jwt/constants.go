package jwt

const (
	keyIssuer         = "iss" // who issued the token
	keySubject        = "sub" // subject of the token
	keyAudience       = "aud" // who or what the token is intended for
	keyExpirationTime = "exp" // expiration time, as unix time
	keyNotBefore      = "nbf" // not before, as unix time
	keyIssuedAt       = "iat" // issued at, as unix time
	keyID             = "jti" // unique identifier for the token
)
