// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"
	"errors"
	"os"

	"encoding/json"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
)

type (
	plugin struct {
		secrets map[string]Secret
	}

	Secret struct {
		Data map[string]interface{} `json:"data,omitempty"`
	}
)

// New returns a new secret plugin that sources secrets
// from the AWS secrets manager.
func New(filename string) (secret.Plugin, error) {
	secrets := map[string]Secret{}

	handle, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	decoder := json.NewDecoder(handle)
	if err := decoder.Decode(&secrets); err != nil {
		return nil, err
	}
	return &plugin{secrets}, nil
}

func (p *plugin) Find(ctx context.Context, req *secret.Request) (*drone.Secret, error) {
	basename := req.Name
	if basename == "" {
		basename = "value"
	}
	dirname := req.Path

	// makes an api call to the aws secrets manager and attempts
	// to retrieve the secret at the requested path.
	params, err := p.find(dirname)
	if err != nil {
		return nil, errors.New("secret not found")
	}
	value, ok := params[basename]
	if !ok {
		return nil, errors.New("secret key not found")
	}

	// the user can filter out requets based on event type
	// using the X-Drone-Events secret key. Check for this
	// user-defined filter logic.
	events := extractEvents(params)
	if !match(req.Build.Event, events) {
		return nil, errors.New("access denied: event does not match")
	}

	// the user can filter out requets based on repository
	// using the X-Drone-Repos secret key. Check for this
	// user-defined filter logic.
	repos := extractRepos(params)
	if !match(req.Repo.Slug, repos) {
		return nil, errors.New("access denied: repository does not match")
	}

	return &drone.Secret{
		Name: basename,
		Data: value,
		Pull: true, // always true. use X-Drone-Events to prevent pull requests.
		Fork: true, // always true. use X-Drone-Events to prevent pull requests.
	}, nil
}

// helper function returns the secret from vault.
func (p *plugin) find(key string) (map[string]string, error) {
	secret, ok := p.secrets[key]
	if !ok || secret.Data == nil {
		return nil, errors.New("secret not found")
	}

	params := map[string]string{}
	for k, v := range secret.Data {
		s, ok := v.(string)
		if !ok {
			continue
		}
		params[k] = s
	}
	return params, nil
}
