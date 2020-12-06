package gorouter

import (
	"encoding/json"
	"io"
	"os"

	"github.com/RageZBla/cf-stale-route-detector/models"
)

type RoutingEntries map[models.ContainerEndpoint][]RoutingEntry
type RoutingEntry struct {
	models.AppID
	Route           string
	DiegoInstanceID models.InstanceID
}

type route struct {
	Address             string            `json:"address"`
	TLS                 bool              `json:"tls"`
	TTL                 int               `json:"ttl"`
	Tags                map[string]string `json:"tags"`
	PrivateInstanceID   string            `json:"private_instance_id"`
	ServerCertDomainSan string            `json:"server_cert_domain_san"`
}

type GorouterTableParser struct {
}

func NewGorouterTableParser() *GorouterTableParser {
	return &GorouterTableParser{}
}

func (p GorouterTableParser) Parse(path string) (RoutingEntries, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	routes, err := p.decodeFile(f)
	if err != nil {
		return nil, err
	}

	return p.mapRoutes(routes), nil
}

func (p GorouterTableParser) decodeFile(f io.Reader) (map[string][]route, error) {
	result := map[string][]route{}
	dec := json.NewDecoder(f)

	for {
		if err := dec.Decode(&result); err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
	}
}

func (p GorouterTableParser) mapRoutes(routes map[string][]route) RoutingEntries {
	results := RoutingEntries{}

	for route, backends := range routes {
		for _, backend := range backends {
			appId, ok := backend.Tags["app_id"]
			if !ok {
				continue
			}

			processInstanceID, ok := backend.Tags["process_instance_id"]
			if !ok {
				continue
			}

			endpoint := models.ContainerEndpoint(backend.Address)
			results[endpoint] = append(results[endpoint], RoutingEntry{AppID: models.AppID(appId), Route: route, DiegoInstanceID: models.InstanceID(processInstanceID)})
		}
	}

	return results
}
