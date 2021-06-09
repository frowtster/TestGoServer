package router

import (
	"fmt"
	"my-simple-server/t_util"
	"net/http"
	"net/http/httputil"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

func init() {
	Router.GET("/test/:name", Test)
}

// Test func
func Test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if t_util.Log.Level == logrus.DebugLevel {
		dump, _ := httputil.DumpRequest(r, true)
		t_util.Log.Debugf("%q", dump)
	}

	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
