package httpcli

// Built in default client instance.
var defaultClient = New()

// Default function returns the built-in default client instance.
func Default() Client {
	return defaultClient
}
