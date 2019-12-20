package importer

import (
	"fmt"
	"sync"

	"github.com/docker/distribution/registry/client/auth"
)

// NewMultiCredentialStore groups multiple other credentials. Every call to
// Next() moves to the next internal CredentialStore allowing caller to use
// every credential store individually.
func NewMultiCredentialStore(stores ...auth.CredentialStore) (*MultiCredentialStore, error) {
	if len(stores) == 0 {
		return nil, fmt.Errorf("no credential stores")
	}
	return &MultiCredentialStore{
		CredentialStore: stores[0],
		stores:          stores,
	}, nil
}

// MultiCredentialStore groups multiple CredentialStores. By default uses
// the first store and moves to the next when a call to Next() is made. This
// is useful when we have multiple credential sources and we want to try with
// them one at a time.
type MultiCredentialStore struct {
	stores []auth.CredentialStore
	cur    int

	auth.CredentialStore
	sync.Mutex
}

// Next moves to the next credential store. If the current store is the last
// store on the group, move it back to the first.
func (m *MultiCredentialStore) Next() {
	m.Lock()
	defer m.Unlock()

	if m.cur++; m.cur == len(m.stores) {
		m.cur = 0
	}
	m.CredentialStore = m.stores[m.cur]
}

// Len returns the number of internal credential stores.
func (m *MultiCredentialStore) Len() int {
	return len(m.stores)
}

// Err returns the internal error.
// XXX remove this?
func (m *MultiCredentialStore) Err() error {
	return nil
}
