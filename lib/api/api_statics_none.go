// +build !msyncdebug

package api

import "net/http"

type staticsServer struct {
}

func newStaticsServer(theme, assetDir string) *staticsServer {
	return &staticsServer{}
}

func (s *staticsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
