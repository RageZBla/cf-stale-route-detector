package diego_test

import (
	"io/ioutil"
	"os"

	"github.com/RageZBla/cf-stale-route-detector/diego"
	"github.com/RageZBla/cf-stale-route-detector/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DesiredLRPMapper", func() {
	var (
		mapper         *diego.DesiredLRPMapper
		content        string
		desiredLRPFile *os.File
		err            error
	)

	JustBeforeEach(func() {
		desiredLRPFile, err = ioutil.TempFile("", "actual.json")
		Expect(err).ToNot(HaveOccurred())
		defer desiredLRPFile.Close()

		_, err = desiredLRPFile.WriteString(content)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		content = `
{"process_guid":"0079a672-c0c4-46f3-a77f-412f35d0b32b-25022c53-b820-435c-b2d8-cd4f2d9e29cb","routes": { "cf-router": [ { "hostnames": [ "mydomain.apps.pcf.com" ], "port": 8080, "route_service_url": null, "isolation_segment": null }, { "hostnames": [ "myotherdomain.apps.pcf.com" ], "port": 8080, "route_service_url": null, "isolation_segment": null } ], "diego-ssh": { "container_port": 2222, "private_key": "", "host_fingerprint": "" }, "internal-router": [], "tcp-router": [] },"domain":"","rootfs":"","instances":0,"env":null,"start_timeout_ms":0,"disk_mb":0,"memory_mb":0,"cpu_weight":0,"privileged":false,"log_source":"","log_guid":"","metrics_guid":"","annotation":"","max_pids":0,"metric_tags":{"app_name":{"static":"app-1"}}}
{"process_guid":"007ca330-9738-4d06-9769-51954e29d471-af3d65f3-2658-4951-bff9-4ba4f531b4ba","domain":"","rootfs":"","instances":0,"env":null,"start_timeout_ms":0,"disk_mb":0,"memory_mb":0,"cpu_weight":0,"privileged":false,"log_source":"","log_guid":"","metrics_guid":"","annotation":"","max_pids":0,"metric_tags":{"app_name":{"static":"app-2"}}}
{"process_guid":"00b51698-9468-4256-9a3f-0df9a6727f5f-4fd0ce85-2805-4fad-824f-55faf1067af8","domain":"","rootfs":"","instances":0,"env":null,"start_timeout_ms":0,"disk_mb":0,"memory_mb":0,"cpu_weight":0,"privileged":false,"log_source":"","log_guid":"","metrics_guid":"","annotation":"","max_pids":0,"metric_tags":{"app_name":{"static":"app-3"}}}`
		mapper = diego.NewDesiredLRPMapper()
	})

	Describe("Map", func() {
		It("maps app ID to signifiant application properties", func() {
			actual, err := mapper.Map(desiredLRPFile.Name())

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(HaveLen(3))
			Expect(actual).To(HaveKeyWithValue(models.ProcessGuid("0079a672-c0c4-46f3-a77f-412f35d0b32b-25022c53-b820-435c-b2d8-cd4f2d9e29cb"), map[string]string{"routes": "mydomain.apps.pcf.com,myotherdomain.apps.pcf.com", "app_name": "app-1"}))
			Expect(actual).To(HaveKeyWithValue(models.ProcessGuid("007ca330-9738-4d06-9769-51954e29d471-af3d65f3-2658-4951-bff9-4ba4f531b4ba"), map[string]string{"routes": "", "app_name": "app-2"}))
			Expect(actual).To(HaveKeyWithValue(models.ProcessGuid("00b51698-9468-4256-9a3f-0df9a6727f5f-4fd0ce85-2805-4fad-824f-55faf1067af8"), map[string]string{"routes": "", "app_name": "app-3"}))
		})

		When("JSON is invalid", func() {
			BeforeEach(func() {
				content = `/this is not json/`
			})

			It("returns error", func() {
				_, err := mapper.Map(desiredLRPFile.Name())

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("invalid character '/' looking for beginning of value"))
			})
		})

	})
})
