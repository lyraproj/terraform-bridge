package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/lyraproj/issue/issue"
	"github.com/lyraproj/terraform-bridge/cmd/goplugin-kubernetes/handler"
)

func init() {
	// Configuring hclog like this allows Lyra to handle log levels automatically
	hclog.DefaultOptions = &hclog.LoggerOptions{
		Name:            "kubernetes",
		Level:           hclog.LevelFromString(os.Getenv("LYRA_LOG_LEVEL")),
		JSONFormat:      true,
		IncludeLocation: false,
		Output:          os.Stderr,
	}
	issue.IncludeStacktrace(hclog.DefaultOptions.Level <= hclog.Debug)
}

func main() {
	handler.Start()
}
