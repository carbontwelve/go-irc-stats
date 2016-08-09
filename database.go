package main

import (
    "fmt"
    "encoding/gob"
    "errors"
    "os"
    "bytes"
)

type MaxDay struct{
    Day string
    Lines uint
}

type Channel struct{
    UserCount uint
    LineCount uint
    WordCount uint
    MaxDay MaxDay
    Mean float64
    First int64
    Last int64
}

type Database struct{
    Channel Channel
    Users map[string]User
    LastGenerated int64
    ActiveUsers []string
    Stats
}

/**
 * Load Database from binary file
 * 
 * @param path string
 * @return error|nil
 */
func (d *Database) Load(path string) (err error)  {
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

func (d Database) Save(path string) (err error)  {
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

    n,e := fh.Write(b.Bytes())
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
}

func (d *Database) calculateDailyMeanLines() {
    var (
        sum uint
        size uint
    )

    for _, u := range(d.Users) {
        sum++
        size += u.LineCount
    }
    d.Channel.Mean = float64(sum)/float64(size)
    //fmt.Printf("Sum / Size : %d/%d = %f \n", sum, size, d.Channel.Mean)
}

func (d *Database) calculatePeakActivity() {
    d.Channel.MaxDay.Day, d.Channel.MaxDay.Lines = d.FindPeakDay()
}
