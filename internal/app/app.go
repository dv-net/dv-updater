package app

import (
	"context"
	"errors"
	"github.com/dv-net/dv-updater/internal/config"
	"github.com/dv-net/dv-updater/internal/distro"
	"github.com/dv-net/dv-updater/internal/server"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/pkg/logger"
	"net/http"
)

func Run(ctx context.Context, conf *config.Config, l logger.Logger, currentAppVersion, currentAppCommitHash string) error {
	d := distro.New(l)
	dist, err := d.DiscoverDistro()
	if err != nil {
		return err
	}

	svc, err := service.NewServices(l, dist, currentAppVersion, currentAppCommitHash)
	if err != nil {
		return err
	}

	if err = initTickers(ctx, svc, l, &conf.AutoUpdate); err != nil {
		return err
	}

	srv := server.NewServer(conf.HTTP, svc, l)

	l.Info("DV-Updater Server Start")

	if err := srv.Stop(); err != nil {
		l.Error("failed to stop server", err)
	}

	serverErrCh := make(chan error, 1)
	go func() {
		defer close(serverErrCh)
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case srvErr := <-serverErrCh:
		return srvErr
	}
}
