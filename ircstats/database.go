package ircstats

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
)

//
// Cached data that has been parsed is stored in the Database struct. Once execution is complete is can save to disk
// this saves us from restarting from the beginning of the file.
//
// The primary data that is stored here is parsed user stats
//
type Database struct {
	Users         map[string]User // Users
	LastGenerated int64           // Unix timestamp of last generated
	Version       string          // Version of application that this works with, if the version changes we need to check for incompatibility
	Channel       Channel         // Channel counters for storage
}

//
// Load Database from disk
//
func (d *Database) Load(path string) (err error) {
	d.Users = make(map[string]User)
	d.Channel.Initiate()

	fh, err := os.Open(path)
	if err != nil {
		return err
	}

	dec := gob.NewDecoder(fh)

	err = dec.Decode(d)
	if err != nil {
		return err
	}
	return
}

//
// Save Database to disk
//
func (d Database) Save(path string) (err error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)

	err = enc.Encode(d)
	if err != nil {
		return err
	}

	fh, eopen := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	defer fh.Close()

	if eopen != nil {
		return eopen
	}

	n, e := fh.Write(b.Bytes())
	if e != nil {
		return e
	}
	fmt.Fprintf(os.Stderr, "%d bytes successfully written to file\n", n)
	return nil
}

//
// Add new User
//
func (d *Database) AddUser(u User) {
	if d.HasUser(u.Username) == true {
		panic("Adding a user that already exists in database")
	}
	d.SetUser(u.Username, u)
}

//
// Check to see if Database contains User by nick
//
func (d Database) HasUser(nick string) bool {
	if _, ok := d.Users[nick]; ok {
		return true
	}
	return false
}

//
// Set User by nick
//
func (d *Database) SetUser(nick string, u User) {
	d.Users[nick] = u
}

//
// Get User by nick, returns an error if User not found
//
func (d Database) GetUser(nick string) (user User, err error) {
	if d.HasUser(nick) == false {
		err = errors.New("User doesn't exist")
		return
	}
	user = d.Users[nick]
	return
}

// Return count of Users in data store
func (d Database) CountUsers() int64 {
	return int64(len(d.Users))
}

//
// At some point this needs to check the version
// @todo implement this
//
func (d Database) CheckVersion(version string) (ok bool, err error) {
	ok = (d.Version != version)
	if ok == false {
		err = errors.New("This version of go-irc-stats is incompatible with the current cached datastore.")
	}
	return
}
