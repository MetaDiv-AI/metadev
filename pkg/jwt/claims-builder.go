package jwt

// NewClaimsBuilder creates a new ClaimsBuilder
// it returns a ClaimsBuilder interface that can be used to build claims
func NewClaimsBuilder() ClaimsBuilder {
	return &ClaimsBuilderImpl{
		values: make(map[string]any),
	}
}

type ClaimsBuilder interface {
	// Issuer sets the issuer of the claims
	// issuer is about who issued the token
	// issuer is the standard claim key: "iss"
	Issuer(issuer string) ClaimsBuilder
	// Subject sets the subject of the claims
	// subject is about who or what the token is intended for
	// subject is the standard claim key: "sub"
	Subject(subject string) ClaimsBuilder
	// Audience sets the audience of the claims
	// audience is about who or what the token is intended for
	// audience is the standard claim key: "aud"
	Audience(audience string) ClaimsBuilder
	// ExpirationTime sets the expiration time of the claims
	// expirationTime is the time when the token will expire
	// expirationTime is the standard claim key: "exp"
	ExpirationTime(expirationTime int64) ClaimsBuilder
	// NotBefore sets the not before time of the claims
	// notBefore is the time before which the token is not valid
	// notBefore is the standard claim key: "nbf"
	NotBefore(notBefore int64) ClaimsBuilder
	// IssuedAt sets the issued at time of the claims
	// issuedAt is the time when the token was issued
	// issuedAt is the standard claim key: "iat"
	IssuedAt(issuedAt int64) ClaimsBuilder
	// ID sets the unique identifier of the claims
	// id is the unique identifier of the token
	// id is the standard claim key: "jti"
	ID(id string) ClaimsBuilder
	// Key sets the key of the claims
	// key is the key of the claim and then set the value of the claim
	Key(key string) *claimsBuilderKeySetter
	// Build builds the claims
	Build() Claims
}

// ClaimsBuilderImpl is the implementation of the ClaimsBuilder interface
type ClaimsBuilderImpl struct {
	values map[string]any
}

// Issuer sets the issuer of the claims
// issuer is about who issued the token
// issuer is the standard claim key: "iss"
func (c *ClaimsBuilderImpl) Issuer(issuer string) ClaimsBuilder {
	c.values[keyIssuer] = issuer
	return c
}

// Subject sets the subject of the claims
// subject is about who or what the token is intended for
// subject is the standard claim key: "sub"
func (c *ClaimsBuilderImpl) Subject(subject string) ClaimsBuilder {
	c.values[keySubject] = subject
	return c
}

// Audience sets the audience of the claims
// audience is about who or what the token is intended for
// audience is the standard claim key: "aud"
func (c *ClaimsBuilderImpl) Audience(audience string) ClaimsBuilder {
	c.values[keyAudience] = audience
	return c
}

// ExpirationTime sets the expiration time of the claims
// expirationTime is the time when the token will expire
// expirationTime is the standard claim key: "exp"
func (c *ClaimsBuilderImpl) ExpirationTime(expirationTime int64) ClaimsBuilder {
	c.values[keyExpirationTime] = expirationTime
	return c
}

// NotBefore sets the not before time of the claims
// notBefore is the time before which the token is not valid
// notBefore is the standard claim key: "nbf"
func (c *ClaimsBuilderImpl) NotBefore(notBefore int64) ClaimsBuilder {
	c.values[keyNotBefore] = notBefore
	return c
}

// IssuedAt sets the issued at time of the claims
// issuedAt is the time when the token was issued
// issuedAt is the standard claim key: "iat"
func (c *ClaimsBuilderImpl) IssuedAt(issuedAt int64) ClaimsBuilder {
	c.values[keyIssuedAt] = issuedAt
	return c
}

// ID sets the unique identifier of the claims
// id is the unique identifier of the token
// id is the standard claim key: "jti"
func (c *ClaimsBuilderImpl) ID(id string) ClaimsBuilder {
	c.values[keyID] = id
	return c
}

// Key sets the key of the claims
// key is the key of the claim and then set the value of the claim
func (c *ClaimsBuilderImpl) Key(key string) *claimsBuilderKeySetter {
	return &claimsBuilderKeySetter{
		claimsBuilder: c,
		key:           key,
	}
}

// Build builds the claims
func (c *ClaimsBuilderImpl) Build() Claims {
	return &ClaimsImpl{
		values: c.values,
	}
}

type claimsBuilderKeySetter struct {
	claimsBuilder *ClaimsBuilderImpl
	key           string
}

// Set sets the value of the claim
func (c *claimsBuilderKeySetter) Set(value any) ClaimsBuilder {
	c.claimsBuilder.values[c.key] = value
	return c.claimsBuilder
}
