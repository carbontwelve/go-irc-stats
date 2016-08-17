package ircstats

import (
	"regexp"
	"os"
	"bufio"
	"fmt"
	"time"
	"strings"
)

//
// This IRC Log Reader parses through each line of the input file and pushes it to Database
//
type IrcLogReader struct {
	RegexAction       *regexp.Regexp
	RegexMessage      *regexp.Regexp
	RegexParseAction  *regexp.Regexp
	RegexParseMessage *regexp.Regexp
	linesParsed       int64				// Total number of lines parsed
	newLinesFound     int64				// Number of new lines found in log that were parsed (not to be confused with \n)
	NickNameHashTable map[string]string		// Nickname hash table from configuration
	Profiles          map[string]map[string]string  // Profile hash table from configuration
	Ignore            []string			// Ignore list from configuration
}

func NewIrcLogReader(c Config) IrcLogReader {
	return IrcLogReader{
		RegexAction: regexp.MustCompile(`^\[(.+)\] \* (.+)$`),
		RegexMessage: regexp.MustCompile(`^\[(.+)\] <(.+)> (.+)$`),
		RegexParseAction: regexp.MustCompile(`^\[(.+)\] \* (\S+) (.+)$`),
		RegexParseMessage: regexp.MustCompile(`^\[(.+)\] <(\S+)> (.+)$`),
		NickNameHashTable: c.NickNameHashTable,
		Profiles: c.Profiles,
		Ignore: c.Ignore,
	}
}

func (lr IrcLogReader) NewLinesFound() int64 {
	return lr.newLinesFound
}

func (lr *IrcLogReader) Load(filename string, db *Database) (err error) {
	// Open the input filename
	file, err := os.Open(filename)
	if (err != nil) {
		return
	}

	// Create Buffered Scanner
	scanner := bufio.NewScanner(file)

	// Loop over all lines and parse them
	for scanner.Scan() {
		lr.linesParsed++
		line := scanner.Text()
		if lr.RegexAction.MatchString(line) == true {
			lr.parseLine(line, true, db)
		} else if lr.RegexMessage.MatchString(line) == true {
			lr.parseLine(line, false, db)
		} else {
			fmt.Printf("error reading line [%d]\n", lr.linesParsed)
		}
	}

	db.LastGenerated = time.Now().Unix()
	return
}

func (lr *IrcLogReader) parseLine(line string, isAction bool, db *Database) {
	var (
		parsed[][]string
		user User
	)

	// timestamp = [0][1]
	// nick = [0][2]
	// message/action = [0][3]
	if isAction == true {
		parsed = lr.RegexParseAction.FindAllStringSubmatch(line, -1)
	} else {
		parsed = lr.RegexParseMessage.FindAllStringSubmatch(line, -1)
	}

	// Convert timestamp into unix timestamp and if this is a line that
	// we have already parsed do nothing and return
	lineTime := lr.ParseTime(parsed[0][1])
	if lineTime.Unix() <= db.Channel.LastSeen {
		return
	}

	// Parse nick and check against ignore list
	lineNick := strings.Trim(strings.ToLower(parsed[0][2]), " ")
	if lr.IsUserIgnored(lineNick) == true {
		return
	}

	// Map the nickname to one set in configuration
	lineNick = lr.MapNick(lineNick)

	// Parse message
	lineMessage := strings.Trim(parsed[0][3], " ")

	// If this is an empty line lets ignore it
	if lineMessage == "" {
		return
	}

	// Get user, if not found make a new user
	if db.HasUser(lineNick) == true {
		user, _ = db.GetUser(lineNick)
	} else {
		user = NewUser(lineNick, lineTime.Unix())
	}

	lineMessageCharCount := int64(strings.Count(lineMessage, "") - 1)
	lineMessageWords := strings.Split(strings.ToLower(lineMessage), " ")
	lineMessageWordCount := int64(len(lineMessageWords))

	// Append to user words array (@todo this should be unique words right, maybe a map of words with a count of use?)
	user.Words = append(user.Words, lineMessageWords...)

	// Increment Line Counters
	db.Channel.LineCount++
	user.LineCount++

	// Increment Word Counters
	user.WordCount += lineMessageWordCount
	db.Channel.WordCount += lineMessageWordCount

	// Increment Character Counters
	user.CharCount += lineMessageCharCount

	// Increment words per day
	user.IncrementDay(lineTime.Format("2006-02-01"), lineMessageCharCount)

	// Increment lines per day
	db.Channel.IncrementDay(lineTime.Format("2006-02-01"), 1)

	// Increment lines per hour
	user.IncrementHour(uint(lineTime.Hour()))
	db.Channel.IncrementHour(uint(lineTime.Hour()))

	// Update First & Last Timestamps for User and Channel
	user.UpdateSeen(lineTime.Unix())
	db.Channel.UpdateSeen(lineTime.Unix())

	// Store User into DB
	db.SetUser(lineNick, user)

	lr.newLinesFound++;
}

//
// Parse the input time into a time.Time
//
func (lr IrcLogReader) ParseTime(inputDate string) time.Time {
	lineTime, err := time.Parse("2006-01-02 15:04:05 -0700", inputDate)
	if err != nil {
		panic(err)
	}
	return lineTime
}

//
// Map the users nick if found in the configuration against a mapping.
//
func (lr IrcLogReader) MapNick(username string) string {
	if val, ok := lr.NickNameHashTable[username]; ok{
		return val
	}
	return username
}

//
// Identify if the user is to be ignored
//
func (lr IrcLogReader) IsUserIgnored(username string) bool {
	for _, n := range(lr.Ignore) {
		if username == n {
			return true
		}
	}
	return false
}