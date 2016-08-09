package main

import (
	"fmt"
	"encoding/gob"
	"errors"
	"os"
	"bytes"
	"time"
)

type MaxDay struct {
	Day   string
	Lines uint
}

type Channel struct {
	UserCount uint
	LineCount uint
	WordCount uint
	MaxDay    MaxDay
	Mean      float64
	First     int64
	Last      int64
}

type Database struct {
	Channel       Channel
	Users         map[string]User
	LastGenerated int64
	ActiveUsers   map[string]User
	Stats
}

/**
 * Load Database from binary file
 * 
 * @param path string
 * @return error|nil
 */
func (d *Database) Load(path string) (err error) {
	d.InitiateStats()
	d.Users = make(map[string]User)

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

func (d Database) Save(path string) (err error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)

	err = enc.Encode(d)
	if err != nil {
		return err
	}

	fh, eopen := os.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0666)
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

func (d *Database) AddUser(u User) {
	if d.HasUser(u.Username) == true {
		panic("Adding a user that already exists in database")
	}
	d.SetUser(u.Username, u)
}

func (d Database) HasUser(nick string) bool {
	if _, ok := d.Users[nick]; ok {
		return true
	}
	return false
}

func (d *Database) SetUser(nick string, u User) {
	d.Users[nick] = u
}

func (d Database) GetUser(nick string) (user User, err error) {
	if d.HasUser(nick) == false {
		err = errors.New("User doesn't exist")
		return
	}
	user = d.Users[nick]
	return
}

func (d *Database) Calculate() {
	// Set Channel User counter
	d.Channel.UserCount = uint(len(d.Users))

	// Get Average lines / day
	d.calculateDailyMeanLines()

	// Get Peak Activity
	d.calculatePeakActivity()

	// Get active Users
	d.calculateActiveUsers()
}

func (d *Database) calculateDailyMeanLines() {
	var (
		sum uint
		size uint
	)

	for _, u := range (d.Users) {
		sum++
		size += u.LineCount
	}
	d.Channel.Mean = float64(sum) / float64(size)
}

func (d *Database) calculatePeakActivity() {
	d.Channel.MaxDay.Day, d.Channel.MaxDay.Lines = d.FindPeakDay()
}

func (d *Database) calculateActiveUsers() {
	timePeriod := make(map[string]bool)
	d.ActiveUsers = make(map[string]User)

	for i := 1; i < 30; i++ {
		timePeriod[time.Now().AddDate(0, 0, -i).Format("2006-02-01")] = true
	}

	for _, u := range (d.Users) {
		var (
			wordCount uint
			daysActive uint
		)

		// Check to see if user has been active within our time period (default past 30 days)
		for timePeriodDate := range (timePeriod) {
			if _, ok := u.Days[timePeriodDate]; ok {
				daysActive++
				wordCount += u.Days[timePeriodDate]
			}
		}

		// If the user is active, copy the struct and push it to the active users hash
		if daysActive > 0 {
			uc := &u
			uc.WordCount = wordCount
			uc.DaysActive = daysActive
			uc.WordsDay = wordCount / daysActive
			d.ActiveUsers[uc.Username] = *uc
		}
	}
}
