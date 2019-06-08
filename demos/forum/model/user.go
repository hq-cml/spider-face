package model

import "time"

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Create a new session for an existing user ***
func (user *User) CreateSession() (session Session, err error) {
	statement := "insert into sessions (uuid, email, user_id, created_at) values (?, ?, ?, ?) "
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(createUUID(), user.Name, user.Id, time.Now())
	if err != nil {
		return
	}

	return user.Session()
}

// Get the session for an existing user ***
func (user *User) Session() (session Session, err error) {
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = ?", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// Create a new user, save user info into the database ***
func (user *User) Create() (id int64, err error) {
	statement := "insert into users (uuid, name, email, password, created_at) values (?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now())
	if err != nil {
		return 0, err
	}
	id, err = ret.LastInsertId()
	return id, err
}

// Delete user from database ***
func (user *User) Delete() (err error) {
	statement := "delete from users where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

// Update user information in the database ***
func (user *User) Update() (err error) {
	statement := "update users set name = ?, email = ? where id = ?"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.Id)
	return
}

// Delete all users from database ***
func UserDeleteAll() (err error) {
	statement := "delete from users"
	_, err = Db.Exec(statement)
	return
}

// Get all users in the database and returns it ***
func Users() (users []User, err error) {
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	rows.Close()
	return
}

// Get a single user given the email ***
func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = ?", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Get a single user given the UUID ***
func GetUserByUUID(uuid string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = ?", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Get a single user given the UUID ***
func GetUserById(id string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE id = ?", id).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// Create a new Issue
func (user *User) CreateIssue(topic string) (conv Issue, err error) {
	statement := "insert into issues (uuid, topic, user_id, created_at) values (?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(createUUID(), topic, user.Id, time.Now())
	if err != nil {
		return
	}
	id, _ := ret.LastInsertId()

	conv = Issue{}
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM issues WHERE id = ?", id).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

// Create a new reply for a issue
func (user *User) CreateReply(conv Issue, body string) (reply Reply, err error) {
	statement := "insert into replies (uuid, body, user_id, thread_id, created_at) values (?, ?, ?, ?, ?)"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	ret, err := stmt.Exec(createUUID(), body, user.Id, conv.Id, time.Now())
	if err != nil {
		return
	}
	id, _ := ret.LastInsertId()

	reply = Reply{}
	err = Db.QueryRow("SELECT id, uuid, body, user_id, thread_id, created_at FROM replies where id = ?", id).
		Scan(&reply.Id, &reply.Uuid, &reply.Body, &reply.UserId, &reply.ThreadId, &reply.CreatedAt)
	return
}