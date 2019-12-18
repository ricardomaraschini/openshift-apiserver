package importer

import (
	"fmt"
	"net/url"
	"os"
	"testing"
)

func TestNewNodeCredentialStore(t *testing.T) {
	if _, err := NewNodeCredentialStore(); err == nil {
		t.Error("able to create credential store with invalid path")
	}
}

func TestBasic(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	oldDir := NodePullSecretsDir
	NodePullSecretsDir = fmt.Sprintf("%s/test/", dir)
	store, err := NewNodeCredentialStore()
	if err != nil {
		t.Fatalf("unexpected error creating credential store: %v", err)
	}
	NodePullSecretsDir = oldDir

	for _, tt := range []struct {
		name string
		url  *url.URL
		user string
		pass string
	}{
		{
			name: "valid registry",
			url:  &url.URL{Host: "registry0.redhat.io"},
			user: "registry0",
			pass: "registry0",
		},
		{
			name: "invalid registry",
			url:  &url.URL{Host: "invalidregistry.redhat.io"},
		},
		{
			name: "nil url",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			user, pass := store.Basic(tt.url)
			if user != tt.user {
				t.Errorf("invalid user for %s", user)
			}
			if pass != tt.pass {
				t.Errorf("invalid user for %s", user)
			}
		})
	}
}
