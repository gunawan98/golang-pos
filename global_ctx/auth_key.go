package globalctx

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// // userIDKey is the key for storing the user ID in the context.
const userIDKey contextKey = "userId"

// UserIDKey returns the key for userId to be used in context.
func UserIDKey() contextKey {
	return userIDKey
}
