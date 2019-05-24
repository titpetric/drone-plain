// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"net/http"

	"github.com/drone/drone-go/plugin/secret"
	"github.com/titpetric/drone-plain/plugin"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	Address string `envconfig:"SERVER_ADDRESS"`
	Secret  string `envconfig:"SECRET_KEY"`
	Debug   bool   `envconfig:"DEBUG"`
	Source  string `envconfig:"SOURCE"`
}

func main() {
	spec := new(config)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	service, err := plugin.New(spec.Source)
	if err != nil {
		logrus.Fatal(err)
	}

	http.Handle("/", secret.Handler(
		spec.Secret,
		service,
		logrus.StandardLogger(),
	))

	logrus.Infof("server listening on address %s", spec.Address)
	if err := http.ListenAndServe(spec.Address, nil); err != nil {
		logrus.Fatal(err)
	}
}
