/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"net/url"
	"os"
	"os/exec"
)

/*
 * SDK for Titan remotes. Currently supports only client-side remote actions (parsing URIs and creating parameter
 * objects).
 */

type Remote interface {

	/*
	 * Returns the canonical name of this provider, such as "ssh" or "s3". This must be globally unique, and must
	 * match the leading URI component (ssh://...).
	 */
	Type() (string, error)

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
	remoteType, error := remote.Type()
	if error != nil {
		panic(error)
	}
	registeredRemotes[remoteType] = remote
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

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "titan",
	MagicCookieValue: "dba4fe2b-56ff-4a16-9bfc-bf651b8f12d6",
}

/*
 * Run the remote as a plugin server, to be invoked from the main method of the remote implementation.
 */
func Serve(remoteType string) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "remote",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	remote := Get(remoteType)
	var pluginMap = map[string]plugin.Plugin {
		"remote": &remotePlugin{Impl: remote},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig:  plugin.HandshakeConfig{},
		Plugins:          pluginMap,
		Logger:           logger,
	})
}

/*
 * Load a remote via the plugin interface. The caller should invoke the Kill() method on the client when complete.
 */
func Load(remoteType string, pluginPath string) (Remote, *plugin.Client, error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "remote",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	remote := Get(remoteType)
	var pluginMap = map[string]plugin.Plugin {
		"remote": &remotePlugin{Impl: remote},
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: pluginMap,
		Cmd: exec.Command(fmt.Sprintf("%s/%s", pluginPath, remoteType)),
		Logger: logger,
	})

	rpcClient, err := client.Client()
	if err != nil {
		return nil, client, err
	}

	raw, err := rpcClient.Dispense("remote")
	if err != nil {
		return nil, client, err
	}

	return raw.(Remote), client, nil
}
