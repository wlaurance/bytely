package database

import "testing"

func TestUserSetup(t *testing.T) {
	DB_URL = "user=rocconicosia dbname=linkthingtest sslmode=disable"
}

func TestUserCreate(t *testing.T) {
	user := NewUser("herp@derp.io", "hotdogs")
	err := user.Insert()
	if err != nil {
		t.Fatal(err)
	}

	if user.Id == 0 {
		t.Fatal("User object not properly initialized")
	}
}

func TestUserGet(t *testing.T) {
	user, err := GetUser("herp@derp.io", "hotdogs")
	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("User object not properly initialized")
	}
}

func TestUserSave(t *testing.T) {
	user, _ := GetUser("herp@derp.io", "hotdogs")

	// Try updating user info.
	user.Email = "meow@derp.io"
	err := user.Save()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the user information was actually updated.
	_, err = GetUser("meow@derp.io", "hotdogs")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserDelete(t *testing.T) {
	user, _ := GetUser("meow@derp.io", "hotdogs")
	err := user.Delete()
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that the record was actually deleted.
	user, err = GetUser("meow@derp.io", "hotdogs")
	if err == nil {
		t.Fatal("User record not properly deleted")
	}
}
