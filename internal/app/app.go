package app

import (
	"context"
	"dv-updater/internal/config"
	"dv-updater/internal/distro"
	"dv-updater/internal/server"
	"dv-updater/internal/service"
	"dv-updater/pkg/logger"
	"errors"
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
