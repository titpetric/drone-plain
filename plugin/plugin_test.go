// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package plugin

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"

	"github.com/google/go-cmp/cmp"
)

var noContext = context.Background()

func TestPlugin(t *testing.T) {
	req := &secret.Request{
		Path: "secret/docker",
		Name: "username",
		Build: drone.Build{
			Event: "push",
		},
		Repo: drone.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin, _ := New("testdata/secrets.json")
	got, err := plugin.Find(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	want := &drone.Secret{
		Name: "username",
		Data: "david",
		Pull: true,
		Fork: true,
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
		return
	}
}

func TestPlugin_FilterRepo(t *testing.T) {
	req := &secret.Request{
		Path: "secret/docker",
		Name: "username",
		Build: drone.Build{
			Event: "push",
		},
		Repo: drone.Repo{
			Slug: "spaceghost/hello-world",
		},
	}
	plugin, _ := New("testdata/secrets.json")
	_, err := plugin.Find(noContext, req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if want, got := err.Error(), "access denied: repository does not match"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
}

func TestPlugin_FilterEvent(t *testing.T) {
	req := &secret.Request{
		Path: "secret/docker",
		Name: "username",
		Build: drone.Build{
			Event: "pull_request",
		},
		Repo: drone.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin, _ := New("testdata/secrets.json")
	_, err := plugin.Find(noContext, req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if want, got := err.Error(), "access denied: event does not match"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
}

func TestPlugin_NotFound(t *testing.T) {
	req := &secret.Request{
		Path: "secret/docker/registring",
		Name: "username",
		Build: drone.Build{
			Event: "pull_request",
		},
		Repo: drone.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin, _ := New("testdata/secrets.json")
	_, err := plugin.Find(noContext, req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if want, got := err.Error(), "secret not found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
		return
	}
}

func TestPlugin_KeyNotFound(t *testing.T) {
	req := &secret.Request{
		Path: "secret/docker",
		Name: "token",
		Build: drone.Build{
			Event: "push",
		},
		Repo: drone.Repo{
			Slug: "octocat/hello-world",
		},
	}
	plugin, _ := New("testdata/secrets.json")
	_, err := plugin.Find(noContext, req)
	if err == nil {
		t.Errorf("Expect error")
		return
	}
	if got, want := err.Error(), "secret key not found"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
		return
	}
}
