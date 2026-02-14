package jwt

type Claims interface {
	Values() map[string]any

	Get(key string) any
	String(key string) string
	Int(key string) int
	Uint(key string) uint
	Float64(key string) float64
	Bool(key string) bool
	Int64(key string) int64
	StringSlice(key string) []string

	Issuer() string
	Subject() string
	Audience() string
	ExpirationTime() int64
	NotBefore() int64
	IssuedAt() int64
	ID() string

	Builder() ClaimsBuilder
}

// Claims represents the claims in a JWT
type ClaimsImpl struct {
	values map[string]any
}

// Values returns the values of the claims
func (c *ClaimsImpl) Values() map[string]any {
	return c.values
}

// Get returns the value of the claim for the given key
func (c *ClaimsImpl) Get(key string) any {
	return c.values[key]
}

// String returns the string value of the claim for the given key
func (c *ClaimsImpl) String(key string) string {
	s, ok := c.values[key].(string)
	if !ok {
		return ""
	}
	return s
}

// Int returns the int value of the claim for the given key
func (c *ClaimsImpl) Int(key string) int {
	i, ok := c.values[key].(int)
	if !ok {
		f, ok := c.values[key].(float64)
		if !ok {
			return 0
		}
		return int(f)
	}
	return i
}

// Uint returns the uint value of the claim for the given key
func (c *ClaimsImpl) Uint(key string) uint {
	u, ok := c.values[key].(uint)
	if !ok {
		i, ok := c.values[key].(float64)
		if !ok {
			return 0
		}
		return uint(i)
	}
	return u
}

// Float64 returns the float64 value of the claim for the given key
func (c *ClaimsImpl) Float64(key string) float64 {
	f, ok := c.values[key].(float64)
	if !ok {
		i, ok := c.values[key].(float64)
		if !ok {
			return 0
		}
		return float64(i)
	}
	return f
}

// Bool returns the bool value of the claim for the given key
func (c *ClaimsImpl) Bool(key string) bool {
	b, ok := c.values[key].(bool)
	if !ok {
		i, ok := c.values[key].(float64)
		if !ok {
			return false
		}
		return i != 0
	}
	return b
}

// Int64 returns the int64 value of the claim for the given key
func (c *ClaimsImpl) Int64(key string) int64 {
	i, ok := c.values[key].(int64)
	if !ok {
		i, ok := c.values[key].(float64)
		if !ok {
			return 0
		}
		return int64(i)
	}
	return i
}

// StringSlice returns the []string value of the claim for the given key
func (c *ClaimsImpl) StringSlice(key string) []string {
	val := c.values[key]
	if val == nil {
		return []string{}
	}
	
	// Try direct type assertion
	if strSlice, ok := val.([]string); ok {
		return strSlice
	}
	
	// Try []interface{} (common when unmarshaling JSON)
	if ifaceSlice, ok := val.([]interface{}); ok {
		result := make([]string, 0, len(ifaceSlice))
		for _, v := range ifaceSlice {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	
	return []string{}
}

// Issuer returns the issuer of the claims
func (c *ClaimsImpl) Issuer() string {
	return c.String(keyIssuer)
}

// Subject returns the subject of the claims
func (c *ClaimsImpl) Subject() string {
	return c.String(keySubject)
}

// Audience returns the audience of the claims
func (c *ClaimsImpl) Audience() string {
	return c.String(keyAudience)
}

// ExpirationTime returns the expiration time of the claims
func (c *ClaimsImpl) ExpirationTime() int64 {
	return c.Int64(keyExpirationTime)
}

// NotBefore returns the not before time of the claims
func (c *ClaimsImpl) NotBefore() int64 {
	return c.Int64(keyNotBefore)
}

// IssuedAt returns the issued at time of the claims
func (c *ClaimsImpl) IssuedAt() int64 {
	return c.Int64(keyIssuedAt)
}

// ID returns the unique identifier of the claims
func (c *ClaimsImpl) ID() string {
	return c.String(keyID)
}

// Builder returns the builder of the claims
func (c *ClaimsImpl) Builder() ClaimsBuilder {
	return &ClaimsBuilderImpl{
		values: c.values,
	}
}
