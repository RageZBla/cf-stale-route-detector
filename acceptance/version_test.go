package acceptance

import (
	"fmt"
	"time"

	"github.com/RageZBla/cf-stale-route-detector/cmd"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("version command", func() {
	var version string
	var buffer *gbytes.Buffer

	BeforeEach(func() {
		version = fmt.Sprintf("v0.0.0-dev.%d", time.Now().Unix())
		buffer = gbytes.NewBuffer()
	})

	When("given the -v short flag", func() {
		It("returns the binary version", func() {
			err := cmd.Main(buffer, buffer, version, []string{"cf-stale-route-detector", "-v"})
			Expect(err).ToNot(HaveOccurred())

			Expect(buffer).To(gbytes.Say(fmt.Sprintf("%s\n", version)))
		})
	})

	When("given the --version long flag", func() {
		It("returns the binary version", func() {
			err := cmd.Main(buffer, buffer, version, []string{"cf-stale-route-detector", "--version"})
			Expect(err).ToNot(HaveOccurred())

			Expect(buffer).To(gbytes.Say(fmt.Sprintf("%s\n", version)))
		})
	})

	When("given the version command", func() {
		It("returns the binary version", func() {
			err := cmd.Main(buffer, buffer, version, []string{"cf-stale-route-detector", "version"})
			Expect(err).ToNot(HaveOccurred())

			Expect(buffer).To(gbytes.Say(fmt.Sprintf("%s\n", version)))
		})
	})
})
