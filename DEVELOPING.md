# Project Development

For general information about contributing changes, see the
[Contributor Guidelines](https://github.com/titan-data/.github/blob/master/CONTRIBUTING.md).

## How it Works

The Remote SDK currently provides interfaces only for use by the client. This includes the ability to register
remote providers, and parse URIs. Future work will add server-side capabilities for use in titan-server, enabling
the EOL of the legacy kotlin remote providers.

## Building

Run `go build -v ./...`.

## Testing

Run `go test -v ./...`.

## Releasing

Push a tag of the form `v<X>.<Y>.<Z>`, and publish the draft release in GitHub.
