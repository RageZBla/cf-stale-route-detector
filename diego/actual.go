package diego

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	bbsModels "code.cloudfoundry.org/bbs/models"

	"github.com/RageZBla/cf-stale-route-detector/models"
)

type ActualLRPMapping map[models.ContainerEndpoint]struct {
	models.InstanceID
	models.ProcessGuid
}

type ActualLRPMapper struct {
}

const appIdLen = 36             // 0079a672-c0c4-46f3-a77f-412f35d0b32b
const processGuidLen = 1 + 36*2 // 0079a672-c0c4-46f3-a77f-412f35d0b32b-25022c53-b820-435c-b2d8-cd4f2d9e29cb

func NewActualLRPMapper() *ActualLRPMapper {
	return &ActualLRPMapper{}
}

func (m ActualLRPMapper) Map(path string) (ActualLRPMapping, error) {
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

func (m ActualLRPMapper) decodeFile(f io.Reader) ([]bbsModels.ActualLRP, error) {
	result := []bbsModels.ActualLRP{}
	dec := json.NewDecoder(f)

	for {
		var actualLRP bbsModels.ActualLRP
		if err := dec.Decode(&actualLRP); err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
		result = append(result, actualLRP)
	}
}

func (m ActualLRPMapper) mapLRPs(lrps []bbsModels.ActualLRP) (results ActualLRPMapping, err error) {
	var (
		// appId             gorouter.AppId
		instanceID models.InstanceID
		endpoint   models.ContainerEndpoint
	)
	results = ActualLRPMapping{}

	for _, lrp := range lrps {
		// appId, err = extractAppId(lrp.ActualLRPKey.ProcessGuid)
		// if err != nil {
		// 	return
		// }

		instanceID = models.InstanceID(lrp.InstanceGuid)

		endpoint, err = m.extractEndpoint(lrp)
		if err != nil {
			return
		}

		results[endpoint] = struct {
			models.InstanceID
			models.ProcessGuid
		}{instanceID, models.ProcessGuid(lrp.ProcessGuid)}
	}

	return results, nil
}

func (m ActualLRPMapper) extractEndpoint(lrp bbsModels.ActualLRP) (models.ContainerEndpoint, error) {
	cellAddress := lrp.Address
	if len(lrp.Ports) == 0 {
		return "", nil
	}
	containerPort := lrp.ActualLRPNetInfo.Ports[0].HostTlsProxyPort

	return models.ContainerEndpoint(fmt.Sprintf("%s:%d", cellAddress, containerPort)), nil
}
