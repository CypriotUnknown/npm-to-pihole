package app

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/CypriotUnknown/npm-to-pihole/models"
)

func (app *AppType) getPiholeCnames() map[string]string {
	cnamesURL, err := url.Parse(app.PiholeURL)

	if err != nil {
		slog.Error("error parsing url", "error", err)
		panic("could not join path for url")
	}

	cnamesURL.Path = "api/config/dns/cnameRecords"

	cnamesRequest, err := http.NewRequestWithContext(app.ctx, "GET", cnamesURL.String(), nil)
	if err != nil {
		slog.Error("error creating Pihole get cnames request", "error", err)
		app.cancelCtx(fmt.Errorf("this is a bug. please create an issue"))
		return nil
	}

	cnamesRequest.Header.Set("X-FTL-SID", app.PiholeAuth.Session.Sid)

	cnamesHttpResponse, err := app.httpClient.Do(cnamesRequest)
	if err != nil {
		slog.Error("error making Pihole get cnames request")
		app.cancelCtx(err)
		return nil
	}

	cnamesResponseData, err := io.ReadAll(cnamesHttpResponse.Body)
	if err != nil {
		slog.Error("could not read Pihole get cnames response")
		app.cancelCtx(err)
		return nil
	}

	piholeConfigResponse, err := models.UnmarshalPiholeGetConfigRespone(cnamesResponseData)
	if err != nil {
		slog.Error("bad response from Pihole get cnames request")
		app.cancelCtx(err)
		return nil
	}

	cnameRecords := make(map[string]string, len(piholeConfigResponse.Config.DNS.CnameRecords))

	for _, record := range piholeConfigResponse.Config.DNS.CnameRecords {
		split := strings.Split(record, ",")
		if len(split) != 2 {
			slog.Error("corrupt pihole data. cname record must be split by a ',' holding both the record and domain")
			app.cancelCtx(nil)
			return nil
		}

		cnameRecords[split[0]] = split[1]
	}

	return cnameRecords
}
