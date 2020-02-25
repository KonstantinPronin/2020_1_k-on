package main

import (
	"sync"
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Image    string `json:"image"`
}

func (user *User) Update(upUser *User) {
	if upUser.Username != "" {
		user.Username = upUser.Username
	}
	if upUser.Password != "" {
		user.Username = upUser.Password
	}
	if upUser.Email != "" {
		user.Email = upUser.Email
	}
}

type UserStorage struct {
	mutex sync.RWMutex
	users map[string]*User
	count uint
}

func CreateUserStorage() *UserStorage {
	//for test
	user := User{
		Id:       0,
		Username: "admin",
		Password: "admin",
		Email:    "admin@example.ru",
	}

	return &UserStorage{
		mutex: sync.RWMutex{},
		users: map[string]*User{user.Username: &user},
		count: 1,
	}
}

func (userStorage *UserStorage) Contains(login string) bool {
	_, ok := userStorage.users[login]
	return ok
}

func (userStorage *UserStorage) Add(user *User) (uint, error) {
	userStorage.mutex.Lock()
	defer userStorage.mutex.Unlock()

	id := userStorage.count
	user.Id = id
	userStorage.users[user.Username] = user
	userStorage.count++

	return id, nil
}

func (userStorage *UserStorage) GetById(id uint) (*User, bool) {
	userStorage.mutex.RLock()
	defer userStorage.mutex.RUnlock()

	for _, user := range userStorage.users {
		if user.Id == id {
			return user, true
		}
	}

	return nil, false
}

func (userStorage *UserStorage) GetByName(login string) (*User, bool) {
	userStorage.mutex.RLock()
	defer userStorage.mutex.RUnlock()

	if !userStorage.Contains(login) {
		return nil, false
	}

	return userStorage.users[login], true
}
