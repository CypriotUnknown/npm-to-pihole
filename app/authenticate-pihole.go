package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"gitlab.com/CypriotUnknown/npm-to-pihole/models"
)

func (app *AppType) authenticatePihole() {
	requestURL := url.URL{
		Scheme: "https",
		Host:   "Pihole",
		Path:   "/api/auth",
	}

	body := struct {
		Password string `json:"password"`
	}{Password: app.PiholePassword}

	bodyData, _ := json.Marshal(body)

	authRequest, err := http.NewRequestWithContext(app.ctx, "POST", requestURL.String(), bytes.NewReader(bodyData))

	if err != nil {
		slog.Error("error creating Pihole auth request", "error", err)
		app.cancelCtx(fmt.Errorf("this is a bug. please create an issue"))
		return
	}

	authHttpResponse, err := app.httpClient.Do(authRequest)

	if err != nil {
		slog.Error("error making Pihole auth request")
		app.cancelCtx(err)
		return
	}

	authResponseData, err := io.ReadAll(authHttpResponse.Body)

	if err != nil {
		slog.Error("could not read Pihole auth response")
		app.cancelCtx(err)
		return
	}

	authResponse, err := models.UnmarshalPiholeAuthResponse(authResponseData)
	if err != nil {
		slog.Error("bad response from Pihole auth request")
		app.cancelCtx(err)
		return
	}

	slog.Info("Pihole authentication success")

	app.mutex.Lock()
	app.PiholeAuth = authResponse
	app.mutex.Unlock()

	go app.handleAuth()
}
