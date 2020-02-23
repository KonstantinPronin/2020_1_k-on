package main

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestUserStorage_Add(t *testing.T) {
	userStorage := UserStorage{
		mutex: sync.RWMutex{},
		users: make(map[string]*User),
		count: 0,
	}
	user := User{
		Username: "test",
		Password: "test",
	}

	_, err := userStorage.Add(&user)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, user, *userStorage.users[user.Username])
}

func TestUserStorage_Contains(t *testing.T) {
	user := User{
		Username: "test",
		Password: "test",
	}
	userStorage := UserStorage{
		mutex: sync.RWMutex{},
		users: map[string]*User{user.Username: &user},
		count: 0,
	}

	contains := userStorage.Contains(user.Username)

	assert.True(t, contains)
}

func TestUserStorage_GetById(t *testing.T) {
	user := User{
		Id:       10,
		Username: "test",
		Password: "test",
	}
	userStorage := UserStorage{
		mutex: sync.RWMutex{},
		users: map[string]*User{user.Username: &user},
		count: 0,
	}

	stored, ok := userStorage.GetById(user.Id)

	assert.True(t, ok)
	assert.Equal(t, user, *stored)
}

func TestUserStorage_GetByName(t *testing.T) {
	user := User{
		Username: "test",
		Password: "test",
	}
	userStorage := UserStorage{
		mutex: sync.RWMutex{},
		users: map[string]*User{user.Username: &user},
		count: 0,
	}

	stored, ok := userStorage.GetByName(user.Username)

	assert.True(t, ok)
	assert.Equal(t, user, *stored)
}
