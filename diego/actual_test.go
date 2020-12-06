package diego_test

import (
	"io/ioutil"
	"os"

	"github.com/RageZBla/cf-stale-route-detector/diego"
	"github.com/RageZBla/cf-stale-route-detector/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActualLRPMapper", func() {
	var (
		mapper        *diego.ActualLRPMapper
		content       string
		actualLRPFile *os.File
		err           error
	)

	JustBeforeEach(func() {
		actualLRPFile, err = ioutil.TempFile("", "actual.json")
		Expect(err).ToNot(HaveOccurred())
		defer actualLRPFile.Close()

		_, err = actualLRPFile.WriteString(content)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		content = `
{"process_guid":"0079a672-c0c4-46f3-a77f-412f35d0b32b-25022c53-b820-435c-b2d8-cd4f2d9e29cb","index":0,"domain":"cf-apps","instance_guid":"01c8092c-ce02-469e-65df-0368","cell_id":"9f592063-12b9-4544-a161-6cf389840c35","address":"10.223.32.36","ports":[{"container_port":8080,"host_port":61048,"container_tls_proxy_port":61001,"host_tls_proxy_port":61050},{"container_port":2222,"host_port":61049,"container_tls_proxy_port":61002,"host_tls_proxy_port":61051}],"instance_address":"10.255.57.25","preferred_address":"HOST","crash_count":0,"state":"RUNNING","since":1605773993281003611,"modification_tag":{"epoch":"bf6e09d5-2ba6-4fa5-4269-ee7c6cfdc6f1","index":32},"presence":"ORDINARY"}
{"process_guid":"007ca330-9738-4d06-9769-51954e29d471-af3d65f3-2658-4951-bff9-4ba4f531b4ba","index":0,"domain":"cf-apps","instance_guid":"daadd3a7-d11a-4754-54c2-692a","cell_id":"47bf3c20-70da-4fa8-b9d2-ee984c331711","address":"10.223.32.34","ports":[{"container_port":8080,"host_port":61060,"container_tls_proxy_port":61001,"host_tls_proxy_port":61062},{"container_port":2222,"host_port":61061,"container_tls_proxy_port":61002,"host_tls_proxy_port":61063}],"instance_address":"10.255.83.32","preferred_address":"HOST","crash_count":0,"state":"RUNNING","since":1605773280280989062,"modification_tag":{"epoch":"5b459ad1-3228-4d72-6fcb-9f842d6945d1","index":5},"presence":"ORDINARY"}
{"process_guid":"00b51698-9468-4256-9a3f-0df9a6727f5f-4fd0ce85-2805-4fad-824f-55faf1067af8","index":0,"domain":"cf-apps","instance_guid":"60bd91a3-a27f-404c-736a-48ca","cell_id":"d670b3d6-1d80-454d-bd35-70335ea9b095","address":"10.223.48.19","ports":[{"container_port":8080,"host_port":61044,"container_tls_proxy_port":61001,"host_tls_proxy_port":61046},{"container_port":2222,"host_port":61045,"container_tls_proxy_port":61002,"host_tls_proxy_port":61047}],"instance_address":"10.255.58.33","preferred_address":"HOST","crash_count":0,"state":"RUNNING","since":1605834322039361140,"modification_tag":{"epoch":"035fbd79-9b4c-42f6-5b90-2b6bc97582c8","index":2},"presence":"ORDINARY"}`
		mapper = diego.NewActualLRPMapper()
	})

	Describe("Map", func() {
		It("maps diego cell host TLS endpoint to app GUID", func() {
			actual, err := mapper.Map(actualLRPFile.Name())

			expected := []struct {
				models.InstanceID
				models.ProcessGuid
			}{
				{"01c8092c-ce02-469e-65df-0368", "0079a672-c0c4-46f3-a77f-412f35d0b32b-25022c53-b820-435c-b2d8-cd4f2d9e29cb"},
				{"daadd3a7-d11a-4754-54c2-692a", "007ca330-9738-4d06-9769-51954e29d471-af3d65f3-2658-4951-bff9-4ba4f531b4ba"},
				{"60bd91a3-a27f-404c-736a-48ca", "00b51698-9468-4256-9a3f-0df9a6727f5f-4fd0ce85-2805-4fad-824f-55faf1067af8"},
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(HaveLen(3))
			Expect(actual).To(HaveKeyWithValue(models.ContainerEndpoint("10.223.32.36:61050"), expected[0]))
			Expect(actual).To(HaveKeyWithValue(models.ContainerEndpoint("10.223.32.34:61062"), expected[1]))
			Expect(actual).To(HaveKeyWithValue(models.ContainerEndpoint("10.223.48.19:61046"), expected[2]))
		})

		When("application is not listening to any ports", func() {
			It("ignores the LRPs", func() {

			})
		})

		When("the actual LRPs are not running", func() {
			It("ignores the LRPs", func() {

			})
		})

		When("JSON is invalied", func() {
			BeforeEach(func() {
				content = `/this is not json/`
			})

			It("returns error", func() {
				_, err := mapper.Map(actualLRPFile.Name())

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("invalid character '/' looking for beginning of value"))
			})
		})

	})
})
