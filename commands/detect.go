package commands

import (
	"errors"

	"github.com/RageZBla/cf-stale-route-detector/detector"
	"github.com/RageZBla/cf-stale-route-detector/presenters"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o fakes/service.go --fake-name Service . service
type service interface {
	DetectFromFiles(routingTablePath, actualLRPPath, desiredLRPPath string) (staleRoutes map[string][]detector.StaleRoute, err error)
}

type Detect struct {
	service   service
	presenter presenters.Presenter
	logger    logger
	Options   struct {
		RoutingTablePath string `long:"routing-table" required:"true" description:"Gorouter routing table export"`
		ActualLRPSPath   string `long:"actual-lrps" required:"true" description:"Diego actual LRPS export"`
		DesiredLRPSPath  string `long:"desired-lrps" description:"Diego desired LRPS export"`
		Verbose          bool   `long:"verbose" required:"false" description:"print details about stale route(s)"`
	}
}

var ErrStaleRouteDetected = errors.New("stale route(s) detected")

func NewDetect(service service, presenter presenters.Presenter) *Detect {
	return &Detect{
		service:   service,
		presenter: presenter,
	}
}

func (d Detect) Execute(args []string) error {
	staleRoutes, err := d.service.DetectFromFiles(d.Options.RoutingTablePath, d.Options.ActualLRPSPath, d.Options.DesiredLRPSPath)

	if err != nil {
		return err
	}

	d.presenter.StaleRoutes(staleRoutes, d.Options.Verbose)

	if len(staleRoutes) != 0 {
		return ErrStaleRouteDetected
	}

	return nil
}
