package handlers

//TODO: define a handler context struct that
//will be a receiver on any of your HTTP
//handler functions that need access to
//globals, such as the key used for signing
//and verifying SessionIDs, the session store
//and the user store

import (
	"github.com/nehay100/assignments-nehay100/servers/gateway/indexes"
	"github.com/nehay100/assignments-nehay100/servers/gateway/models/users"
	"github.com/nehay100/assignments-nehay100/servers/gateway/sessions"
)

// MyHandler is a struct that has the signing key, getting/saving session state
// and finding/saving user profiles
type SessionHandler struct {
	SigningKey string
	Store      sessions.Store
	Users      users.Store
	Trie       indexes.Trie
}
