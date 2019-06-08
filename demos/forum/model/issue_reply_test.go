package model

//
import (
	"testing"
)

//Delete all threads from database
func IssuesDeleteAll() (err error) {
	statement := "delete from issues"
	_, err = Db.Exec(statement)
	return
}

func RepliesDeleteAll() (err error) {
	statement := "delete from replies"
	_, err = Db.Exec(statement)
	return
}

func Test_CreateIssue(t *testing.T) {
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

	conv, err := u.CreateIssue("My first issue.")
	if err != nil {
		t.Error(err, "Cannot create thread")
	}
	if conv.UserId != u.Id {
		t.Error("User not linked with thread")
	}
}
