package importer

import (
	"net/url"

	"github.com/openshift/library-go/pkg/image/registryclient"
	"k8s.io/kubernetes/pkg/credentialprovider"
)

var (
	// NodePullSecretsDir points to the directory from where we read node pullsecrets.
	NodePullSecretsDir = "/node/var/lib/kubelet/"
)

// NewNodeCredentialStore returns a credential store holding content of node's pullsecrets.
func NewNodeCredentialStore() (*NodeCredentialStore, error) {
	config, err := credentialprovider.ReadDockerConfigJSONFile(
		[]string{NodePullSecretsDir},
	)
	if err != nil {
		return nil, err
	}

	keyring := &credentialprovider.BasicDockerKeyring{}
	keyring.Add(config)

	return &NodeCredentialStore{
		keyring:           keyring,
		RefreshTokenStore: registryclient.NewRefreshTokenStore(),
	}, nil
}

// NodeCredentialStore holds node's pull secrets in a keyring.
type NodeCredentialStore struct {
	keyring credentialprovider.DockerKeyring
	registryclient.RefreshTokenStore
}

// Basic returns basic authentication for url. If keyring does not have credentials for the url,
// empty strings are returned.
func (n *NodeCredentialStore) Basic(url *url.URL) (string, string) {
	return basicCredentialsFromKeyring(n.keyring, url)
}

// Err returns credential store internal error.
func (n *NodeCredentialStore) Err() error {
	return nil
}
