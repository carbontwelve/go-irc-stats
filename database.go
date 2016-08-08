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
func (this *Database) Load(path string) (err error)  {
    this.Users = make(map[string]User)

    fh, err := os.Open(path)
    if err != nil {
        return err
    }

    dec := gob.NewDecoder(fh)

    err = dec.Decode(this)
    if err != nil {
        return err
    }

    return
}

func (this Database) Save(path string) (err error)  {
    b := new(bytes.Buffer)
    enc := gob.NewEncoder(b)

    err = enc.Encode(this)
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

func (this *Database) AddUser(u User) {
    if this.HasUser(u.Username) == true {
        panic("Adding a user that already exists in database")
    }
    this.SetUser(u.Username, u)
}

func (this Database) HasUser(nick string) bool {
    if _, ok := this.Users[nick]; ok {
        return true
    }
    return false
}

func (this *Database) SetUser(nick string, u User) {
    this.Users[nick] = u
}

func (this Database) GetUser(nick string) (user User, err error) {
    if this.HasUser(nick) == false {
        err = errors.New("User doesn't exist")
        return
    }
    user = this.Users[nick]
    return
}
