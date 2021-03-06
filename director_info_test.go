package gogobosh_test

import (
	gogobosh "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("get director info", func() {
	It("GET /info to return Director{}", func() {
		request := gogobosh.NewDirectorTestRequest(gogobosh.TestRequest{
			Method: "GET",
			Path:   "/info",
			Response: gogobosh.TestResponse{
				Status: http.StatusOK,
				Body: `{
				  "name": "Bosh Lite Director",
				  "uuid": "bd462a15-213d-448c-aa5b-66624dad3f0e",
				  "version": "1.5.0.pre.1657 (14bc162c)",
				  "user": "admin",
				  "cpi": "warden",
				  "features": {
				    "dns": {
				      "status": false,
				      "extras": {
				        "domain_name": "bosh"
				      }
				    },
				    "compiled_package_cache": {
				      "status": true,
				      "extras": {
				        "provider": "local"
				      }
				    },
				    "snapshots": {
				      "status": false
				    }
				  }
				}`}})
		ts, handler, repo := createDirectorRepo(request)
		defer ts.Close()

		info, apiResponse := repo.GetInfo()
		
		Expect(info.Name                           ).To(Equal("Bosh Lite Director"))
		Expect(info.UUID                           ).To(Equal("bd462a15-213d-448c-aa5b-66624dad3f0e"))
		Expect(info.Version                        ).To(Equal("1.5.0.pre.1657 (14bc162c)"))
		Expect(info.User                           ).To(Equal("admin"))
		Expect(info.CPI                            ).To(Equal("warden"))
		Expect(info.DNSEnabled                     ).To(Equal(false))
		Expect(info.DNSDomainName                  ).To(Equal("bosh"))
		Expect(info.CompiledPackageCacheEnabled    ).To(Equal(true))
		Expect(info.CompiledPackageCacheProvider   ).To(Equal("local"))
		Expect(info.SnapshotsEnabled               ).To(Equal(false))

		Expect(apiResponse.IsSuccessful()).To(Equal(true))
		Expect(handler.AllRequestsCalled()).To(Equal(true))
	})
})

func createDirectorRepo(reqs ...gogobosh.TestRequest) (ts *httptest.Server, handler *gogobosh.TestHandler, repo gogobosh.DirectorRepository) {
	ts, handler = gogobosh.NewTLSServer(reqs)
	config := &gogobosh.Director{
		TargetURL: ts.URL,
		Username:  "admin",
		Password:  "admin",
	}
	gateway := gogobosh.NewDirectorGateway()
	repo = gogobosh.NewBoshDirectorRepository(config, gateway)
	return
}

