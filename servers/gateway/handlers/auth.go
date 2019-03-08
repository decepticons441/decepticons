package handlers

//TODO: define HTTP handler functions as described in the
//assignment description. Remember to use your handler context
//struct as the receiver on these functions so that you have
//access to things like the session store and user store.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	// "path"
	"strconv"
	"strings"
	"time"
	"sort"
	"github.com/gorilla/mux"
	"github.com/nehay100/assignments-nehay100/servers/gateway/models/users"
	"github.com/nehay100/assignments-nehay100/servers/gateway/sessions"
)

const AuthorizationHeader = "Authorization"
const ContentTypeHeader = "Content-Type"
const ContentTypeApplicationJSON = "application/json"
const paramAuthorization = "auth"

// func (sh *SessionHandler) UpdateUserInTrie(user *users.User, username bool, remove bool) {
// 	var names [][]string
// 	fname := strings.Split(user.FirstName, " ")
// 	lname := strings.Split(user.LastName, " ")
// 	if username {
// 		useName := strings.Split(user.UserName, " ")
// 		names = append(names, useName)
// 	}	
// 	id := user.ID
// 	names = append(names, fname)
// 	names = append(names, lname)
	
// 	for _, name := range names {
// 		for _ ,word := range name {
// 			if remove {
// 				sh.Trie.Remove(word, id)
// 			} else {
// 				sh.Trie.Add(word, id)
// 			}			
// 		}
// 	}
// }

// func (sh *SessionHandler) UpdateUserInTrie(user *users.User, changeUsername bool, addUser bool) {
// 	var wordSplit []string
// 	if changeUsername {
// 		username := strings.ToLower(user.UserName)
// 		wordSplit = strings.Split(username, " ")
// 	}
// 	log.Println(user.FirstName)
// 	log.Println(user.LastName)
// 	firstname := strings.ToLower(user.FirstName)
// 	lastname := strings.ToLower(user.LastName)
// 	wordSplit = strings.Split(firstname, " ")
// 	wordSplit = strings.Split(lastname, " ")
// 	// if strings.Contains(firstname, " ") {
// 	// 	wordSplit = strings.Split(firstname, " ")
// 	// } else if strings.Contains(lastname, " ") {
// 	// 	wordSplit = strings.Split(lastname, " ")
// 	// }

// 	for _, elem := range wordSplit {
// 		log.Println("elem supposedly added: ", elem)
// 		if addUser {
// 			sh.Trie.Add(elem, user.ID)
// 			currNode := sh.Trie.Root
// 			for _, letter := range elem {
// 				if _, ok := currNode.Children[letter]; ok {
// 					currNode = currNode.Children[letter]
// 					log.Println("added letter into trie: ", string(letter))
// 				} else {
// 					log.Println("didn't add letter into trie: ", string(letter))
// 				}
// 			}
// 		} else {
// 			sh.Trie.Remove(elem, user.ID)
// 		}
// 	}
// }

func (sh *SessionHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if !strings.HasPrefix(r.Header.Get(ContentTypeHeader), ContentTypeApplicationJSON) {
			http.Error(w, "User-Handler: Request body content must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		newUser := &users.NewUser{}
		if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
			http.Error(w, fmt.Sprintf("User-Handler: Error decoding JSON: %v", err),
				http.StatusBadRequest)
			return
		}

		user, err := newUser.ToUser()
		if err != nil {
			http.Error(w, fmt.Sprintf("User-Handler: Error when converting new user to user: %v", err), http.StatusBadRequest)
			return
		}

		insertUser, err := sh.Users.Insert(user)

		if err != nil {
			http.Error(w, fmt.Sprintf("User-Handler: Error when inserting new user to user %v", err), http.StatusBadRequest)
			return
		}

		// upload all the users into trie
		// sh.UpdateUserInTrie(insertUser, true, false)
		sh.Trie.Add(strings.ToLower(insertUser.FirstName), insertUser.ID)
		sh.Trie.Add(strings.ToLower(insertUser.LastName), insertUser.ID)
		sh.Trie.Add(strings.ToLower(insertUser.UserName), insertUser.ID)
		
		// _, err = sh.Users.AddAllToTrie()
		// if err != nil {
		// 	http.Error(w, fmt.Sprintf("User-Handler: Error when inserting new user to trie %v", err), http.StatusInternalServerError)
		// 	return
		// }

		state := &SessionState{
			SessionTime: time.Now(),
			User:        insertUser,
		}

		_, err = sessions.BeginSession(sh.SigningKey, sh.Store, state, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("User-Handler: Error when beginning session: %v", err), http.StatusInternalServerError)
			return
		}
		sh.Trie.Add(strings.ToLower(insertUser.FirstName), insertUser.ID)
		sh.Trie.Add(strings.ToLower(insertUser.LastName), insertUser.ID)
		sh.Trie.Add(strings.ToLower(insertUser.UserName), insertUser.ID)
		

		w.Header().Set(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(insertUser); err != nil {
			http.Error(w, fmt.Sprintf("User-Handler: Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}
	case http.MethodGet:
		var state SessionState
		_, err := sessions.GetState(r, sh.SigningKey, sh.Store, &state)
		if err != nil {
			http.Error(w, fmt.Sprintf("UserHandler: error getting session state/session unauthorized %v", err),
				http.StatusUnauthorized)
			return
		}
		query := r.FormValue("q")
		log.Println("Query: ", query)
		query = strings.ToLower(query)
		if len(query) == 0 {
			http.Error(w, fmt.Sprintf("UserHandler: error parsing id string into int64: %v", err),
				http.StatusBadRequest)
			return
		}
		idSet := sh.Trie.Find(query, 20)
		log.Println("Ids: ", idSet)

		users := []*users.User{}
		for _, id := range idSet {
			user, err := sh.Users.GetByID(id)
			log.Println("User Gotten: ", user)
			if err != nil {
				http.Error(w, fmt.Sprintf("UserHandler: error getting user by id/user not found: %v", err),
					http.StatusNotFound)
				return
			}
			users = append(users, user)
		}
		sort.Slice(users, func(i, j int) bool { 
			return users[i].UserName < users[j].UserName 
		})

		log.Println("Users: ", users)
		w.Header().Set(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, fmt.Sprintf("User-Handler: Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}
	default:
		// Handles any methods not allowed at this resource path
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (sh *SessionHandler) SpecificUserHandler(w http.ResponseWriter, r *http.Request) {
	// idString := path.Base(r.URL.Path)
	idString := mux.Vars(r)["id"]
	sessionState := &SessionState{}
	_, err := sessions.GetState(r, sh.SigningKey, sh.Store, sessionState)
	if err != nil {
		http.Error(w, fmt.Sprintf("SpecificUserHandler: error getting session state/session unauthorized %v", err),
			http.StatusUnauthorized)
		return
	}

	// SHOULD WE HAVE A CHECK IF IDSTRING == SESSIONSTATE.USER.ID?????????
	var userID int64
	if idString == "me" {
		userID = sessionState.User.ID
	} else {
		userID, err = strconv.ParseInt(idString, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("SpecificUserHandler: error parsing id string into int64: %v", err),
				http.StatusBadRequest)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		user, err := sh.Users.GetByID(userID)
		if err != nil {
			http.Error(w, fmt.Sprintf("SpecificUserHandler: error getting user by id/user not found: %v", err),
				http.StatusNotFound)
			return
		}
		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("SpecificUserHandler: Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}
	case http.MethodPatch:
		if userID != sessionState.User.ID {
			http.Error(w, "SpecificUserHandler: Forbidden", http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(r.Header.Get(ContentTypeHeader), ContentTypeApplicationJSON) {
			http.Error(w, "SpecificUserHandler: Request body must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		updates := &users.Updates{}
		if err := json.NewDecoder(r.Body).Decode(updates); err != nil {
			http.Error(w, fmt.Sprintf("SpecificUserHandler: Error decoding JSON: %v", err),
				http.StatusBadRequest)
			return
		}

		user, err := sh.Users.Update(userID, updates)
		if err != nil {
			http.Error(w, fmt.Sprintf("SpecificUserHandler: error updating user: %v", err),
				http.StatusBadRequest)
			return
		}
		sh.Trie.Remove(sessionState.User.FirstName, sessionState.User.ID)
		sh.Trie.Remove(sessionState.User.LastName, sessionState.User.ID)
		sh.Trie.Add(strings.ToLower(user.FirstName), user.ID)
		sh.Trie.Add(strings.ToLower(user.FirstName), user.ID)

		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, fmt.Sprintf("SpecificUserHandler: Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}

	default:
		// Handles any methods not allowed at this resource path
		http.Error(w, "SpecificUserHandler: Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (sh *SessionHandler) SessionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if !strings.HasPrefix(r.Header.Get(ContentTypeHeader), ContentTypeApplicationJSON) {
			http.Error(w, "SessionsHandler: Request body content must be in JSON", http.StatusUnsupportedMediaType)
			return
		}

		state := &SessionState{}

		cred := &users.Credentials{}
		if err := json.NewDecoder(r.Body).Decode(cred); err != nil {
			http.Error(w, fmt.Sprintf("SessionsHandler: Error decoding JSON: %v", err),
				http.StatusBadRequest)
			return
		}

		// HOW TO MAKE THE AUTHENTICATING PROCESS REPLICATE TIME FOR NON-USER PROFILE
		credUser, err := sh.Users.GetByEmail(cred.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("SessionsHandler: Invalid email credentials for given id: %v", err),
				http.StatusBadRequest)
			return
		}

		err = credUser.Authenticate(cred.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("SessionsHandler: Invalid password credentials for given id: %v", err),
				http.StatusUnauthorized)
			return
		}

		ipaddress := r.Header.Get("X-Forwarded-For")
		ipaddress = strings.Split(ipaddress, ", ")[0]
		if len(ipaddress) < 1 {
			ipaddress = r.RemoteAddr
		}
		userSignedIn := &users.UserSignIn{
			ID:        credUser.ID,
			DateTime:  time.Now(),
			IPAddress: ipaddress,
		}
		sh.Users.SignInInsert(userSignedIn)

		state.User = credUser
		state.SessionTime = time.Now()
		_, err = sessions.BeginSession(sh.SigningKey, sh.Store, state, w)
		if err != nil {
			http.Error(w, fmt.Sprintf("SessionsHandler: Error when beginning session: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add(ContentTypeHeader, ContentTypeApplicationJSON)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(credUser); err != nil {
			http.Error(w, fmt.Sprintf("SessionsHandler: Error encoding JSON: %v", err),
				http.StatusInternalServerError)
			return
		}

	default:
		// Handles any methods not allowed at this resource path
		http.Error(w, "SessionsHandler: Method not supported", http.StatusMethodNotAllowed)
		return
	}
}

func (sh *SessionHandler) SpecificSessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		// lastElement := path.Base(r.URL.Path)
		lastElement := mux.Vars(r)["id"]

		if lastElement != "mine" {
			http.Error(w, "SpecificSessionHandler: Status forbidden", http.StatusForbidden)
			return
		}

		state := &SessionState{}
		_, err := sessions.GetState(r, sh.SigningKey, sh.Store, state)
		if err != nil {
			http.Error(w, fmt.Sprintf("SpecificSessionHandler: error getting session state/unauthorized: %v", err),
				http.StatusUnauthorized)
			return
		}

		_, err = sessions.EndSession(r, sh.SigningKey, sh.Store)
		if err != nil {
			http.Error(w, fmt.Sprintf("SpecificSessionHandler: error ending session: %v", err),
				http.StatusInternalServerError)
			return
		}
		w.Header().Add(ContentTypeHeader, "text/plain")
		io.WriteString(w, "signed out")
	default:
		// Handles any methods not allowed at this resource path
		http.Error(w, "SpecificSessionHandler: Method not supported", http.StatusMethodNotAllowed)
		return
	}
}