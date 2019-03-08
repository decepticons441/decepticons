package sessions

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const headerAuthorization = "Authorization"
const paramAuthorization = "auth"
const schemeBearer = "Bearer "

//ErrNoSessionID is used when no session ID was found in the Authorization header
var ErrNoSessionID = errors.New("no session ID found in " + headerAuthorization + " header")

//ErrInvalidScheme is used when the authorization scheme is not supported
var ErrInvalidScheme = errors.New("authorization scheme not supported")

//BeginSession creates a new SessionID, saves the `sessionState` to the store, adds an
//Authorization header to the response with the SessionID, and returns the new SessionID
func BeginSession(signingKey string, store Store, sessionState interface{}, w http.ResponseWriter) (SessionID, error) {
	// Create a new SessionID
	sessionID, err := NewSessionID(signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error when creating signing key: %s", err)
	}

	// Save the sessionState to the store
	log.Println("BeginSession: " + sessionID)
	// store.Save(sessionID, sessionState)
	err = store.Save(sessionID, sessionState)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error when saving signing key: %s", err)
	}

	// Add a Header to the ResponseWriter that looks like this
	//    "Authorization: Bearer <sessionID>"
	// where "<sessionID>" is replaced with the newly-created SessionID
	// (note the constants declared for you above, which will help you avoid typos)

	w.Header().Add(headerAuthorization, schemeBearer+sessionID.String())

	return sessionID, nil
}

//GetSessionID extracts and validates the SessionID from the request headers
func GetSessionID(r *http.Request, signingKey string) (SessionID, error) {
	//TODO: get the value of the Authorization header,
	//or the "auth" query string parameter if no Authorization header is present,
	//and validate it. If it's valid, return the SessionID. If not
	//return the validation error.

	sessionID := r.Header.Get(headerAuthorization)
	if len(sessionID) == 0 {
		sessionID = r.FormValue(paramAuthorization)
		// if len(sessionID) == 0 {
		// 	return InvalidSessionID, ErrNoSessionID
		// }
	}
	if len(sessionID) == 0 {
		// sessionID = r.FormValue(paramAuthorization)
		// if len(sessionID) == 0 {
		return InvalidSessionID, ErrNoSessionID
		// }
	}
	// if !strings.Contains(sessionID, schemeBearer) {
	// 	return InvalidSessionID, fmt.Errorf("error when validating session id from header")
	// }

	sessionID = strings.Replace(sessionID, schemeBearer, "", 1)
	// sessionID = strings.TrimPrefix(sessionID, schemeBearer)

	id, err := ValidateID(sessionID, signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error when validating session id from header: %v", err)
	}
	return id, nil
}

//GetState extracts the SessionID from the request,
//gets the associated state from the provided store into
//the `sessionState` parameter, and returns the SessionID
func GetState(r *http.Request, signingKey string, store Store, sessionState interface{}) (SessionID, error) {
	//TODO: get the SessionID from the request, and get the data
	//associated with that SessionID from the store.
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, ErrStateNotFound
	}

	err = store.Get(sessionID, sessionState)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error when getting session: %v", err)
	}

	return sessionID, nil
}

//EndSession extracts the SessionID from the request,
//and deletes the associated data in the provided store, returning
//the extracted SessionID.
func EndSession(r *http.Request, signingKey string, store Store) (SessionID, error) {
	//TODO: get the SessionID from the request, and delete the
	//data associated with it in the store.
	sessionID, err := GetSessionID(r, signingKey)
	if err != nil {
		return InvalidSessionID, fmt.Errorf("error when getting a valid session id: %v", err)
	}

	err = store.Delete(SessionID(sessionID))

	if err != nil {
		return InvalidSessionID, fmt.Errorf("error when getting session: %v", err)
	}
	r.Header.Del(headerAuthorization)
	r.URL.Query().Del(paramAuthorization)
	return sessionID, nil
}
