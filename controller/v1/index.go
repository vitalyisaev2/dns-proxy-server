package v1

import (
	"net/http"

	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/utils"
	. "github.com/mageddo/go-httpmap"
	"golang.org/x/net/context"
)

func init() {
	Get("/", func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Location", "/static")
		res.WriteHeader(301)
	})

	Get("/configuration/", func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		if conf, err := local.LoadConfiguration(); conf != nil {
			utils.GetJsonEncoder(res).Encode(conf)
			return
		} else {
			confLoadError(res, err)
		}
	})
}
