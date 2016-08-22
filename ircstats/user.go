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
	MaxHours   int64
	Words      map[string]int64
	HoursAndDaysStats
	Seen
}

func (u *User) CalculateTotals() {
	u.Vocabulary = int64(len(u.Words))
	u.DaysTotal = int64(len(u.Days))

	if u.LineCount > 0 {
		u.WordsLine = int64(u.WordCount / u.LineCount)
	} else {
		u.WordsLine = 0
	}

	// @todo finish
}

func (u *User) AddWord(word string) {
	if u.HasWord(word) == true {
		u.Words[word]++
	} else {
		u.Words[word] = 1;
	}
}

// Check to see if User contains word
func (u User) HasWord(word string) bool {
	if _, ok := u.Words[word]; ok {
		return true
	}
	return false
}

func NewUser(nick string, timestamp int64) *User {
	u := User{Username: nick}
	u.Initiate()
	u.Words = make(map[string]int64)
	u.UpdateSeen(timestamp)
	return &u
}
