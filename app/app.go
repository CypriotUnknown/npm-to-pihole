package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"gitlab.com/CypriotUnknown/npm-to-pihole/models"
)

type AppType struct {
	mutex          sync.Mutex
	ctx            context.Context
	cancelCtx      context.CancelCauseFunc
	PiholePassword string
	mginxProxyDir  string
	PiholeAuth     *models.PiholeAuthResponse
	httpClient     *http.Client
	authTimer      *time.Timer
}

func (app *AppType) handleAuth() {
	app.authTimer = time.AfterFunc(time.Second*time.Duration(app.PiholeAuth.Session.Validity-10), app.authenticatePihole)
}

func Start() {
	password, passwordExists := os.LookupEnv("PIHOLE_PASSWORD")
	if !passwordExists {
		panic("'PIHOLE_PASSWORD' not set")
	}

	nginxProxyDir, nginxProxyDirExists := os.LookupEnv("NGINX_PROXY_DIR")
	if !nginxProxyDirExists {
		const nginx_dir = "/data/nginx/nginx/proxy_host"
		nginxProxyDir = nginx_dir
	}

	appContext, appContextCancel := context.WithCancelCause(context.Background())
	app := AppType{
		PiholePassword: password,
		mutex:          sync.Mutex{},
		ctx:            appContext,
		cancelCtx:      appContextCancel,
		mginxProxyDir:  nginxProxyDir,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	defer func() {
		if app.authTimer != nil {
			app.authTimer.Stop()
		}
	}()

	watcher := app.createNginxProxyHostsDirectoryWatcher()
	defer watcher.Close()

	executeChanges, cancelExecuter := app.createExecuter()
	defer cancelExecuter()

	runOnStart := os.Getenv("RUN_ON_START")
	if runOnStart == "true" {
		executeChanges()
	} else if runOnStart != "" && runOnStart != "false" {
		panic("invalid value for 'RUN_ON_START'. if specified, it can only be true or false.")
	}

	go app.authenticatePihole()

	slog.Info("Watching for changes in Nginx proxy configuration files...")

	for {
		select {
		case <-signalChannel:
			app.cancelCtx(fmt.Errorf("interrupt signal received"))
		case <-app.ctx.Done():
			if err := app.ctx.Err(); err != nil {
				slog.Error("Quitting", "reason", context.Cause(app.ctx), "error", err)
				return
			}
		case err, ok := <-watcher.Errors:
			if ok && err != nil {
				slog.Error("error occured watching nginx proxy hosts directory", "error", err)
			}
		case event, ok := <-watcher.Events:
			if !ok {
				app.cancelCtx(fmt.Errorf("no longer watching nginx proxy hosts directory"))
			} else {
				slog.Info("change in proxy hosts directory", "path", event.Name, "operation", event.Op.String())
				executeChanges()
			}

		}
	}
}
