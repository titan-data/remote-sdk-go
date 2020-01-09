/*
 * Copyright The Titan Project Contributors.
 */
package remote

import (
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
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
	 * such as "-p keyFile=/path/to/sshKey". This should return an error for a bad URL format or invalid properties.
	 * The calling context will have stripped out any query parameters or fragments.
	 */
	FromURL(url string, properties map[string]string) (map[string]interface{}, error)

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
	GetParameters(properties map[string]interface{}) (map[string]interface{}, error)
}

type loadedRemote struct {
	r Remote
	c *plugin.Client
}

var registeredRemotes = map[string]Remote{}
var loadedRemotes = map[string]loadedRemote{}

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
 * Clear any registered or loaded remotes. Should only be used for testing.
 */
func Clear() {
	registeredRemotes = map[string]Remote{}
	for _, v := range loadedRemotes {
		v.c.Kill()
	}
	loadedRemotes = map[string]loadedRemote{}
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
	remote := Get(remoteType)
	var pluginMap = map[string]plugin.Plugin{
		"remote": &remotePlugin{Impl: remote},
	}
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

/*
 * Load a remote via the plugin interface. These plugins will remain loaded until Unload() or Clear() is called.
 */
func Load(remoteType string, pluginPath string) (Remote, error) {
	if v, ok := loadedRemotes[remoteType]; ok {
		return v.r, nil
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "remote",
		Output: os.Stdout,
		Level:  hclog.Error,
	})

	remote := Get(remoteType)
	var pluginMap = map[string]plugin.Plugin{
		"remote": &remotePlugin{Impl: remote},
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		Cmd:              exec.Command(fmt.Sprintf("%s/%s", pluginPath, remoteType)),
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, err
	}

	raw, err := rpcClient.Dispense("remote")
	if err != nil {
		client.Kill()
		return nil, err
	}

	loadedRemotes[remoteType] = loadedRemote{
		r: raw.(Remote),
		c: client,
	}

	return raw.(Remote), nil
}

func Unload(remoteType string) {
	if val, ok := loadedRemotes[remoteType]; ok {
		val.c.Kill()
		delete(loadedRemotes, remoteType)
	}
}
