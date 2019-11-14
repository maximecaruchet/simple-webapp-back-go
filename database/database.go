package database

import (
	"github.com/asdine/storm/v3"
	"log"
)

const DefaultDbName = "simple-app.db"

type Database struct {
	db *storm.DB
}

type User struct {
	ID        int    `storm:"id,increment" json:"id"` // primary key with auto increment
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

func Open(filename string) *Database {
	db, err := storm.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{
		db: db,
	}
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) AddUser(user User) (User, error) {
	err := d.db.Save(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (d *Database) ListUsers() ([]User, error) {
	var users []User

	err := d.db.All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
