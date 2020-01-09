package main

import (
	"github.com/titan-data/remote-sdk-go/internal/echo"
	"github.com/titan-data/remote-sdk-go/remote"
)

type MockRemote struct {
}

func main() {
	remote.Register(echo.EchoRemote{})
	remote.Serve("echo")
}
