package handler

import (
	"fmt"
	"net/http"
)

type Session struct{}

var matchbooktoken string

func (o Session) Login(w http.ResponseWriter, r *http.Request) {
	matchbooktoken = loadMatchboookToken()
	fmt.Printf("Got session token %v", matchbooktoken)
}

func (o Session) GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Current session token %v", matchbooktoken)
}

func (o Session) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout of current session")
}

func (o Session) CreateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create event data for Id")
}
