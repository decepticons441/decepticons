package handlers

//TODO: define a session state struct for this web server
//see the assignment description for the fields you should include
//remember that other packages can only see exported fields!

import (
	"time"

	"github.com/nehay100/assignments-nehay100/servers/gateway/models/users"
)

// SessionState struct contains time at which session began and authenticated user
// who started the session
type SessionState struct {
	SessionTime time.Time
	User        *users.User
}
