package handler

import (
	"fmt"
	"net/http"
)

type Order struct{}

func (o Order) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get session token")
}

func (o Order) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Logout of current session")
}

func (o Order) CreateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create event data for Id")
}
