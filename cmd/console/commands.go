package console

import (
	"fmt" //nolint:goimports
	"github.com/dv-net/dv-updater/internal/app"
	"github.com/dv-net/dv-updater/internal/config"
	"github.com/dv-net/dv-updater/internal/distro"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/pkg/logger"
	"os" //nolint:goimports

	"github.com/tkcrm/mx/cfg"
	"github.com/urfave/cli/v2"
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
				if err := cfg.GenerateMarkdown(new(config.Config), "ENVS.md"); err != nil {
					return fmt.Errorf("failed to generate markdown: %w", err)
				}

				if err := cfg.GenerateYamlTemplate(new(config.Config), "configs/config.template.yaml"); err != nil {
					return fmt.Errorf("failed to generate yaml template: %w", err)
				}
				return nil
			},
		},
		{
			Name:  "flags",
			Usage: "print available config flags",
			Action: func(_ *cli.Context) error {
				res, err := cfg.GenerateFlags(new(config.Config))
				if err != nil {
					return err
				}

				fmt.Println(res)

				return nil
			},
		},
	}
}

func loadConfig(args, configPaths []string) (*config.Config, error) {
	conf := new(config.Config)
	if err := cfg.Load(conf,
		cfg.WithLoaderConfig(cfg.Config{
			Args:       args,
			Files:      configPaths,
			MergeFiles: true,
			SkipFiles:  true,
		}),
	); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return conf, nil
}
