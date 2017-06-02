// maestro-cli api
// https://github.com/topfreegames/maestro-cli
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/topfreegames/maestro-cli/interfaces"
	"github.com/topfreegames/maestro-cli/login"
)

const Index = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Maestro</title>
</head>
<body>
  <h1>Thanks for logging in</h1>
  You can go back to your terminal
</body>
</html>
`

const UnauthorizedIndex = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Maestro</title>
</head>
<body>
  <h1>Unauthorized</h1>
  Your email is not authorized to use Maestro
</body>
</html>
`

//OAuthCallbackHandler handles the callback after user approves/deny auth
type OAuthCallbackHandler struct {
	app    *App
	fs     interfaces.FileSystem
	client interfaces.Client
}

func NewOAuthCallbackHandler(
	app *App,
	fs interfaces.FileSystem,
	client interfaces.Client,
) *OAuthCallbackHandler {
	return &OAuthCallbackHandler{
		app:    app,
		fs:     fs,
		client: client,
	}
}

//ServeHTTP method
func (o *OAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	l := o.app.Logger
	l.Debugf("Returned state %s and code %s", state, code)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := login.SaveAccessToken(
		state, code, o.app.Login.OAuthState, o.app.Login.ServerURL,
		o.fs,
		o.client,
	)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, UnauthorizedIndex)
		o.app.Listener.Close()
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, Index)
	o.app.Listener.Close()
}
