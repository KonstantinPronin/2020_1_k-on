package main

import (
	"sync"
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserStorage struct {
	mutex sync.RWMutex
	users map[uint]*User
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
		users: map[uint]*User{0: &user},
		count: 1,
	}
}

func (userStorage *UserStorage) Contains(id uint) bool {
	_, ok := userStorage.users[id]
	return ok
}

func (userStorage *UserStorage) Add(user *User) (uint, error) {
	userStorage.mutex.Lock()
	defer userStorage.mutex.Unlock()

	id := userStorage.count
	user.Id = id
	userStorage.users[id] = user
	userStorage.count++

	return id, nil
}

func (userStorage *UserStorage) GetById(id uint) (*User, bool) {
	userStorage.mutex.RLock()
	defer userStorage.mutex.RUnlock()

	if !userStorage.Contains(id) {
		return nil, false
	}

	return userStorage.users[id], true
}

func (userStorage *UserStorage) GetByName(login string) (*User, bool) {
	userStorage.mutex.RLock()
	defer userStorage.mutex.RUnlock()

	for _, user := range userStorage.users {
		if user.Username == login {
			return user, true
		}
	}

	return nil, false
}
