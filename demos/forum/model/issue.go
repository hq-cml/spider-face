package model

/*
 * 帖子
 */
import (
	"time"
)

type Issue struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

// format the CreatedAt date to display nicely on the screen
func (issue *Issue) CreatedAtDate() string {
	return issue.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

// get the number of posts in a thread
func (issue *Issue) NumReplies() (count int) {
	rows, err := Db.Query("SELECT count(*) FROM replies where issue_id = ?", issue.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()
	return
}

// get posts to a thread
func (issue *Issue) Replies() (replies []Reply, err error) {
	rows, err := Db.Query("SELECT id, uuid, body, user_id, issue_id, created_at FROM replies where issue_id = ?", issue.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		reply := Reply{}
		if err = rows.Scan(&reply.Id, &reply.Uuid, &reply.Body, &reply.UserId, &reply.IssueId, &reply.CreatedAt); err != nil {
			return
		}
		replies = append(replies, reply)
	}
	rows.Close()
	return
}

// Get all threads in the database and returns it
func GetAllIssues() (issues []Issue, err error) {
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM issues ORDER BY created_at DESC")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Issue{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		issues = append(issues, conv)
	}
	rows.Close()
	return
}

// Get a reply by the UUID
func GetIssueByUUID(uuid string) (conv Issue, err error) {
	conv = Issue{}
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM issues WHERE uuid = ?", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

// Get the user who started this issue
func (issue *Issue) User() (user User) {
	user = User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = ?", issue.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}


