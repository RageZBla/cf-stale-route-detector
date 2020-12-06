package presenters_test

import (
	"fmt"

	"github.com/RageZBla/cf-stale-route-detector/detector"
	"github.com/RageZBla/cf-stale-route-detector/logger/fakes"
	"github.com/RageZBla/cf-stale-route-detector/presenters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoggerPresenter", func() {
	var (
		stdoutFakeLogger *fakes.Logger
		stderrFakeLogger *fakes.Logger
		presenter        *presenters.DefaultPresenter
	)

	BeforeEach(func() {
		stdoutFakeLogger = &fakes.Logger{}
		stderrFakeLogger = &fakes.Logger{}
		presenter = presenters.NewLoggerPresenter(stdoutFakeLogger, stderrFakeLogger)
	})

	Describe("StaleRoutes", func() {
		When("there is no stale route", func() {
			It("prints no stale route detected to standard out", func() {
				presenter.StaleRoutes(map[string][]detector.StaleRoute{}, false)

				Expect(stdoutFakeLogger.PrintlnCallCount()).To(Equal(1))
				Expect(stdoutFakeLogger.PrintlnArgsForCall(0)).To(Equal([]interface{}{"No stale route detected"}))
				Expect(stderrFakeLogger.PrintlnCallCount()).To(Equal(0))
			})
		})

		When("there is stale route", func() {
			It("prints 'Dectected stale route(s)' to standard error", func() {
				staleRoutes := map[string][]detector.StaleRoute{
					"foo.domain.com": {
						{
							AppID:             "app-1",
							DiegoInstanceID:   "",
							ContainerEndpoint: "192.168.0.1:61000",
							Extra:             nil,
						},
						{
							AppID:             "app-2",
							DiegoInstanceID:   "other-app",
							ContainerEndpoint: "192.168.0.1:61001",
							Extra:             nil,
						},
					},
				}

				presenter.StaleRoutes(staleRoutes, false)

				Expect(stdoutFakeLogger.PrintlnCallCount()).To(Equal(0))
				Expect(stderrFakeLogger.PrintlnCallCount()).To(Equal(1))
				Expect(stderrFakeLogger.PrintlnArgsForCall(0)).To(Equal([]interface{}{"Detected stale route(s)"}))
			})

			When("verbose mode is activated", func() {

				It("prints stale routes details", func() {
					staleRoutes := map[string][]detector.StaleRoute{
						"foo.domain.com": {
							{
								AppID:             "app-1",
								DiegoInstanceID:   "",
								ContainerEndpoint: "192.168.0.1:61000",
								Extra:             nil,
							},
							{
								AppID:             "app-2",
								DiegoInstanceID:   "other-app",
								ContainerEndpoint: "192.168.0.1:61001",
								Extra:             nil,
							},
						},
					}

					presenter.StaleRoutes(staleRoutes, true)

					Expect(stderrFakeLogger.PrintfCallCount()).To(Equal(3))
					expectedOutput := []string{
						"domain: foo.domain.com\n",
						"app guid: app-1\nDiego instance ID: \ncontainer endpoint: 192.168.0.1:61000\n",
						"app guid: app-2\nDiego instance ID: other-app\ncontainer endpoint: 192.168.0.1:61001\n",
					}
					for i, expected := range expectedOutput {
						format, content := stderrFakeLogger.PrintfArgsForCall(i)
						Expect(fmt.Sprintf(format, content...)).To(Equal(expected))
					}
				})
				When("data contains extra information", func() {
					It("outputs the extra information with indentation", func() {

						staleRoutes := map[string][]detector.StaleRoute{
							"bar.domain.com": {
								{
									AppID:             "app-1",
									DiegoInstanceID:   "",
									ContainerEndpoint: "192.168.0.1:61000",
									Extra: map[string]string{
										"foo": "extra",
									},
								},
							},
						}

						presenter.StaleRoutes(staleRoutes, true)

						Expect(stderrFakeLogger.PrintfCallCount()).To(Equal(4))
						expectedOutput := []string{
							"domain: bar.domain.com\n",
							"app guid: app-1\nDiego instance ID: \ncontainer endpoint: 192.168.0.1:61000\n",
							"extra:\n",
							"  foo: extra\n",
						}
						for i, expected := range expectedOutput {
							format, content := stderrFakeLogger.PrintfArgsForCall(i)
							Expect(fmt.Sprintf(format, content...)).To(Equal(expected))
						}
					})
				})
			})
		})
	})
})
