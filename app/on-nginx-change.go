package app

import (
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/romdo/go-debounce"
	"gitlab.com/CypriotUnknown/npm-to-pihole/constants"
	"gitlab.com/CypriotUnknown/npm-to-pihole/models"
	"gitlab.com/CypriotUnknown/npm-to-pihole/utils"
)

func (app *AppType) createExecuter() (debounced func(), cancel func()) {
	return debounce.New(3*time.Second, func() {

		proxyHosts := app.getNginxProxyHostServerNames()

		cnames := app.getPiholeCnames()

		piholeConfigRequests := make([]models.PiholeConfigRequest, 0)

		wg := sync.WaitGroup{}

		wg.Go(func() {
			for record, domain := range cnames {
				if !slices.Contains(proxyHosts, record) {
					piholeConfigRequests = append(piholeConfigRequests, models.PiholeConfigRequest{
						Cname: strings.Join([]string{record, domain}, ","),
						Type:  constants.PiholeDeleteRequest,
					})
				}
			}

		})

		wg.Go(func() {
			for _, host := range proxyHosts {
				if _, exists := cnames[host]; !exists {
					baseDomain, err := utils.GetBaseDomainFromHostString(host)
					if err != nil {
						slog.Error("invalid nginx host", "host", host, "error", err)
						continue
					}

					piholeConfigRequests = append(piholeConfigRequests, models.PiholeConfigRequest{
						Cname: strings.Join([]string{host, baseDomain}, ","),
						Type:  constants.PiholePutRequest,
					})
				}
			}
			wg.Done()
		})

		wg.Wait()

		app.editCnames(piholeConfigRequests)
	})

}
