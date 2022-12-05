package modeltests

import (
	"log"
	"testing"

	"github.com/rajaghaneshan/go-fullstack/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}
	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the users: %v", err)
		return
	}
	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		ID:       1,
		Email:    "test@gmail.com",
		Nickname: "test",
		Password: "password",
	}
	saveduser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("this is the error creating user: %v", err)
		return
	}

	assert.Equal(t, newUser.ID, saveduser.ID)
	assert.Equal(t, newUser.Email, saveduser.Email)
	assert.Equal(t, newUser.Nickname, saveduser.Nickname)
}

func TestGetUserByID(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	foundUser, err := userInstance.FindUserById(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error getting one user: %v", err)
		return
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
}

func TestUpdateAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed user: %v", err)
	}

	userUpdate := models.User{
		ID:       1,
		Nickname: "updatedUser",
		Email:    "updatedemail@gmail.com",
		Password: "password",
	}
	updatedUser, err := userUpdate.UpdateAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating one user %v", err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("cannot seed user table :%v", err)
	}
	isDeleted, err := userInstance.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("this is the error updating the user: %v", err)
		return
	}
	//use any one - both are same
	assert.Equal(t, isDeleted, int64(1))
	// assert.Equal(t, int(isDeleted), 1)
}
