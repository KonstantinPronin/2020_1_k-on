package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io/ioutil"
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

func (userHandler *UserHandler) isAuth(r *http.Request) bool {
	session, err := r.Cookie("session_id")
	if session == nil || err == http.ErrNoCookie {
		return false
	}

	if _, ok := userHandler.sessions[session.Value]; !ok {
		return false
	}

	return true
}

func (userHandler *UserHandler) auth(login, password string, w http.ResponseWriter) {
	if login == "" {
		http.Error(w, `bad login`, http.StatusBadRequest)
		return
	}

	user, ok := userHandler.users.GetByName(login)
	if !ok {
		http.Error(w, `no user`, http.StatusUnauthorized)
		return
	}
	if user.Password != password {
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

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Method", "POST")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	if userHandler.isAuth(r) {
		return
	}

	in := new(User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(in)
	if err != nil {
		http.Error(w, `bad parameters`, http.StatusBadRequest)
		return
	}

	userHandler.auth(in.Username, in.Password, w)
}

func (userHandler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if !userHandler.isAuth(r) {
		http.Error(w, `no session`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")

	delete(userHandler.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
}

func (userHandler *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	if userHandler.isAuth(r) {
		http.Error(w, `already login`, http.StatusForbidden)
		return
	}

	user := new(User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, `bad parameters`, http.StatusBadRequest)
		return
	}
	if userHandler.users.Contains(user.Username) {
		http.Error(w, `such user already exist`, http.StatusBadRequest)
		return
	}

	_, err = userHandler.users.Add(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userHandler.auth(user.Username, user.Password, w)
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (userHandler *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userHandler.Get(w, r)
	case http.MethodPost:
		userHandler.Update(w, r)
	}
}

func (userHandler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	if !userHandler.isAuth(r) {
		http.Error(w, `no session`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")
	id := userHandler.sessions[session.Value]

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `no user`, http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (userHandler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if !userHandler.isAuth(r) {
		http.Error(w, `no session`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")
	id := userHandler.sessions[session.Value]

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `no user`, http.StatusNotFound)
		return
	}

	upUser := new(User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(upUser)
	if err != nil {
		http.Error(w, `bad parameters`, http.StatusBadRequest)
		return
	}

	user.Update(upUser)
}

func (userHandler *UserHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if !userHandler.isAuth(r) {
		http.Error(w, `no session`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")
	id := userHandler.sessions[session.Value]

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `no user`, http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		http.Error(w, `bad data form`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `bad data form`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, `cannot read file`, http.StatusInternalServerError)
		return
	}

	filepath, err := WriteFile(ImageDirectory, header.Filename, data)
	if err != nil {
		http.Error(w, `cannot save file`, http.StatusInternalServerError)
		return
	}

	user.Image = filepath
}

func (userHandler *UserHandler) GetImage(w http.ResponseWriter, r *http.Request) {
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

	data, err := ReadFile(user.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	image := ImageJson{Image: data}
	err = json.NewEncoder(w).Encode(&image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
