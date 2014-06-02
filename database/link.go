package database

import (
	"bytes"
	"database/sql"

	_ "github.com/lib/pq"
)

// Link is a Golang representation of a Link record in Postgres.
type Link struct {
	Id          uint64 `json:"-"`
	OriginalURL string `json:"original_url"`
	Hash        string `json:"hash"`
	Hits        uint   `json:"hits"`
	LastHit     string `json:"last_hit"`
	MobileHits  uint   `json:"mobile_hits"`
	CreatedOn   string `json:"created_on"`
	UserId      uint64 `json:"-"`
}

// NewLink creates a new Link struct which is owned by the user with the
// provided id. Links returned by this function will not be valid, as their
// Id field will be set to 0. They should either be inserted into the
// database using Insert() immediately or GetLink() should be used instead.
func NewLink(url string, userId uint64) *Link {
	return &Link{0, url, "", 0, "", 0, "", userId}
}

// GetLink retrives a single Link record by its hash value.
func GetLink(hash string) (*Link, error) {
	link := new(Link)

	err := WithDatabase(func(db *sql.DB) error {
		row := db.QueryRow(
			`select id, url, hash, hits, last_hit, mobile_hits, created_on, user_id 
			from links 
			where hash = $1`,
			hash,
		)

		return row.Scan(&link.Id, &link.OriginalURL, &link.Hash,
			&link.Hits, &link.LastHit, &link.MobileHits, &link.CreatedOn,
			&link.UserId)
	})

	if err != nil {
		return nil, err
	}

	return link, nil
}

// GetLinksForUser retrieves all link records that were created by the
// user with the specified id.
func GetLinksForUser(userId uint64) ([]*Link, error) {
	links := make([]*Link, 0)

	var rows *sql.Rows
	err := WithDatabase(func(db *sql.DB) error {
		dbrows, dberr := db.Query(
			`select id, url, hits, hash, last_hit, mobile_hits, created_on from links 
			where user_id = $1`,
			userId,
		)

		// Not sure why this is necessary - := is broken with closures?
		rows = dbrows
		return dberr
	})

	if err != nil {
		return nil, err
	}

	// Iterate through all returned rows and populate a Link for each.
	for rows.Next() {
		link := new(Link)
		link.UserId = userId

		if err = rows.Scan(&link.Id, &link.OriginalURL, &link.Hits,
			&link.Hash, &link.LastHit, &link.MobileHits, &link.CreatedOn); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}

// Insert attempts to insert the Link into the database as a record.
// It first runs an insert operation which returns the primary id of
// the new record if successful. The link hash is then generated with
// this unique id, which calls for an update query to add this value
// to the new record.
func (this *Link) Insert() error {
	return WithDatabase(func(db *sql.DB) error {
		// Add new record with empty hash field and return primary id.
		var id uint64
		dberr := db.QueryRow(
			`insert into links(url, hash, hits, last_hit, mobile_hits, created_on, user_id) 
			 values($1, $2, $3, $4, $5, $6, $7) returning id`,
			this.OriginalURL, "", this.Hits, "", this.MobileHits, this.CreatedOn, this.UserId,
		).Scan(&id)
		if dberr != nil {
			return dberr
		}

		// Save the table id so the current Link object is valid.
		this.Id = id

		// Generate link hash.
		this.Hash = hashInt(id)

		// Update new record with link hash.
		_, dberr = db.Exec(
			"update links set hash = $1 where id = $2",
			this.Hash, this.Id,
		)
		return dberr
	})
}

// Save commits all of the Link struct's fields to its database record. If
// any of them have changed, it will update the database record accordingly.
func (this *Link) Save() error {
	return WithDatabase(func(db *sql.DB) error {
		_, dberr := db.Exec(
			`update links 
			set url = $1, hash = $2, hits = $3, last_hit = $4, mobile_hits = $5,
			created_on = $6 
			where id = $7`,
			this.OriginalURL, this.Hash, this.Hits, this.LastHit,
			this.MobileHits, this.CreatedOn, this.Id,
		)
		return dberr
	})
}

// Delete removes the Link's corresponding database record.
func (this *Link) Delete() error {
	return WithDatabase(func(db *sql.DB) error {
		_, dberr := db.Exec("delete from links where id = $1", this.Id)
		return dberr
	})
}

// hashInt produces a six-character hash which is guaranteed to be unique for every
// integer up to the value 68,719,476,736.
func hashInt(id uint64) string {
	var alphabet string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	// Align 1-based to 0-based alphabet index.
	id--

	var buffer bytes.Buffer

	for i := 0; i < 6; i++ {
		buffer.WriteString(string(alphabet[id%62]))
		id /= 62
	}

	return buffer.String()
}
