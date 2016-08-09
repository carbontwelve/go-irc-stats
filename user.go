package main

type User struct {
    Username string
    Url string
    Avatar string
    LineCount uint
    WordCount uint
    CharCount uint
    WordsLine uint
    LineLength uint
    LinesDay uint
    WordsDay uint
    Vocabulary uint
    DaysTotal uint
    FirstSeen int64
    LastSeen int64
    MaxHours uint
    Words []string
    Stats
}

func (u *User) CalculateTotals() {
    u.Vocabulary = uint(len(u.Words))
    u.DaysTotal  = uint(len(u.Days))

    if (u.LineCount > 0) {
        u.WordsLine = uint(u.WordCount / u.LineCount)
    }else{
        u.WordsLine = 0
    }

    // @todo finish
}

func NewUser(nick string, timestamp int64) User {
    var (
        hours [23]uint
        i uint
        days map[string]uint
    )

    for i = 0; i < 23; i++ {
        hours[i] = 0
    }

    days = make(map[string]uint)

    u := User{Username: nick, FirstSeen: timestamp, LastSeen: timestamp, Stats: Stats{Hours: hours, Days: days}}
    return u
}

