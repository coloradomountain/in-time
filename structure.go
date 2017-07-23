package main

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

// Secrets ...
type Secrets []secret
