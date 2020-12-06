package diego

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	bbsModels "code.cloudfoundry.org/bbs/models"

	"github.com/RageZBla/cf-stale-route-detector/models"
)

type DesiredLRPMapping map[models.ProcessGuid]map[string]string

type DesiredLRPMapper struct {
}

const routerKey = "cf-router"

func NewDesiredLRPMapper() *DesiredLRPMapper {
	return &DesiredLRPMapper{}
}

func (m DesiredLRPMapper) Map(path string) (DesiredLRPMapping, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	lrps, err := m.decodeFile(f)
	if err != nil {
		return nil, err
	}

	return m.mapLRPs(lrps)
}

func (m DesiredLRPMapper) decodeFile(f io.Reader) ([]bbsModels.DesiredLRP, error) {
	result := []bbsModels.DesiredLRP{}
	dec := json.NewDecoder(f)

	for {
		var desiredLRP bbsModels.DesiredLRP
		if err := dec.Decode(&desiredLRP); err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
		result = append(result, desiredLRP)
	}
}

func (m DesiredLRPMapper) mapLRPs(lrps []bbsModels.DesiredLRP) (DesiredLRPMapping, error) {
	var (
		processGUID models.ProcessGuid
		err         error
	)

	result := DesiredLRPMapping{}

	for _, lrp := range lrps {
		processGUID = models.ProcessGuid(lrp.ProcessGuid)
		result[processGUID] = m.extractMetricTags(lrp)
		if result[processGUID]["routes"], err = m.extractAppRoutes(lrp); err != nil {
			return nil, err
		}

	}

	return result, nil
}

func (m DesiredLRPMapper) extractAppInfo(lrp bbsModels.DesiredLRP) (map[string]string, error) {
	var err error
	result := m.extractMetricTags(lrp)

	if result["routes"], err = m.extractAppRoutes(lrp); err != nil {
		return nil, err
	}

	return result, nil
}

func (m DesiredLRPMapper) extractMetricTags(lrp bbsModels.DesiredLRP) map[string]string {
	result := map[string]string{}

	for k, v := range lrp.MetricTags {
		if v.Static != "" {
			result[k] = v.Static
		} else {
			result[k] = v.Dynamic.String()
		}
	}

	return result
}

func (m DesiredLRPMapper) extractAppRoutes(lrp bbsModels.DesiredLRP) (string, error) {
	routes := []string{}

	r := lrp.Routes
	if r == nil {
		return "", nil
	}

	jsonRoutes, ok := (*r)[routerKey]
	if !ok {
		return "", nil
	}

	var rawRoutes []struct {
		Hostnames        []string    `json:"hostnames"`
		Port             int         `json:"port"`
		RouteServiceURL  interface{} `json:"route_service_url"`
		IsolationSegment interface{} `json:"isolation_segment"`
	}

	if err := json.Unmarshal(*jsonRoutes, &rawRoutes); err != nil {
		return "", err
	}

	for _, r := range rawRoutes {
		for _, h := range r.Hostnames {
			routes = append(routes, h)
		}
	}

	return strings.Join(routes, ","), nil
}
