// maestro api
// https://github.com/topfreegames/maestro
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/topfreegames/maestro-cli/interfaces"
	"github.com/topfreegames/maestro-cli/metadata"
	"go.uber.org/zap"
)

// App is the api application
type App struct {
	Address    string
	Logger     zap.Logger
	Router     *mux.Router
	Server     *http.Server
	FileSystem interfaces.FileSystem
	Client     interfaces.Client
	Listener   net.Listener
	Context    string
}

func NewApp(
	fs interfaces.FileSystem,
	client interfaces.Client,
	logger zap.Logger,
	context string,
) (*App, error) {
	app := &App{
		Address:    ":57460",
		Logger:     logger,
		FileSystem: fs,
		Client:     client,
		Context:    context,
	}
	err := app.configureApp()
	return app, err
}

func (a *App) configureApp() error {
	a.configureLogger()
	a.configureServer()
	return nil
}

func (a *App) configureLogger() {
	a.Logger = *a.Logger.With(
		zap.String("source",    "maestro-cli"),
		zap.String("operation", "initializeApp"),
		zap.String("version",   metadata.Version),
	)
}

func (a *App) configureServer() {
	a.Router = a.getRouter()
	a.Server = &http.Server{Addr: a.Address, Handler: a.Router}
}

func (a *App) getRouter() *mux.Router {
	r := mux.NewRouter()
	r.Handle("/google-callback",
		NewOAuthCallbackHandler(a, a.FileSystem, a.Client),
	).Methods("GET").Name("oauth2")

	return r
}

//ListenAndLoginAndServe logins and starts local server to get access token from google
func (a *App) ListenAndLoginAndServe() (io.Closer, error) {
	listener, err := net.Listen("tcp", a.Address)
	if err != nil {
		return nil, err
	}
	a.Listener = listener

	if err != nil {
		return nil, err
	}

	err = a.Server.Serve(listener)
	//TODO: do a better check, in case a real "use of closed network connection" happens
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		listener.Close()
		return nil, err
	}

	return listener, nil
}
