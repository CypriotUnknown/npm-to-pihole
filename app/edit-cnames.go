package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gitlab.com/CypriotUnknown/npm-to-pihole/constants"
	"gitlab.com/CypriotUnknown/npm-to-pihole/models"
)

func (app *AppType) editCnames(requests []models.PiholeConfigRequest) {
	for i, request := range requests {
		endpoint := url.URL{
			Scheme: "https",
			Host:   "Pihole",
			Path:   fmt.Sprintf("/api/config/dns/cnameRecords/%s", request.Cname),
		}

		query := endpoint.Query()
		restartFTL := i == (len(requests) - 1)
		query.Set("restart", strconv.FormatBool(restartFTL))

		endpoint.RawQuery = query.Encode()

		cnamesRequest, err := http.NewRequestWithContext(app.ctx, string(request.Type), endpoint.String(), nil)
		if err != nil {
			slog.Error("error creating Pihole edit cnames request", "error", err)
			app.cancelCtx(fmt.Errorf("this is a bug. please create an issue"))
			return
		}

		cnamesRequest.Header.Set("X-FTL-SID", app.PiholeAuth.Session.Sid)

		cnamesHttpResponse, err := app.httpClient.Do(cnamesRequest)
		if err != nil {
			slog.Error("error making Pihole edit cnames request")
			app.cancelCtx(err)
			return
		}

		var expectedStatus int
		var infoMessage string

		switch request.Type {
		case constants.PiholeDeleteRequest:
			expectedStatus = 204
			infoMessage = "removed record from Pihole"
		case constants.PiholePutRequest:
			expectedStatus = 201
			infoMessage = "added record to Pihole"
		}

		if cnamesHttpResponse.StatusCode != expectedStatus {
			errorData, _ := io.ReadAll(cnamesHttpResponse.Body)

			var errorMap map[string]any
			json.Unmarshal(errorData, &errorMap)

			slog.Error("could not edit cnames", "pihole_response", errorMap)
		}

		slog.Info(infoMessage, "cname", strings.Split(request.Cname, ",")[0])
	}
}
