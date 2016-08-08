package main

import (
    "fmt"
    "encoding/gob"
    "errors"
    "os"
    "bytes"
)

type Channel struct{
    UserCount uint
    LineCount uint
    WordCount uint
    MaxDay uint
    Mean int64
    First int64
    Last int64
}

type Database struct{
    Channel Channel
    Users map[string]User
    LastGenerated int64
    ActiveUsers []string
}

/**
 * Load Database from binary file
 * 
 * @param path string
 * @return error|nil
 */
func (d *Database) Load(path string) (err error)  {
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
