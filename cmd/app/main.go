package main

import (
	"fmt"
	"os"
	"runtime"

	"dv-updater/cmd/console"

	"github.com/urfave/cli/v2"
)

var (
	appName    = "dv-updater"
	version    = "local"
	commitHash = "unknown"
)

func main() {
	application := &cli.App{
		Name:                 appName,
		Description:          "This is an API for DV updater",
		Version:              getBuildVersion(),
		Suggest:              true,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			cli.HelpFlag,
			cli.VersionFlag,
			cli.BashCompletionFlag,
		},
		Commands: console.InitCommands(version, commitHash),
	}
	if err := application.Run(os.Args); err != nil {
		_, _ = fmt.Println(err.Error())
		os.Exit(1)
	}
}

func getBuildVersion() string {
	return fmt.Sprintf(
		"\n\nrelease: %s\ncommit hash: %s\ngo version: %s",
		version,
		commitHash,
		runtime.Version(),
	)
}
