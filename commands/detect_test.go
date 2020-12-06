package commands_test

import (
	"errors"

	"github.com/RageZBla/cf-stale-route-detector/commands"
	"github.com/RageZBla/cf-stale-route-detector/commands/fakes"
	"github.com/RageZBla/cf-stale-route-detector/detector"
	fakesPresenter "github.com/RageZBla/cf-stale-route-detector/presenters/fakes"

	"github.com/jessevdk/go-flags"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func executeCommand(command interface{}, args []string) error {
	parser := flags.NewParser(command, flags.HelpFlag|flags.PassDoubleDash)
	_, err := parser.ParseArgs(args)
	Expect(err).NotTo(HaveOccurred())

	commander, ok := command.(flags.Commander)
	Expect(ok).To(BeTrue())

	return commander.Execute(nil)
}

var _ = Describe("Detect", func() {
	var (
		fakeDetector  *fakes.Service
		fakePresenter *fakesPresenter.Presenter
		command       *commands.Detect
	)

	BeforeEach(func() {
		fakeDetector = &fakes.Service{}
		fakePresenter = &fakesPresenter.Presenter{}
		command = commands.NewDetect(fakeDetector, fakePresenter)
	})

	Describe("Execute", func() {
		When("there is no stale route", func() {
			It("returns no error", func() {
				fakeDetector.DetectFromFilesReturns(nil, nil)

				err := executeCommand(command, []string{
					"--routing-table", "table.json",
					"--actual-lrps", "actual.json",
				})
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeDetector.DetectFromFilesCallCount()).To(Equal(1))
				routePath, actualPath, desiredPath := fakeDetector.DetectFromFilesArgsForCall(0)
				Expect(routePath).To(Equal("table.json"))
				Expect(actualPath).To(Equal("actual.json"))
				Expect(desiredPath).To(Equal(""))
			})
		})
		When("there is stale route", func() {
			var (
				staleRoutes map[string][]detector.StaleRoute
			)
			BeforeEach(func() {
				staleRoutes = map[string][]detector.StaleRoute{
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
				fakeDetector.DetectFromFilesReturns(staleRoutes, nil)
			})

			It("returns stale route detected error", func() {
				err := executeCommand(command, []string{
					"--routing-table", "table.json",
					"--actual-lrps", "actual.json",
				})
				Expect(err).To(MatchError("stale route(s) detected"))
			})

			It("prints the result of the detection", func() {

				err := executeCommand(command, []string{
					"--routing-table", "table.json",
					"--actual-lrps", "actual.json",
				})
				Expect(err).To(HaveOccurred())

				Expect(fakeDetector.DetectFromFilesCallCount()).To(Equal(1))
				routePath, actualPath, desiredPath := fakeDetector.DetectFromFilesArgsForCall(0)
				Expect(routePath).To(Equal("table.json"))
				Expect(actualPath).To(Equal("actual.json"))
				Expect(desiredPath).To(Equal(""))

				Expect(fakePresenter.StaleRoutesCallCount()).To(Equal(1))
				actual, verbose := fakePresenter.StaleRoutesArgsForCall(0)
				Expect(actual).To(Equal(staleRoutes))
				Expect(verbose).To(BeFalse())
			})
		})

		When("the service fails to detect stale routes", func() {
			It("returns an error", func() {
				fakeDetector.DetectFromFilesReturns(nil, errors.New("it went wrong"))

				err := executeCommand(command, []string{
					"--routing-table", "table.json",
					"--actual-lrps", "actual.json",
				})
				Expect(err).To(MatchError("it went wrong"))
			})
		})
	})
})
