package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Structure for tomb container
type tomb struct {
	SecretID     string `json:"SecretID"`
	UserID       uint   `json:"UserID"` // ID of tree
	Expiration   string `json:"Expires"`
	Contents     string `json:"Contents"`
	ContentsMeta string `json:"ContentsMeta"`
}

// initalise the DB.
func init() {
	var err error
	DB, err = sql.Open("mysql", "root:toor@tcp(localhost:3306)/secrets")
	checkErrorf(err)
	err = DB.Ping() // Ensure connection
	checkErrorf(err)
}

func main() {
	// Router for site
	r := mux.NewRouter()
	r.Headers("Content-Type", "text/html")

	r.HandleFunc("/secret/{secretID}", wHandler(viewSecretHandler)).Methods("GET")
	r.HandleFunc("/secret/{secretID}", wHandler(deleteSecretHandler)).Methods("DELETE")
	r.HandleFunc("/secret/{secretID}", wHandler(modifySecretHandler)).Methods("PUT")
	r.HandleFunc("/secret/{secretID}", wHandler(addSecretHandler)).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func isAuth(r *http.Request) bool {
	return true
}

func isPermissable(r *http.Request) bool {
	return true
}

func viewSecretHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	var tempSecret tomb
	tempSecret.SecretID, _ = vs["secretID"]
	// Is user authenticated to service?
	if isAuth(r) {
		if isPermissable(r) { // is user allowed to access resource?
			row := DB.QueryRow("SELECT * from secrets where secretID = ?", tempSecret.SecretID)

			var err error
			err = row.Scan(&tempSecret.UserID, &tempSecret.Expiration, &tempSecret.Contents, &tempSecret.ContentsMeta, &tempSecret.UserID)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				s, _ := json.Marshal(tempSecret)
				fmt.Fprintln(w, string(s))
			}
		} else { // return satus StatusForbidden, user is NOT allowed to access resource
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
	} else { // return satus StatusUnauthorized, user is NOT authenticated with service
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}

}

func deleteSecretHandler(w http.ResponseWriter, r *http.Request) {

}

func modifySecretHandler(w http.ResponseWriter, r *http.Request) {
}

func addSecretHandler(w http.ResponseWriter, r *http.Request) {

}

// Error handling
func checkErrorf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkErrorp(err error) { // panic
	if err != nil {
		log.Panic(err)
	}
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
		handler(w, r)
	}
	return h
}
