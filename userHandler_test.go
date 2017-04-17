package main

import (
	mgo "gopkg.in/mgo.v2"
)

type FakeUserRepository struct{}

func (f FakeUserRepository) insertUser(db *mgo.Database, name string, age int64) error {
	return nil
}

func (f FakeUserRepository) getAllUsers(db *mgo.Database) ([]user, error) {
	return nil, nil
}

func (f FakeUserRepository) getUserByID(db *mgo.Database, id string) (user, error) {
	var u user
	return u, nil
}
