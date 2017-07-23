package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// initialise DB
	InitDB("root:toor@/secrets")
	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/secret/{secretID}", wHandler(viewSecretHandler)).Methods("GET")
	r.HandleFunc("/secret/{secretID}", wHandler(deleteSecretHandler)).Methods("DELETE")
	r.HandleFunc("/secret/{secretID}", wHandler(modifySecretHandler)).Methods("PUT")
	r.HandleFunc("/secret/{secretID}", wHandler(addSecretHandler)).Methods("POST")
	r.HandleFunc("/user/{userID}", wHandler(viewUserHandler)).Methods("GET")
	// Serve
	http.ListenAndServe(":8080", r)
}

// Is the user authenticated?
func isAuth(r *http.Request) bool {
	return true
}

// Does the user have permisson to access resource?
func isPermissable(r *http.Request) bool {
	return true
}

// curl -i -X GET http://localhost:8080/user/{userID}
func viewUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

}

// curl -i -X GET http://localhost:8080/secret/{secretID}
func viewSecretHandler(w http.ResponseWriter, r *http.Request) {

}

// curl --data "secretID=0&message=new-message" -i -X PUT http://localhost:8080/secret/{secretID}
func modifySecretHandler(w http.ResponseWriter, r *http.Request) {
}

// curl --data "userID=0&expiration=2017-07-22%2023%3A48%3A2&secretID=0&contents=new-message&name=New-secret-for-you" -i -X PUT http://localhost:8080/secret/{secretID}
func addSecretHandler(w http.ResponseWriter, r *http.Request) {

}

// curl -i -X GET http://localhost:8080/secret/{secretID}
func deleteSecretHandler(w http.ResponseWriter, r *http.Request) {

}

// Wrap handler function
func wHandler(
	handler func(w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {

	h := func(w http.ResponseWriter, r *http.Request) {
		if !isAuth(r) { // Ensure user is authenticated
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized) // return status 401 Unauthorized
			return
		}
		if !isPermissable(r) { // Ensure user has proper access level
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden) // return status 403 forbidden
			return
		}
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
	return h
}
