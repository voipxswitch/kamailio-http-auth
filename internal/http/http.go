package http

import (
	"encoding/json"
	"net/http"

	"github.com/romana/rlog"
	"goji.io"
	"goji.io/pat"

	"github.com/voipxswitch/kamailio-http-auth/internal/userdata"
)

const (
	requestPath  = "/api/*"
	userEndpoint = "/user/:realm/:username"
)

var (
	h httpHandler
)

type httpHandler struct {}

func New(httpAddress string, root *goji.Mux) {
	v := goji.SubMux()
	root.Handle(pat.New(requestPath), v)
	rlog.Debugf("registered http handler [%s]", requestPath)
	registerMux(v)
	http.ListenAndServe(httpAddress, root)
}

func registerMux(v *goji.Mux) {
	rlog.Debugf("registered user endpoint [%s]", userEndpoint)
	v.HandleFunc(pat.Get(userEndpoint), h.userHandler)
}

// accepts input for devices
func (h *httpHandler) userHandler(w http.ResponseWriter, r *http.Request) {
	realm := pat.Param(r, "realm")
	username := pat.Param(r, "username")
	rlog.Debugf("received http realm [%s] and user [%s] requeset", realm, username)

	type returnData struct {
		Username string `json:"username"`
		Realm    string `json:"realm"`
		Password string `json:"password"`
		Status   string `json:"status"`
	}
	d := returnData{}
	d.Username = username
	d.Realm = realm
	d.Status = "success"
	password, ok := userdata.LoadPassword(username, realm)
	if !ok {
		d.Status = "failure"
	}
	d.Password = password

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
	return
}
