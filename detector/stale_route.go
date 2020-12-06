package detector

import (
	"github.com/RageZBla/cf-stale-route-detector/diego"
	"github.com/RageZBla/cf-stale-route-detector/gorouter"
	"github.com/RageZBla/cf-stale-route-detector/models"

	"github.com/pkg/errors"
)

type StaleRoute struct {
	models.AppID
	models.ContainerEndpoint
	DiegoInstanceID models.InstanceID
	Extra           map[string]string
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o ./fakes/routing_table_parser.go --fake-name RoutingTableParser . routingTableParser
type routingTableParser interface {
	Parse(path string) (gorouter.RoutingEntries, error)
}

//counterfeiter:generate -o ./fakes/actual_lrps_mapper.go --fake-name ActualLRPMapper . actualLRPMapper
type actualLRPMapper interface {
	Map(path string) (diego.ActualLRPMapping, error)
}

//counterfeiter:generate -o ./fakes/desired_lrps_mapper.go --fake-name DesiredLRPMapper . desiredLRPMapper
type desiredLRPMapper interface {
	Map(path string) (diego.DesiredLRPMapping, error)
}

type StaleRouteDetector struct {
	routingTableParser routingTableParser
	desiredLRPMapper   desiredLRPMapper
	actualLRPMapper    actualLRPMapper
}

func NewStaleRouteDetector(routingTableParser routingTableParser, actualLRPMapper actualLRPMapper, desiredLRPMapper desiredLRPMapper) *StaleRouteDetector {
	return &StaleRouteDetector{
		routingTableParser: routingTableParser,
		desiredLRPMapper:   desiredLRPMapper,
		actualLRPMapper:    actualLRPMapper,
	}
}

func (d StaleRouteDetector) DetectFromFiles(routingTablePath, actualLRPPath, desiredLRPPath string) (staleRoutes map[string][]StaleRoute, err error) {

	var desiredLRPSMapping diego.DesiredLRPMapping

	routingTable, err := d.routingTableParser.Parse(routingTablePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse routing table")
	}

	lrpsMapping, err := d.actualLRPMapper.Map(actualLRPPath)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse actual LRPs")
	}

	if desiredLRPPath != "" {
		desiredLRPSMapping, err = d.desiredLRPMapper.Map(desiredLRPPath)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse desired LRPs")
		}
	}

	staleRoutes = d.Detect(routingTable, lrpsMapping, desiredLRPSMapping)
	return
}

func (d StaleRouteDetector) Detect(routingTable gorouter.RoutingEntries, lrpsMapping diego.ActualLRPMapping, desiredLRPSMapping diego.DesiredLRPMapping) (staleRoutes map[string][]StaleRoute) {
	staleRoutes = map[string][]StaleRoute{}

	for endpoint, routes := range routingTable {
		for _, routeDetails := range routes {
			lrp, ok := lrpsMapping[endpoint]
			extra := desiredLRPSMapping[lrp.ProcessGuid]
			if ok {
				if routeDetails.DiegoInstanceID != lrp.InstanceID {
					staleRoutes = d.append(staleRoutes, endpoint, routeDetails, lrp.InstanceID, extra)
				}
			}
		}
	}

	return
}

func (d StaleRouteDetector) append(staleRoutes map[string][]StaleRoute, endpoint models.ContainerEndpoint, routeEntry gorouter.RoutingEntry, diegoProcessID models.InstanceID, extra map[string]string) map[string][]StaleRoute {
	staleRoute := StaleRoute{
		AppID:             routeEntry.AppID,
		DiegoInstanceID:   diegoProcessID,
		ContainerEndpoint: endpoint,
		Extra:             extra,
	}
	staleRoutes[routeEntry.Route] = append(staleRoutes[routeEntry.Route], staleRoute)

	return staleRoutes
}
