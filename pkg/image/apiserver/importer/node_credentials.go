package importer

import (
	"net/url"

	"k8s.io/kubernetes/pkg/credentialprovider"
)

var (
	// NodePullSecretsDir points to the directory from where to read node
	// Docker pull secrets.
	NodePullSecretsDir = "/node/var/lib/kubelet/"
)

// NewNodeCredentialStore returns a credential store holding the content of
// node's Docker pull secrets. If something wrong happens during the object
// initialization an internal error is stored.
func NewNodeCredentialStore() *NodeCredentialStore {
	keyring := &credentialprovider.BasicDockerKeyring{}

	config, err := credentialprovider.ReadDockerConfigJSONFile(
		[]string{NodePullSecretsDir},
	)
	if err == nil {
		keyring.Add(config)
	}

	return &NodeCredentialStore{
		err:     err,
		keyring: keyring,
	}
}

// NodeCredentialStore holds node's Docker pull secrets in an internal
// keyring. It allows callers to query for BasicAuth information by registry
// URL.
type NodeCredentialStore struct {
	keyring credentialprovider.DockerKeyring
	err     error
}

// Basic returns BasicAuth information for the given url. If keyring does not
// have credentials for the url, empty strings are returned.
func (n *NodeCredentialStore) Basic(url *url.URL) (string, string) {
	return basicCredentialsFromKeyring(n.keyring, url)
}

// Err returns NodeCredentialStore's internal error.
func (n *NodeCredentialStore) Err() error {
	return n.err
}
