package acceptance

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("detect stale routes", func() {

	It("detects no stale routes", func() {
		command := exec.Command(pathToMain,
			"detect",
			"--routing-table", "fixtures/routing-table.json",
			"--actual-lrps", "fixtures/actual-lrps.json",
			"--desired-lrps", "fixtures/desired-lrps.json",
		)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(string(session.Out.Contents())).To(Equal("No stale route detected\n"))
	})

	It("detects stale routes", func() {
		command := exec.Command(pathToMain,
			"detect",
			"--routing-table", "fixtures/stale-routing-table.json",
			"--actual-lrps", "fixtures/actual-lrps.json",
			"--desired-lrps", "fixtures/desired-lrps.json",
		)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())

		Eventually(session).Should(gexec.Exit(1))
		Expect(string(session.Err.Contents())).To(Equal("Detected stale route(s)\n"))
	})

	Describe("--verbose flag", func() {
		It("outputs details about detected routes", func() {
			command := exec.Command(pathToMain,
				"detect",
				"--routing-table", "fixtures/stale-routing-table.json",
				"--actual-lrps", "fixtures/actual-lrps.json",
				"--desired-lrps", "fixtures/desired-lrps.json",
				"--verbose",
			)

			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			expectedOutput := []string{
				"Detected stale route(s)",
				"domain: spring-music-wise-grysbok-tb.apps.cloudfoundry.org",
				"app guid: f256797e-a9ac-41bf-a2cd-1b6d37dbe0c4",
				"Diego instance ID: 548a2730-0321-425c-7e6d-96ef",
				"container endpoint: 10.213.60.39:61046",
				"extra:",
				"  app_id: f256797e-a9ac-41bf-a2cd-1b6d37dbe0c4",
				"  organization_id: 49cb11e9-ec12-4d28-9786-dc40be63372c",
				"  routes: spring-music-wise-grysbok-tb.apps.cloudfoundry.org",
				"  process_type: web",
				"  space_id: 6724b327-2e5f-434e-8b03-6d41d66a2279",
				"  app_name: spring-music",
				"  source_id: f256797e-a9ac-41bf-a2cd-1b6d37dbe0c4",
				"  organization_name: brianoc",
				"  space_name: myspace",
				"  instance_id: INDEX",
				"  process_id: f256797e-a9ac-41bf-a2cd-1b6d37dbe0c4",
				"  process_instance_id: INSTANCE_GUID",
			}

			actual := string(session.Err.Contents())
			for _, o := range expectedOutput {
				Expect(actual).To(ContainSubstring(o + "\n"))
			}
		})
	})

})
