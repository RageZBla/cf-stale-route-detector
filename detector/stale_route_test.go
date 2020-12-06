package detector_test

import (
	"errors"

	"github.com/RageZBla/cf-stale-route-detector/detector"
	"github.com/RageZBla/cf-stale-route-detector/detector/fakes"
	"github.com/RageZBla/cf-stale-route-detector/diego"
	"github.com/RageZBla/cf-stale-route-detector/gorouter"
	"github.com/RageZBla/cf-stale-route-detector/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StaleRouteDetector", func() {
	var (
		fakeParser        *fakes.RoutingTableParser
		fakeActualMapper  *fakes.ActualLRPMapper
		fakeDesiredMapper *fakes.DesiredLRPMapper
		service           *detector.StaleRouteDetector
	)

	BeforeEach(func() {
		fakeParser = &fakes.RoutingTableParser{}
		fakeActualMapper = &fakes.ActualLRPMapper{}
		fakeDesiredMapper = &fakes.DesiredLRPMapper{}
		service = detector.NewStaleRouteDetector(fakeParser, fakeActualMapper, fakeDesiredMapper)
	})

	Describe("Detect", func() {
		When("There is no routes", func() {
			It("does not detect any stale route", func() {
				routingTable := gorouter.RoutingEntries{}
				actualLRPMapping := diego.ActualLRPMapping{"192.168.0.1:61000": struct {
					models.InstanceID
					models.ProcessGuid
				}{"instance-id", "process-guid"}}
				desiredLRPMapping := diego.DesiredLRPMapping{}

				staleRoutes := service.Detect(routingTable, actualLRPMapping, desiredLRPMapping)

				Expect(staleRoutes).To(HaveLen(0))
			})
		})

		When("all the routes matches", func() {
			It("does not detect any stale route", func() {
				routingTable := gorouter.RoutingEntries{"192.168.0.1:61000": {{AppID: "app-guid", Route: "foo.domain.com", DiegoInstanceID: "instance-id"}}}
				actualLRPMapping := diego.ActualLRPMapping{"192.168.0.1:61000": struct {
					models.InstanceID
					models.ProcessGuid
				}{"instance-id", "process-guid"}}
				desiredLRPMapping := diego.DesiredLRPMapping{}

				staleRoutes := service.Detect(routingTable, actualLRPMapping, desiredLRPMapping)

				Expect(staleRoutes).To(HaveLen(0))
			})
		})

		When("there is no record of such app on Diego actual LRPS", func() {
			It("does not detect any stale route", func() {
				routingTable := gorouter.RoutingEntries{"192.168.0.1:61000": {{AppID: "app-guid", Route: "foo.domain.com", DiegoInstanceID: "instance-id"}}}
				actualLRPMapping := diego.ActualLRPMapping{}
				desiredLRPMapping := diego.DesiredLRPMapping{}

				staleRoutes := service.Detect(routingTable, actualLRPMapping, desiredLRPMapping)

				Expect(staleRoutes).To(HaveLen(0))
			})
		})

		When("the diego and router instance ID does not match", func() {
			It("detects and return stale route", func() {
				routingTable := gorouter.RoutingEntries{"192.168.0.1:61000": {{AppID: "app-guid", Route: "foo.domain.com", DiegoInstanceID: "instance-id"}}}
				actualLRPMapping := diego.ActualLRPMapping{"192.168.0.1:61000": struct {
					models.InstanceID
					models.ProcessGuid
				}{"other-instance-id", "process-guid"}}
				desiredLRPMapping := diego.DesiredLRPMapping{}

				staleRoutes := service.Detect(routingTable, actualLRPMapping, desiredLRPMapping)

				expectedStaleRoutes := map[string][]detector.StaleRoute{"foo.domain.com": {{AppID: "app-guid", ContainerEndpoint: "192.168.0.1:61000", DiegoInstanceID: "other-instance-id", Extra: nil}}}
				Expect(staleRoutes).To(Equal(expectedStaleRoutes))
			})
		})

		When("desired LRPS are being provided", func() {
			It("adds the desired information to the stale route", func() {
				routingTable := gorouter.RoutingEntries{"192.168.0.1:61000": {{AppID: "app-guid", Route: "foo.domain.com"}}}
				actualLRPMapping := diego.ActualLRPMapping{"192.168.0.1:61000": struct {
					models.InstanceID
					models.ProcessGuid
				}{"other-instance-id", "process-guid"}}
				desiredLRPMapping := diego.DesiredLRPMapping{"process-guid": {"foo": "extra"}}

				staleRoutes := service.Detect(routingTable, actualLRPMapping, desiredLRPMapping)

				expectedStaleRoutes := map[string][]detector.StaleRoute{"foo.domain.com": {{AppID: "app-guid", ContainerEndpoint: "192.168.0.1:61000", DiegoInstanceID: "other-instance-id", Extra: map[string]string{"foo": "extra"}}}}
				Expect(staleRoutes).To(Equal(expectedStaleRoutes))
			})
		})
	})

	Describe("DetectFromFiles", func() {
		When("route parser can not parse the routing table provided", func() {
			It("returns an error", func() {
				fakeParser.ParseReturns(gorouter.RoutingEntries{}, errors.New("this is fine"))

				_, err := service.DetectFromFiles("table.json", "actual.json", "")
				Expect(err).To(MatchError("could not parse routing table: this is fine"))
			})
		})

		When("actual LRP mapper can not map the file provided", func() {
			It("returns an error", func() {
				fakeActualMapper.MapReturns(nil, errors.New("this is fine"))

				_, err := service.DetectFromFiles("table.json", "actual.json", "")
				Expect(err).To(MatchError("could not parse actual LRPs: this is fine"))
			})
		})

		When("desired LRP mapper can not map the file provided", func() {
			It("returns an error", func() {

				fakeDesiredMapper.MapReturns(nil, errors.New("this is fine"))

				_, err := service.DetectFromFiles("table.json", "actual.json", "desired.json")

				Expect(fakeDesiredMapper.MapCallCount()).To(Equal(1))
				Expect(fakeDesiredMapper.MapArgsForCall(0)).To(Equal("desired.json"))
				Expect(err).To(MatchError("could not parse desired LRPs: this is fine"))
			})
		})

		It("detects stale routes", func() {

			routingTable := gorouter.RoutingEntries{"192.168.0.1:61000": {{AppID: "app-guid", Route: "foo.domain.com", DiegoInstanceID: "instance-id"}}}
			actualLRPMapping := diego.ActualLRPMapping{"192.168.0.1:61000": struct {
				models.InstanceID
				models.ProcessGuid
			}{"other-instance-id", "process-guid"}}
			fakeParser.ParseReturns(routingTable, nil)
			fakeActualMapper.MapReturns(actualLRPMapping, nil)
			fakeDesiredMapper.MapReturns(diego.DesiredLRPMapping{"process-guid": {"extra": "stuff"}}, nil)

			staleRoutes, err := service.DetectFromFiles("table.json", "actual.json", "desired.json")

			Expect(err).ToNot(HaveOccurred())
			expectedStaleRoutes := map[string][]detector.StaleRoute{"foo.domain.com": {{AppID: "app-guid", ContainerEndpoint: "192.168.0.1:61000", DiegoInstanceID: "other-instance-id", Extra: map[string]string{"extra": "stuff"}}}}
			Expect(staleRoutes).To(Equal(expectedStaleRoutes))
		})
	})

})
