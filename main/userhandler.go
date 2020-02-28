package main

import (
	"encoding/json"
	"fmt"
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

func (userHandler *UserHandler) auth(login, password string, w http.ResponseWriter) bool {
	w.Header().Set("Content-Type", "application/json")
	if login == "" {
		http.Error(w, `{"error":"bad login"}`, http.StatusBadRequest)
		return false
	}

	user, ok := userHandler.users.GetByName(login)
	if !ok {
		http.Error(w, `{"error":"no user"}`, http.StatusUnauthorized)
		return false
	}
	if user.Password != password {
		http.Error(w, `{"error":"bad password"}`, http.StatusBadRequest)
		return false
	}

	sessionId := uuid.New().String()

	userHandler.sessions[sessionId] = user.Id
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   sessionId,
		Path:    "/",
		Expires: time.Now().Add(10 * time.Hour),
	}

	http.SetCookie(w, cookie)
	return true
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodPost {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if userHandler.isAuth(r) {
		http.Error(w, `{"error":"already login"}`, http.StatusForbidden)
		return
	}

	in := new(User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(in)
	if err != nil {
		http.Error(w, `{"error":"bad parameters"}`, http.StatusBadRequest)
		return
	}

	if userHandler.auth(in.Username, in.Password, w) {
		fmt.Fprint(w, `{"Answer":"OK"}`)
	}
}

func (userHandler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodDelete {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !userHandler.isAuth(r) {
		http.Error(w, `{"error":"no session"}`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")

	delete(userHandler.sessions, session.Value)

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	fmt.Fprint(w, `{"Answer":"OK"}`)
}

func (userHandler *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodPost {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if userHandler.isAuth(r) {
		http.Error(w, `{"error":"already login"}`, http.StatusForbidden)
		return
	}

	user := new(User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, `{"error":"bad in parameters"}`, http.StatusBadRequest)
		return
	}
	if userHandler.users.Contains(user.Username) {
		http.Error(w, `{"error":"such user already exist"`, http.StatusBadRequest)
		return
	}

	_, err = userHandler.users.Add(user)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}

	userHandler.auth(user.Username, user.Password, w)
	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}
}

func (userHandler *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodGet {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !userHandler.isAuth(r) {

		http.Error(w, `{"error":"no session"}`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")
	id := userHandler.sessions[session.Value]

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `{"error":"no user"}`, http.StatusNotFound)
		return
	}

	if user.Image != "" {
		data, err := ReadFile(user.Image)
		if err != nil {
			http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
			return
		}
		user.ImageBase64 = data
	}

	err := json.NewEncoder(w).Encode(&user)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}
}

func (userHandler *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST,GET")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodPost {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !userHandler.isAuth(r) {
		http.Error(w, `{"error":"no session"}`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")
	id := userHandler.sessions[session.Value]

	upUser := new(User)
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := decoder.Decode(upUser)
	if err != nil {
		http.Error(w, `{"error":"bad parameters"}`, http.StatusBadRequest)
		return
	}

	if !userHandler.users.Update(id, upUser) {
		http.Error(w, `{"error":"can not save updates"}`, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, `{"Answer":"OK"}`)
}

func (userHandler *UserHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodPut {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !userHandler.isAuth(r) {
		http.Error(w, `{"error":"no session"}`, http.StatusUnauthorized)
		return
	}
	session, _ := r.Cookie("session_id")
	id := userHandler.sessions[session.Value]

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `{"error":"no user"}`, http.StatusNotFound)
		return
	}

	err := r.ParseMultipartForm(MaxFileSize)
	if err != nil {
		http.Error(w, `{"error":"bad data form"}`, http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, `{"error":"bad data form"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, `{"error":"cannot read file"}`, http.StatusInternalServerError)
		return
	}

	filepath, err := WriteFile(ImageDirectory, header.Filename, data)
	if err != nil {
		http.Error(w, `{"error":"cannot save file"}`, http.StatusInternalServerError)
		return
	}

	user.Image = filepath
	fmt.Fprint(w, `{"Answer":"OK"}`)
}

func (userHandler *UserHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method != http.MethodGet {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if id < 0 || err != nil {
		http.Error(w, `{"error":"bad id"}`, http.StatusBadRequest)
		return
	}

	user, ok := userHandler.users.GetById(uint(id))
	if !ok {
		http.Error(w, `{"error":"no user"}`, http.StatusNotFound)
		return
	}

	data, err := ReadFile(user.Image)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}

	image := ImageJson{Image: data}
	err = json.NewEncoder(w).Encode(&image)
	if err != nil {
		http.Error(w, `{"error":"`+string(err.Error())+`"}`, http.StatusInternalServerError)
		return
	}
}
