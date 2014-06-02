package database

import "testing"

var user *User

func TestLinkSetup(t *testing.T) {
	DB_URL = "user=rocconicosia dbname=linkthingtest sslmode=disable"

	user = NewUser("herp@derp.io", "meow")
	user.Insert()
}

func TestLinkCreate(t *testing.T) {
	if user == nil {
		t.SkipNow()
	}

	link := NewLink("http://rocco.io", user.Id)
	err := link.Insert()
	if err != nil {
		t.Fatal(err)
	}

	if link.Hash == "" || link.Id == 0 {
		t.Fatal("Link object not properly initialized")
	}
}

func TestLinkGet(t *testing.T) {
	link, err := GetLink("aaaaaa")
	if err != nil {
		t.Fatal(err)
	}

	if link == nil {
		t.Fatal("Link object not properly initialized")
	}
}

func TestLinkGetForUser(t *testing.T) {
	if user == nil {
		t.SkipNow()
	}

	NewLink("http://google.com", user.Id).Insert()
	NewLink("http://godoc.org", user.Id).Insert()
	NewLink("http://insightpool.com", user.Id).Insert()

	links, err := GetLinksForUser(user.Id)
	if err != nil {
		t.Fatal(err)
	}

	if len(links) != 4 {
		t.Fatal("Not all links were returned")
	}

	links, err = GetLinksForUser(1337)
	if err != nil {
		t.Fatal(err)
	}

	if len(links) != 0 {
		t.Fatal("Links are returned for nonexistent user")
	}
}

func TestLinkSave(t *testing.T) {
	link, _ := GetLink("aaaaaa")

	link.OriginalURL = "http://wtf.ninja"
	err := link.Save()
	if err != nil {
		t.Fatal(err)
	}

	link, _ = GetLink("aaaaaa")
	if link.OriginalURL != "http://wtf.ninja" {
		t.Fatal("Link information is not saved properly")
	}
}

func TestLinkDelete(t *testing.T) {
	link, _ := GetLink("aaaaaa")

	err := link.Delete()
	if err != nil {
		t.Fatal(err)
	}

	_, err = GetLink("aaaaaa")
	if err == nil {
		t.Fatal("Links deletion is borked")
	}
}

func TestLinkTeardown(t *testing.T) {
	links, _ := GetLinksForUser(user.Id)
	for _, link := range links {
		link.Delete()
	}
	user.Delete()
}
