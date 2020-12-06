package gorouter_test

import (
	"io/ioutil"
	"os"

	"github.com/RageZBla/cf-stale-route-detector/gorouter"
	"github.com/RageZBla/cf-stale-route-detector/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GorouterTableParser", func() {
	var (
		parser           *gorouter.GorouterTableParser
		table            string
		routingTableFile *os.File
		err              error
	)

	JustBeforeEach(func() {
		routingTableFile, err = ioutil.TempFile("", "route.json")
		Expect(err).ToNot(HaveOccurred())
		defer routingTableFile.Close()

		_, err = routingTableFile.WriteString(table)
		Expect(err).ToNot(HaveOccurred())
	})

	BeforeEach(func() {
		parser = gorouter.NewGorouterTableParser()
	})

	Describe("Parse", func() {
		BeforeEach(func() {
			table = `
{
	"*.doppler.system.pcf.com":[{"address":"10.223.32.20:8081","tls":true,"ttl":120,"tags":null,"private_instance_id":"8ddb510a-71ae-4fa9-6058-b4db302e2b18","server_cert_domain_san":"doppler.service.cf.internal"},{"address":"10.223.48.26:8081","tls":true,"ttl":120,"tags":null,"private_instance_id":"8f278fa8-e030-4ebd-62c6-17c016b6136e","server_cert_domain_san":"doppler.service.cf.internal"},{"address":"10.223.16.23:8081","tls":true,"ttl":120,"tags":null,"private_instance_id":"1b6ec550-5205-452a-486c-8a5240bbd0be","server_cert_domain_san":"doppler.service.cf.internal"},{"address":"10.223.16.49:8081","tls":true,"ttl":120,"tags":null,"private_instance_id":"39d18efd-4a9b-4543-6603-982352eedcc3","server_cert_domain_san":"doppler.service.cf.internal"},{"address":"10.223.32.33:8081","tls":true,"ttl":120,"tags":null,"private_instance_id":"81b7af64-4643-453b-7688-305976916248","server_cert_domain_san":"doppler.service.cf.internal"},{"address":"10.223.48.16:8081","tls":true,"ttl":120,"tags":null,"private_instance_id":"57cce225-2f97-47f8-5696-58384955f195","server_cert_domain_san":"doppler.service.cf.internal"}],
	"myapp.apps.pcf.com":[{"address":"10.223.32.39:61086","tls":true,"ttl":120,"tags":{"app_id":"bbfd13e9-17b1-4f3e-9c71-a48243de0d20","app_name":"myapp","component":"route-emitter","instance_id":"0","organization_id":"a3c5755c-067d-4737-9749-bc59e2d1160f","organization_name":"myorg","process_id":"bbfd13e9-17b1-4f3e-9c71-a48243de0d20","process_instance_id":"a62d2f03-307f-49c3-66fa-1778","process_type":"web","source_id":"bbfd13e9-17b1-4f3e-9c71-a48243de0d20","space_id":"640c3f05-689a-4d76-b9f1-610b418d96c9","space_name":"myspace"},"private_instance_id":"a62d2f03-307f-49c3-66fa-1778","server_cert_domain_san":"a62d2f03-307f-49c3-66fa-1778"}]
}`
		})

		It("ignores route with missing app_id tag", func() {
			results, err := parser.Parse(routingTableFile.Name())

			Expect(err).ToNot(HaveOccurred())
			Expect(results).ToNot(HaveKey("10.223.32.20:8081"))
		})

		It("maps containerId to app-guids and routes", func() {
			results, err := parser.Parse(routingTableFile.Name())

			Expect(err).ToNot(HaveOccurred())
			endpoint := models.ContainerEndpoint("10.223.32.39:61086")
			Expect(results).To(HaveKey(endpoint))
			Expect(results[endpoint]).To(HaveLen(1))
			Expect(results[endpoint][0]).To(Equal(gorouter.RoutingEntry{AppID: "bbfd13e9-17b1-4f3e-9c71-a48243de0d20", Route: "myapp.apps.pcf.com", DiegoInstanceID: "a62d2f03-307f-49c3-66fa-1778"}))
		})
	})
})
