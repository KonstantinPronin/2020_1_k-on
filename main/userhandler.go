package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	sessions map[string]uint
	users    *UserStorage
}

func createUserHandler() *UserHandler {
	return &UserHandler{
		sessions: make(map[string]uint),
		users:    CreateUserStorage(),
	}
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	if login == "" {
		http.Error(w, `bad login`, http.StatusBadRequest)
		return
	}
	user, ok := userHandler.users.GetByName(login)

	if !ok {
		http.Error(w, `no user`, http.StatusNotFound)
		return
	}

	if user.Password != r.FormValue("password") {
		http.Error(w, `bad password`, http.StatusBadRequest)
		return
	}

	sessionId := uuid.New().String()

	userHandler.sessions[sessionId] = user.Id
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Expires: time.Now().Add(10 * time.Hour),
	}

	http.SetCookie(w, cookie)
}

func (userHandler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if session == nil || err == http.ErrNoCookie {
		http.Error(w, `no session`, http.StatusUnauthorized)
		return
	}

	if _, ok := userHandler.sessions[session.Value]; !ok {
		http.Error(w, `no session`, http.StatusUnauthorized)
		return
	}

	delete(userHandler.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (userHandler *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	email := r.FormValue("email")
	if login == "" || password == "" {
		http.Error(w, `bad parameters`, http.StatusBadRequest)
		return
	}

	user := &User{
		Username: login,
		Password: password,
		Email:    email,
	}

	_, err := userHandler.users.Add(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (userHandler *UserHandler) UserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if id < 0 || err != nil {
		http.Error(w, `bad id`, http.StatusBadRequest)
		return
	}

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `no user`, http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
