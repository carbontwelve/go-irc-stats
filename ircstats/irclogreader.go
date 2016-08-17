package ircstats

import (
	"regexp"
	"os"
	"bufio"
	"fmt"
	"time"
)

//
// This IRC Log Reader parses through each line of the input file and pushes it to Database
//
type IrcLogReader struct {
	RegexAction       *regexp.Regexp
	RegexMessage      *regexp.Regexp
	RegexParseAction  *regexp.Regexp
	RegexParseMessage *regexp.Regexp
	linesParsed       int64
	NickNameHashTable map[string]string		// Nickname hash table from configuration
	Profiles          map[string]map[string]string  // Profile hash table from configuration
	Ignore            []string			// Ignore list from configuration
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

func (lr IrcLogReader) parseLine(line string, isAction bool, db *Database) {

	var (
		parsed[][]string
		//user User
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
	if lineTime.Unix() <= db.Channel.Last {
		return
	}

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