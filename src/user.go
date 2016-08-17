package ircstats

type User struct {
	Username   string
	Url        string
	Avatar     string
	LineCount  int64
	WordCount  int64
	DaysActive int64
	CharCount  int64
	WordsLine  int64
	LineLength int64
	LinesDay   int64
	WordsDay   int64
	Vocabulary int64
	DaysTotal  int64
	FirstSeen  int64
	LastSeen   int64
	MaxHours   int64
	Words      []string
	HoursAndDaysStats
}

func (u *User) CalculateTotals() {
	u.Vocabulary = int64(len(u.Words))
	u.DaysTotal = int64(len(u.Days))

	if (u.LineCount > 0) {
		u.WordsLine = int64(u.WordCount / u.LineCount)
	} else {
		u.WordsLine = 0
	}

	// @todo finish
}

func NewUser(nick string, timestamp int64) User {
	u := User{Username: nick, FirstSeen: timestamp, LastSeen: timestamp}
	u.Initiate()
	return u
}

