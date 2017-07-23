package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Structure for tomb container
type secret struct {
	SecretID     int    `json:"SecretID"`
	Name         string `json:"Name"`
	UserID       int    `json:"UserID"` // ID of tree
	Expiration   string `json:"Expires"`
	Contents     string `json:"Contents"`
	ContentsMeta string `json:"ContentsMeta"`
}

type user struct {
	UserID   int    `json:"UserID"` // ID of tree
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Secrets []secret

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

	r.HandleFunc("/secret/{secretID}", wHandler(viewSecretHandler)).Methods("GET")
	r.HandleFunc("/secret/{secretID}", wHandler(deleteSecretHandler)).Methods("DELETE")
	r.HandleFunc("/secret/{secretID}", wHandler(modifySecretHandler)).Methods("PUT")
	r.HandleFunc("/secret/{secretID}", wHandler(addSecretHandler)).Methods("POST")
	r.HandleFunc("/user/{userID}", wHandler(viewUserHandler)).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func isAuth(r *http.Request) bool {
	return true
}

func isPermissable(r *http.Request) bool {
	return true
}

func viewUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	vs := mux.Vars(r)
	var err error

	gd, _ := strconv.Atoi(vs["userID"])
	secretrows, err := DB.Query("SELECT expiration, secretID, name from secrets where userID = ?", gd)
	if err != nil {
		checkErrorf(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	defer secretrows.Close()
	Secrets := []secret{}

	for secretrows.Next() {
		var st secret
		err := secretrows.Scan(&st.Expiration, &st.SecretID, &st.Name)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		Secrets = append(Secrets, st)
	}
	fmt.Fprintf(w, "<h3>%d</h3>", gd)
	for _, st := range Secrets {
		fmt.Fprintf(w, "<a href=\"/secret/%d\">%s</a> <br>", st.SecretID, st.Name)

	}
}

// curl -i -X GET http://localhost:8080/secret/0
func viewSecretHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)
	var tempSecret secret
	var err error

	tempSecret.SecretID, _ = strconv.Atoi(vs["secretID"])
	row := DB.QueryRow("SELECT * from secrets where secretID = ?", tempSecret.SecretID)
	err = row.Scan(&tempSecret.SecretID, &tempSecret.Expiration, &tempSecret.Contents, &tempSecret.ContentsMeta, &tempSecret.UserID, &tempSecret.Name)

	if err != nil {
		checkErrorp(err)

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		s, _ := json.Marshal(tempSecret) // Marshall Tree
		fmt.Fprintln(w, string(s))       // Print to user
	}

}

func deleteSecretHandler(w http.ResponseWriter, r *http.Request) {
	vs := mux.Vars(r)

	_, err := DB.Exec("DELETE FROM secrets WHERE secretID=?", vs["secretID"])
	fmt.Fprintln(w, err)

}

func modifySecretHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.FormValue("a"))
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
		w.Header().Set("Content-Type", "application/json")
		handler(w, r)
	}
	return h
}
