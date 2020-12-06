package presenters

import (
	"github.com/RageZBla/cf-stale-route-detector/detector"
	"github.com/RageZBla/cf-stale-route-detector/logger"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/formatted_presenter.go --fake-name Presenter . Presenter
type Presenter interface {
	StaleRoutes(map[string][]detector.StaleRoute, bool)
}

type DefaultPresenter struct {
	stdout logger.Logger
	stderr logger.Logger
}

func NewLoggerPresenter(stdout, stderr logger.Logger) *DefaultPresenter {
	return &DefaultPresenter{stdout, stderr}
}

func (p DefaultPresenter) StaleRoutes(staleRoutes map[string][]detector.StaleRoute, verbose bool) {
	if len(staleRoutes) == 0 {
		p.stdout.Println("No stale route detected")
		return
	}

	p.stderr.Println("Detected stale route(s)")
	if verbose {
		for domain, routes := range staleRoutes {
			p.stderr.Printf("domain: %s\n", domain)
			for _, route := range routes {
				p.stderr.Printf("app guid: %s\nDiego instance ID: %s\ncontainer endpoint: %s\n", route.AppID, route.DiegoInstanceID, route.ContainerEndpoint)
				if len(route.Extra) > 0 {
					p.stderr.Printf("extra:\n")
					for k, v := range route.Extra {
						p.stderr.Printf("  %s: %s\n", k, v)
					}
				}
			}
		}
	}
}
