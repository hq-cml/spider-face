package model

import (
	"testing"
	"database/sql"
	"strconv"
	"github.com/hq-cml/spider-face/utils/helper"
)

// test data
var users = []User{
	{
		Name:     "Peter Jones",
		Email:    "peter@gmail.com",
		Password: "peter_pass",
	},
	{
		Name:     "John Smith",
		Email:    "john@gmail.com",
		Password: "john_pass",
	},
}

func setup() {
	SessionDeleteAll()
	UserDeleteAll()
}

func Test_UserCreate(t *testing.T) {
	setup()

	id, err := users[0].Create();
	if  err != nil {
		t.Error(err, "Cannot create user.")
	}
	t.Logf("Create User: %v", id)

	u, err := GetUserByEmail(users[0].Email)
	if err != nil {
		t.Errorf("User not created. %v", err)
	}
	if users[0].Email != u.Email {
		t.Errorf("User retrieved is not the same as the one created.")
	}
	t.Logf("Create User: %v", helper.JsonEncode(u))
}

func Test_UserDelete(t *testing.T) {
	setup()
	_, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}

	u, err := GetUserByEmail(users[0].Email)
	if err != nil {
		t.Errorf("User not created. %v", err)
	}

	if err := u.Delete(); err != nil {
		t.Error(err, "- Cannot delete user")
	}
	_, err = GetUserByEmail(users[0].Email)
	if err != sql.ErrNoRows {
		t.Error(err, "- User not deleted.")
	}
}

func Test_UserUpdate(t *testing.T) {
	setup()
	_, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}

	u, err := GetUserByEmail(users[0].Email)
	if err != nil {
		t.Errorf("User not created. %v", err)
	}
	u.Name = "Random User"
	if err := u.Update(); err != nil {
		t.Error(err, "- Cannot update user")
	}

	u, err = GetUserByEmail(users[0].Email)
	if err != nil {
		t.Error(err, "- Cannot get user")
	}
	if u.Name != "Random User" {
		t.Error(err, "- User not updated")
	}
}

func Test_UserByUUID(t *testing.T) {
	setup()
	id, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}
	u, err := GetUserById(strconv.Itoa(int(id)))
	if err != nil {
		t.Errorf("User not created. %v", err)
	}

	u, err = GetUserByUUID(u.Uuid)
	if err != nil {
		t.Error(err, "User not created.")
	}
	if users[0].Email != u.Email {
		t.Errorf("User retrieved is not the same as the one created.")
	}
}

func Test_Users(t *testing.T) {
	setup()
	for _, user := range users {
		if _, err := user.Create(); err != nil {
			t.Error(err, "Cannot create user.")
		}
	}
	u, err := Users()
	if err != nil {
		t.Error(err, "Cannot retrieve users.")
	}
	if len(u) != 2 {
		t.Error(err, "Wrong number of users retrieved")
	}
	if u[0].Email != users[0].Email {
		t.Error(u[0], users[0], "Wrong user retrieved")
	}
}

func Test_CreateSession(t *testing.T) {
	setup()
	id, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}
	u, err := GetUserById(strconv.Itoa(int(id)))
	if err != nil {
		t.Errorf("User not created. %v", err)
	}

	session, err := u.CreateSession()
	if err != nil {
		t.Error(err, "Cannot create session")
	}
	if session.UserId != u.Id {
		t.Error("User not linked with session")
	}
}

func Test_GetSession(t *testing.T) {
	setup()
	id, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}
	u, err := GetUserById(strconv.Itoa(int(id)))
	if err != nil {
		t.Errorf("User not created. %v", err)
	}

	session, err := u.CreateSession()
	if err != nil {
		t.Error(err, "Cannot create session")
	}

	s, err := u.Session()
	if err != nil {
		t.Error(err, "Cannot get session")
	}
	if s.Id == 0 {
		t.Error("No session retrieved")
	}
	if s.Id != session.Id {
		t.Error("Different session retrieved")
	}
}

func Test_checkValidSession(t *testing.T) {
	setup()
	id, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}
	u, err := GetUserById(strconv.Itoa(int(id)))
	if err != nil {
		t.Errorf("User not created. %v", err)
	}

	session, err := u.CreateSession()
	if err != nil {
		t.Error(err, "Cannot create session")
	}

	uuid := session.Uuid

	s := Session{Uuid: uuid}
	valid, err := s.Check()
	if err != nil {
		t.Error(err, "Cannot check session")
	}
	if valid != true {
		t.Error(err, "Session is not valid")
	}
}

func Test_checkInvalidSession(t *testing.T) {
	setup()
	s := Session{Uuid: "123"}
	valid, err := s.Check()
	if err == nil {
		t.Error(err, "Session is not valid but is validated")
	}
	if valid == true {
		t.Error(err, "Session is valid")
	}

}

func Test_DeleteSession(t *testing.T) {
	setup()
	id, err := users[0].Create();
	if err != nil {
		t.Error(err, "Cannot create user.")
	}
	u, err := GetUserById(strconv.Itoa(int(id)))
	if err != nil {
		t.Errorf("User not created. %v", err)
	}

	session, err := u.CreateSession()
	if err != nil {
		t.Error(err, "Cannot create session")
	}

	err = session.DeleteByUUID()
	if err != nil {
		t.Error(err, "Cannot delete session")
	}
	s := Session{Uuid: session.Uuid}
	valid, err := s.Check()
	if err == nil {
		t.Error(err, "Session is valid even though deleted")
	}
	if valid == true {
		t.Error(err, "Session is not deleted")
	}
}
