package database

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/rand"

	"database/sql"
	_ "github.com/lib/pq"
)

// User is just a Golang representation of a User record in Postgres.
type User struct {
	Id       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// NewUser is meant to take raw user data received from a user and prep it
// for database usage. After the user is returned, Insert() can be called
// to add it to the database.
func NewUser(email, password string) *User {
	// Hash your passwords, kids.
	hashedPassword := sha512.Sum512([]byte(password))

	// Generate an API token.
	token := generateAuthToken(email, password)

	return &User{
		0,
		email,
		hex.EncodeToString(hashedPassword[:]),
		token,
	}
}

// generateAuthToken creates a unique auth token given a user's
// credentials. It is guaranteed to be unique in-so-far as the sha1
// algorithm is guaranteed to produce hashes with no collisions
// (sha1 is currently considered secure).
func generateAuthToken(email, password string) string {
	// Generate a random integer.
	rawToken := make([]byte, 20)
	binary.PutVarint(rawToken, rand.Int63())

	credentials := []byte(email + password)
	token := sha1.Sum(append(rawToken, credentials...))

	return hex.EncodeToString(token[:])
}

// GetUser is meant for applications such as logging in, where a raw password
// and username are received and the entire user profile is desired.
func GetUser(email, password string) (*User, error) {
	user := NewUser(email, password)

	err := WithDatabase(func(db *sql.DB) error {
		row := db.QueryRow(
			`select id, email, password, token from users where email = $1 and password = $2`,
			user.Email, user.Password,
		)

		return row.Scan(&user.Id, &user.Email, &user.Password, &user.Token)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByToken retrieves a user by its auth token.
func GetUserByToken(token string) (*User, error) {
	user := new(User)

	err := WithDatabase(func(db *sql.DB) error {
		row := db.QueryRow(`select id, email, password, token 
                            from users where token = $1`, token)
		return row.Scan(&user.Id, &user.Email, &user.Password, &user.Token)
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

// Insert attempts to insert the User struct into the database as a new record.
func (this *User) Insert() error {
	return WithDatabase(func(db *sql.DB) error {
		// Ensure that a user with provided email doesn't already exist.
		rows, _ := db.Query("select 1 from users where email = $1", this.Email)
		if rows.Next() {
			return errors.New("User already exists.")
		}

		dberr := db.QueryRow(
			"insert into users(email, password, token) values($1, $2, $3) returning id",
			this.Email, this.Password, this.Token,
		).Scan(&this.Id)

		return dberr
	})
}

// Save runs an update query on all fields of the User struct. This has the
// effect of only updating those fields which have changed.
func (this *User) Save() error {
	return WithDatabase(func(db *sql.DB) error {
		_, dberr := db.Exec(
			"update users set email = $1, password = $2, token = $3 where id = $4",
			this.Email, this.Password, this.Token, this.Id,
		)
		return dberr
	})
}

// Delete removes the User struct's corresponding record from the database.
func (this *User) Delete() error {
	return WithDatabase(func(db *sql.DB) error {
		_, dberr := db.Exec("delete from users where id = $1", this.Id)
		return dberr
	})
}
