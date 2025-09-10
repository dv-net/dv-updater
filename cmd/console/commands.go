package console

import (
	"bytes"
	"fmt"
	"os"

	"github.com/dv-net/dv-updater/internal/app"
	"github.com/dv-net/dv-updater/internal/config"
	"github.com/dv-net/dv-updater/internal/distro"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/pkg/logger"
	"github.com/dv-net/xconfig"
	"github.com/goccy/go-yaml"

	"github.com/urfave/cli/v2"
)

const (
	envPrefix = "UPDATER"
)

func InitCommands(currentAppVersion, currentAppCommitHash string) []*cli.Command {
	return []*cli.Command{
		{
			Name:        "start",
			Description: "DV updater server",
			Action: func(ctx *cli.Context) error {
				conf, err := loadConfig(ctx.Args().Slice(), ctx.StringSlice("configs"))
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}

				l := logger.New(currentAppVersion, conf.Log)
				l.Info("Logger Init")
				return app.Run(ctx.Context, conf, l, currentAppVersion, currentAppCommitHash)
			},
		},
		{
			Name:        "self-update",
			Description: "Dv self update updater",
			Action: func(ctx *cli.Context) error {
				conf, err := loadConfig(ctx.Args().Slice(), ctx.StringSlice("configs"))
				if err != nil {
					return fmt.Errorf("failed to load config: %w", err)
				}
				l := logger.New(currentAppVersion, conf.Log)
				l.Info("Logger Init")
				d := distro.New(l)
				dist, err := d.DiscoverDistro()
				if err != nil {
					return err
				}

				svc, err := service.NewServices(l, dist, currentAppVersion, currentAppCommitHash)
				if err != nil {
					return err
				}

				if err = app.SelfUpdate(ctx.Context, &conf.AutoUpdate, svc, l); err != nil {
					l.Error("self update failed", err)
				}

				return nil
			},
		},
		{
			Name:        "version",
			Description: "print DV updater server version",
			Action: func(_ *cli.Context) error {
				_, _ = fmt.Fprintln(os.Stdout, currentAppVersion)
				return nil
			},
		},
		{
			Name:        "config",
			Description: "validate, gen envs and flags for config",
			Subcommands: prepareConfigCommands(),
		}, // config
	}
}

func prepareConfigCommands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "genenvs",
			Usage: "generate markdown for all envs and config yaml template",
			Action: func(_ *cli.Context) error {
				conf := new(config.Config)
				envMarkdown, err := xconfig.GenerateMarkdown(conf, xconfig.WithEnvPrefix(envPrefix))
				if err != nil {
					return fmt.Errorf("failed to generate markdown: %w", err)
				}
				envMarkdown = fmt.Sprintf("# Environment variables\n\nAll envs have prefix `%s_`\n\n%s", envPrefix, envMarkdown)
				if err := os.WriteFile("ENVS.md", []byte(envMarkdown), 0o600); err != nil {
					return err
				}

				buf := bytes.NewBuffer(nil)
				enc := yaml.NewEncoder(buf, yaml.Indent(2))
				defer enc.Close()

				if err := enc.Encode(conf); err != nil {
					return fmt.Errorf("failed to encode yaml: %w", err)
				}

				if err := os.WriteFile("configs/config.template.yaml", buf.Bytes(), 0o600); err != nil {
					return fmt.Errorf("failed to write file: %w", err)
				}

				return nil
			},
		},
	}
}

func loadConfig(_, configPaths []string) (*config.Config, error) {
	conf, err := config.Load[config.Config](configPaths, envPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return conf, nil
}
