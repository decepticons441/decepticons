package handlers

import (
	"bytes"
	"encoding/json"

	// "fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/nehay100/decepticons/servers/gateway/indexes"
	"github.com/nehay100/decepticons/servers/gateway/models/users"
	"github.com/nehay100/decepticons/servers/gateway/sessions"
)

func makeSessionHandler(expectedErr bool) *SessionHandler {
	sh := &SessionHandler{
		SigningKey: "testKey",
		Store:      sessions.NewMemStore(time.Hour, time.Minute),
		Users:      users.NewFakeMockStore(expectedErr),
		Trie:       *indexes.NewTrie(),
	}
	return sh
}

func makeSameUser() *users.User {
	user := &users.User{
		ID:        1,
		Email:     "nehay100@gmail.com",
		UserName:  "nehay100",
		FirstName: "Neha",
		LastName:  "Yadav",
		PhotoURL:  "",
	}
	return user
}

func TestUsersHandler(t *testing.T) {
	//	ADD MOCK STORE INSERT CASE IF COVERAGE IS LOW!!!!!!!!
	cases := []struct {
		name                string
		request             string
		jsonString          string
		expectedErr         bool
		expectedStatusCode  int
		expectedContentType string
		expectedReturn      *users.User
	}{
		{
			name:    "Valid User",
			request: http.MethodPost,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure",
				"passwordConf": "thisPasswordIsSecure",
				"userName":     "nehay100",
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			expectedErr:         false,
			expectedStatusCode:  http.StatusCreated,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
		{
			name:    "Non Valid Content User - PNG",
			request: http.MethodPost,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure",
				"passwordConf": "thisPasswordIsSecure",
				"userName":     "nehay100",
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			expectedErr:         true,
			expectedStatusCode:  http.StatusUnsupportedMediaType,
			expectedContentType: "image/png",
			expectedReturn:      makeSameUser(),
		},
		{
			name:    "Non Valid HTTP Method User",
			request: http.MethodPatch,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure",
				"passwordConf": "thisPasswordIsSecure",
				"userName":     "nehay100",
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			expectedErr:         true,
			expectedStatusCode:  http.StatusMethodNotAllowed,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
		// {
		// 	name:    "Non Valid HTTP Method User",
		// 	request: http.MethodGet,
		// 	jsonString: `{"email": "nehay100@gmail.com",
		// 		"password":     "thisPasswordIsSecure",
		// 		"passwordConf": "thisPasswordIsSecure",
		// 		"userName":     "nehay100",
		// 		"firstName":    "Neha",
		// 		"lastName":     "Yadav"}`,
		// 	expectedErr:         true,
		// 	expectedStatusCode:  http.StatusMethodNotAllowed,
		// 	expectedContentType: ContentTypeApplicationJSON,
		// 	expectedReturn:      makeSameUser(),
		// },
		{
			name:    "Non Valid HTTP Method User",
			request: http.MethodDelete,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure",
				"passwordConf": "thisPasswordIsSecure",
				"userName":     "nehay100",
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			expectedErr:         true,
			expectedStatusCode:  http.StatusMethodNotAllowed,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
		{
			name:    "Non Valid HTTP Method User",
			request: http.MethodOptions,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure",
				"passwordConf": "thisPasswordIsSecure",
				"userName":     "nehay100",
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			expectedErr:         true,
			expectedStatusCode:  http.StatusMethodNotAllowed,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
		{
			name:                "Non Valid JSON Body String",
			request:             http.MethodPost,
			jsonString:          "",
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
		{
			name:    "Different Mock Store and JSON Body New User",
			request: http.MethodPost,
			jsonString: `{"email": "nehay100@uw.edu",
				"password":     "thisPasswordIsSecure2.0",
				"passwordConf": "thisPasswordIsSecure2.0",
				"userName":     "nehay1002.0",
				"firstName":    "Neha2.0",
				"lastName":     "Yadav2.0"}`,
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
		{
			name:    "Invalid new user to test ToUser method",
			request: http.MethodPost,
			jsonString: `{"email": "nehay100@uw.edu",
				"password":     "thisPasswordIsSecure2.0",
				"passwordConf": "thisPasswordIsSecure",
				"userName":     "nehay1002.0",
				"firstName":    "Neha2.0",
				"lastName":     "Yadav2.0"}`,
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
		},
	}

	for _, c := range cases {
		sh := makeSessionHandler(c.expectedErr)

		jsonByte := []byte(c.jsonString)
		jsonIO := bytes.NewBuffer(jsonByte)
		request := httptest.NewRequest(c.request, "/v1/users/", jsonIO)
		request.Header.Add(ContentTypeHeader, c.expectedContentType)
		recorder := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/v1/users/", sh.UsersHandler)
		router.ServeHTTP(recorder, request)

		// sh.UsersHandler(recorder, request)

		response := recorder.Result()
		//check the response status code
		if response.StatusCode != c.expectedStatusCode {
			t.Errorf("case %s: incorrect status code: expected %d but got %d",
				c.name, c.expectedStatusCode, response.StatusCode)
		}
		user := &users.User{}
		err := json.NewDecoder(response.Body).Decode(user)
		if c.expectedErr && err == nil {
			t.Errorf("case %s: expected error but revieved none", c.name)
		}

		if !c.expectedErr && c.expectedReturn.Email != user.Email && !reflect.DeepEqual(c.expectedReturn.PassHash, user.PassHash) &&
			c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
			t.Errorf("case %s: incorrect return: expected %v but revieved %v",
				c.name, c.expectedReturn, user)
		}
	}
}
func TestSpecificUserHandler(t *testing.T) {
	//	ADD BAD SIGNING KEY CASE IF COVERAGE IS LOW!!!!!!!!!!

	cases := []struct {
		name                string
		request             string
		jsonString          string
		authRequired        bool
		useMe               bool
		expectedErr         bool
		wrongSignKey        bool
		expectedStatusCode  int
		expectedContentType string
		expectedReturn      *users.User
		updates             *users.Updates
	}{
		{
			name:                "Valid Get User w/ id",
			request:             http.MethodGet,
			jsonString:          "",
			authRequired:        true,
			useMe:               false,
			expectedErr:         false,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusOK,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:    "Valid Patch User w/ id",
			request: http.MethodPatch,
			jsonString: `{
				"firstName":    "Diana",
				"lastName":     "Prince"}`,
			authRequired:        true,
			useMe:               false,
			expectedErr:         false,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusOK,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "Diana",
				LastName:  "Prince",
			},
		},
		{
			name:                "Valid Get User w/ me",
			request:             http.MethodGet,
			jsonString:          "",
			authRequired:        true,
			useMe:               true,
			expectedErr:         false,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusOK,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:    "Valid Patch User w/ me",
			request: http.MethodPatch,
			jsonString: `{
				"firstName":    "Diana",
				"lastName":     "Prince"}`,
			authRequired:        true,
			useMe:               true,
			expectedErr:         false,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusOK,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "Diana",
				LastName:  "Prince",
			},
		},
		{
			name:                "Invalid ParseInt Test",
			request:             http.MethodGet,
			jsonString:          "",
			authRequired:        true,
			useMe:               false,
			expectedErr:         false,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:                "Invalid SigningKey",
			request:             http.MethodPost,
			jsonString:          "",
			authRequired:        true,
			useMe:               false,
			expectedErr:         true,
			wrongSignKey:        true,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:                "Invalid Authorization w/ me",
			request:             http.MethodGet,
			jsonString:          "",
			authRequired:        false,
			useMe:               true,
			expectedErr:         true,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:                "Invalid Authorization w/o me",
			request:             http.MethodGet,
			jsonString:          "",
			authRequired:        false,
			useMe:               false,
			expectedErr:         true,
			wrongSignKey:        false,
			expectedStatusCode:  http.StatusUnauthorized,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:                "Invalid Signing Key",
			request:             http.MethodGet,
			jsonString:          "",
			authRequired:        true,
			useMe:               false,
			expectedErr:         true,
			wrongSignKey:        true,
			expectedStatusCode:  http.StatusNotFound,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:    "Invalid Content User - PNG for Patch",
			request: http.MethodPatch,
			jsonString: `{
				"firstName":    "Steven",
				"lastName":     "Rogers"}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusUnsupportedMediaType,
			expectedContentType: "image/png",
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "Steven",
				LastName:  "Rogers",
			},
		},
		{
			name:    "Invalid HTTP Method User",
			request: http.MethodPost,
			jsonString: `{
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusMethodNotAllowed,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:    "Invalid HTTP Method User",
			request: http.MethodDelete,
			jsonString: `{
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusMethodNotAllowed,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:    "Invalid HTTP Method User",
			request: http.MethodOptions,
			jsonString: `{
				"firstName":    "Neha",
				"lastName":     "Yadav"}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusMethodNotAllowed,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:                "Invalid JSON Body String - Patch",
			request:             http.MethodPatch,
			jsonString:          "",
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates:             &users.Updates{},
		},
		{
			name:    "Invalid GetById",
			request: http.MethodGet,
			jsonString: `{
				"firstName":    "Tony",
				"lastName":     "Stark"}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusNotFound,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "Tony",
				LastName:  "Stark",
			},
		},
		{
			name:    "Invalid Update Object 1",
			request: http.MethodPatch,
			jsonString: `{
				"firstName":    "",
				"lastName":     "Stark"}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "",
				LastName:  "Stark",
			},
		},
		{
			name:    "Invalid Update Object 2",
			request: http.MethodPatch,
			jsonString: `{
				"firstName":"Tony",
				"lastName":""}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "Tony",
				LastName:  "",
			},
		},
		{
			name:    "Invalid Update Object 3",
			request: http.MethodPatch,
			jsonString: `{
				"firstName":    "",
				"lastName":     ""}`,
			authRequired:        true,
			useMe:               false,
			wrongSignKey:        false,
			expectedErr:         true,
			expectedStatusCode:  http.StatusBadRequest,
			expectedContentType: ContentTypeApplicationJSON,
			expectedReturn:      makeSameUser(),
			updates: &users.Updates{
				FirstName: "",
				LastName:  "",
			},
		},
	}

	for _, c := range cases {
		sh := makeSessionHandler(c.expectedErr)
		if c.wrongSignKey {
			sh.SigningKey = ""
		}
		sessionState := &SessionState{
			SessionTime: time.Now(),
			User:        makeSameUser(),
		}

		sessionid, err := sessions.NewSessionID(sh.SigningKey)
		if err != nil && !c.wrongSignKey {
			t.Errorf("error when making new session id: %v", err)
			return
		}

		if !c.wrongSignKey {
			err = sh.Store.Save(sessionid, sessionState)
			if err != nil {
				t.Errorf("error when saving new session id: %v", err)
				return
			}

			sh.AddNewUserToTrie(sessionState.User)

			jsonByte := []byte(c.jsonString)
			jsonIO := bytes.NewBuffer(jsonByte)
			// fmt.Printf("jsonIO: %v\n", jsonIO)
			var url string
			if c.useMe {
				url = "/v1/users/me"
			} else if c.name == "Invalid ParseInt Test" {
				url = "/v1/users/!@#$*()_"
			} else if c.name == "Invalid GetByID" || c.name == "Not Authorized" {
				url = "/v1/users/1234"
				// } else if c.name == "UserID doesn't equal SessionStateID" {
				// 	url = "/v1/users/5"
			} else {
				url = "/v1/users/1"
			}

			request := httptest.NewRequest(c.request, url, jsonIO)
			request.Header.Add(ContentTypeHeader, c.expectedContentType)
			recorder := httptest.NewRecorder()

			if c.authRequired {
				request.Header.Set("Authorization", "Bearer "+sessionid.String())
			} else {
				request.Header.Set("Authorization", "whatever you want")
			}

			if c.name == "UserID doesn't equal SessionStateID" {
				sessionState.User.ID = 5
			}

			// sh.Users.AddAllToTrie()
			// sh.SpecificUserHandler(recorder, request)

			router := mux.NewRouter()
			router.HandleFunc("/v1/users/{id}", sh.SpecificUserHandler)
			router.ServeHTTP(recorder, request)

			response := recorder.Result()
			//check the response status code
			if response.StatusCode != c.expectedStatusCode {
				t.Errorf("case %s: incorrect status code: expected %d but got %d",
					c.name, c.expectedStatusCode, response.StatusCode)
			}

			user := &users.User{}
			err = json.NewDecoder(response.Body).Decode(user)
			if c.expectedErr && err == nil {
				t.Errorf("case %s: expected error but revieved none", c.name)
			}

			if c.request == http.MethodPatch && c.name == "Invalid ApplyUpdates" &&
				user.FirstName != c.updates.FirstName && user.LastName != c.updates.LastName {
				t.Errorf("case %s: incorrect apply updates: user FirstName(%v) and user LastName(%v), update FirstName(%v) and update LastName(%v)",
					c.name, user.FirstName, user.LastName, c.updates.FirstName, c.updates.LastName)
			}
		}
	}
}

func (sh *SessionHandler) AddNewUserToTrie(u *users.User) {
	username := strings.ToLower(u.UserName)
	firstname := strings.ToLower(u.FirstName)
	lastname := strings.ToLower(u.LastName)

	var wordSplit []string
	if strings.Contains(username, " ") {
		wordSplit = strings.Split(username, " ")
	} else if strings.Contains(firstname, " ") {
		wordSplit = strings.Split(firstname, " ")
	} else if strings.Contains(lastname, " ") {
		wordSplit = strings.Split(lastname, " ")
	}
	for _, elem := range wordSplit {
		sh.Trie.Add(elem, u.ID)
	}
}

func TestSessionsHandler(t *testing.T) {
	cases := []struct {
		name                string
		request             string
		expectedContentType string
		jsonString          string
		credential          *users.Credentials
		expectedStatusCode  int
		expectedErr         bool
		authRequired        bool
	}{
		{
			name:                "Valid Case",
			request:             http.MethodPost,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure"}`,
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusCreated,
			expectedErr:        false,
			authRequired:       true,
		},
		{
			name:                "Invalid HTTP Method",
			request:             http.MethodGet,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure"}`,
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			authRequired:       true,
		},
		{
			name:                "Invalid HTTP Method",
			request:             http.MethodPatch,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure"}`,
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			authRequired:       true,
		},
		{
			name:                "Invalid HTTP Method",
			request:             http.MethodDelete,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure"}`,
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			authRequired:       true,
		},
		{
			name:                "Invalid HTTP Method",
			request:             http.MethodOptions,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure"}`,
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			authRequired:       true,
		},
		{
			name:                "Invalid Content Type",
			request:             http.MethodPost,
			expectedContentType: "text/plain",
			jsonString:          "",
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusUnsupportedMediaType,
			expectedErr:        true,
			authRequired:       true,
		},
		{
			name:                "Invalid Authorization",
			request:             http.MethodPost,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "",
				"password":     ""}`,
			credential: &users.Credentials{
				Email:    "",
				Password: "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			authRequired:       false,
		},
		{
			name:                "Invalid JSON String",
			request:             http.MethodPost,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString:          "",
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "thisPasswordIsSecure",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        true,
			authRequired:       true,
		},
		{
			name:                "Invalid Credentials",
			request:             http.MethodPost,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     ""}`,
			credential: &users.Credentials{
				Email:    "nehay100@gmail.com",
				Password: "",
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedErr:        false,
			authRequired:       true,
		},
		{
			name:                "Fake SigningKey",
			request:             http.MethodPost,
			expectedContentType: ContentTypeApplicationJSON,
			jsonString: `{"email": "nehay100@gmail.com",
				"password":     "thisPasswordIsSecure"}`,
			credential: &users.Credentials{
				Email:    "",
				Password: "",
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        true,
			authRequired:       true,
		},
	}

	for _, c := range cases {
		sh := makeSessionHandler(c.expectedErr)
		if c.name == "Fake SigningKey" {
			sh.SigningKey = ""
		}
		// sessionState := &SessionState{
		// 	SessionTime: time.Now(),
		// 	User:        makeSameUser(),
		// }

		sessionid, err := sessions.NewSessionID(sh.SigningKey)
		if err != nil && c.name != "Fake SigningKey" {
			t.Errorf("error when making new session id: %v", err)
			return
		}

		// err = sh.Store.Save(sessionid, sessionState)
		// if err != nil {
		// 	t.Errorf("error when saving new session id: %v", err)
		// 	return
		// }

		if c.name != "Fake SigningKey" {
			jsonByte := []byte(c.jsonString)
			jsonIO := bytes.NewBuffer(jsonByte)
			// fmt.Printf("jsonIO: %v\n", jsonIO)
			request := httptest.NewRequest(c.request, "/v1/sessions", jsonIO)
			request.Header.Add(ContentTypeHeader, c.expectedContentType)
			recorder := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/v1/sessions", sh.SessionsHandler)
			router.ServeHTTP(recorder, request)

			if c.authRequired {
				request.Header.Set("Authorization", "Bearer "+sessionid.String())
			} else {
				request.Header.Set("Authorization", "whatever you want")
			}

			response := recorder.Result()
			//check the response status code
			if response.StatusCode != c.expectedStatusCode {
				t.Errorf("case %s: incorrect status code: expected %d but got %d",
					c.name, c.expectedStatusCode, response.StatusCode)
			}
			user := &users.User{}
			err = json.NewDecoder(response.Body).Decode(user)
			if c.expectedErr && err == nil {
				t.Errorf("case %s: expected error but revieved none", c.name)
			}
		}

		// if !c.expectedErr && c.expectedReturn.Email != user.Email && !reflect.DeepEqual(c.expectedReturn.PassHash, user.PassHash) &&
		// 	c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
		// 	t.Errorf("case %s: incorrect return: expected %v but revieved %v",
		// 		c.name, c.expectedReturn, user)
		// }
	}
}
func TestSpecificSessionHandler(t *testing.T) {
	cases := []struct {
		name               string
		request            string
		expectedStatusCode int
		expectedErr        bool
		useMine            bool
		wrongSignKey       bool
	}{
		{
			name:               "Valid Case",
			request:            "DELETE",
			expectedStatusCode: http.StatusOK,
			expectedErr:        false,
			useMine:            true,
			wrongSignKey:       false,
		},
		{
			name:               "Invalid HTTP Method",
			request:            http.MethodGet,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			useMine:            true,
			wrongSignKey:       false,
		},
		{
			name:               "Invalid HTTP Method",
			request:            http.MethodPost,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			useMine:            true,
			wrongSignKey:       false,
		},
		{
			name:               "Invalid HTTP Method",
			request:            http.MethodPatch,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			useMine:            true,
			wrongSignKey:       false,
		},
		{
			name:               "Invalid HTTP Method",
			request:            http.MethodOptions,
			expectedStatusCode: http.StatusMethodNotAllowed,
			expectedErr:        true,
			useMine:            true,
			wrongSignKey:       false,
		},
		{
			name:               "Invalid Authorization w/o mine",
			request:            "DELETE",
			expectedStatusCode: http.StatusForbidden,
			expectedErr:        true,
			useMine:            false,
			wrongSignKey:       false,
		},
		{
			name:               "Invalid SigningKey",
			request:            http.MethodDelete,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErr:        true,
			useMine:            true,
			wrongSignKey:       true,
		},
	}

	for _, c := range cases {
		sh := makeSessionHandler(c.expectedErr)
		if c.wrongSignKey {
			sh.SigningKey = ""
		}
		sessionState := &SessionState{
			SessionTime: time.Now(),
			User:        makeSameUser(),
		}

		sessionid, err := sessions.NewSessionID(sh.SigningKey)
		if err != nil && !c.wrongSignKey {
			t.Errorf("error when making new session id: %v", err)
			return
		}

		if !c.wrongSignKey {
			err = sh.Store.Save(sessionid, sessionState)
			if err != nil {
				t.Errorf("error when saving new session id: %v", err)
				return
			}

			var url string
			if c.useMine {
				url = "/v1/sessions/mine"
			} else {
				url = "/v1/sessions/1"
			}
			// jsonByte := []byte(c.jsonString)
			// jsonIO := bytes.NewBuffer(jsonByte)
			request := httptest.NewRequest(c.request, url, nil)
			request.Header.Set("Authorization", "Bearer "+sessionid.String())
			// request.Header.Add(ContentTypeHeader, ContentTypeApplicationJSON)
			recorder := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/v1/sessions/{id}", sh.SpecificSessionHandler)
			router.ServeHTTP(recorder, request)

			// sh.SpecificSessionHandler(recorder, request)

			// if c.authRequired {
			// 	request.Header.Set("Authorization", "Bearer "+sessionid.String())
			// } else {
			// 	request.Header.Set("Authorization", "whatever you want")
			// }

			response := recorder.Result()
			//check the response status code
			if response.StatusCode != c.expectedStatusCode {
				t.Errorf("case %s: incorrect status code: expected %d but got %d",
					c.name, c.expectedStatusCode, response.StatusCode)
			}
			user := &users.User{}
			err = json.NewDecoder(response.Body).Decode(user)
			if c.expectedErr && err == nil {
				t.Errorf("case %s: expected error but revieved none", c.name)
			}

		}

		// if !c.expectedErr && c.expectedReturn.Email != user.Email && !reflect.DeepEqual(c.expectedReturn.PassHash, user.PassHash) &&
		// 	c.expectedReturn.FirstName != user.FirstName && c.expectedReturn.LastName != user.LastName {
		// 	t.Errorf("case %s: incorrect return: expected %v but revieved %v",
		// 		c.name, c.expectedReturn, user)
		// }
	}
}
