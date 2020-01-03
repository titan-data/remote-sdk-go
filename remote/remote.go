/*
 * Copyright The Titan Project Contributors.
 */
package remote

import "net/url"

/*
 * SDK for Titan remotes. Currently supports only client-side remote actions (parsing URIs and creating parameter
 * objects).
 */

type Remote interface {

	/*
	 * Returns the canonical name of this provider, such as "ssh" or "s3". This must be globally unique, and must
	 * match the leading URI component (ssh://...).
	 */
	Type() string

	/*
	 * Parse a URL and return the provider-specific remote parameters in structured form. These properties will be
	 * preserved as part of the remote metadata on the server and passed to subsequent server-side operations. The
	 * additional properties map can contain properties specified by the user that don't fit the URI format well,
	 * such as "-p keyFile=/path/to/sshKey". This should return an error a bad URL format or invalid properties.
	 */
	FromURL(url *url.URL, additionalProperties map[string]string) (map[string]interface{}, error)

	/*
	 * Convert a remote back into URI form for display. Since this is for display only, any sensitive information
	 * should be redacted (e.g. "user:****@host"). Any properties that cannot be represented in the URI can be
	 * passed back as the second part of the pair.
	 */
	ToURL(properties map[string]interface{}) (string, map[string]string, error)

	/*
	 * Given a set of remote properties, return a set of parameter properties that will be passed to each operation.
	 * This is invoked in the context of the user CLI. It can access user data, such as SSH or AWS configuration. It
	 * can also interactively prompt the user for additional input (such as a password).
	 */
	GetParameters(remoteProperties map[string]interface{}) (map[string]interface{}, error)
}

var registeredRemotes = map[string]Remote{}

/*
 * Register a new remote. This should be called from the init() function of a remote implementation. The remotes can
 * later be accessed via the Get() method.
 */
func Register(remote Remote) {
	registeredRemotes[remote.Type()] = remote
}

/*
 * Get a remote by type.
 */
func Get(remoteType string) Remote {
	return registeredRemotes[remoteType]
}

/*
 * Clear any registered remotes. Should only be used for testing.
 */
func Clear() {
	registeredRemotes = map[string]Remote{}
}
