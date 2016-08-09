package main

import (
	"fmt"
	"bufio"
	"regexp"
	"strings"
	"time"
	"os"
)

type LogReader struct {
	Database          Database
	RegexAction       *regexp.Regexp
	RegexMessage      *regexp.Regexp
	RegexParseAction  *regexp.Regexp
	RegexParseMessage *regexp.Regexp
	Config            Config
}

func (lr *LogReader) LoadFile() bool {
	// Open the file
	f, err := os.Open(lr.Config.Location)

	if (err != nil) {
		return false
	}

	// Create Buffered Scanner
	scanner := bufio.NewScanner(f)

	// Loop over all found lines and analyse them
	for scanner.Scan() {
		line := scanner.Text()
		if lr.RegexAction.MatchString(line) == true {
			lr.ParseLine(line, true)
		} else if lr.RegexMessage.MatchString(line) == true {
			lr.ParseLine(line, false)
		} else {
			fmt.Printf("error reading line [%d]\n", lr.Database.Channel.LineCount)
		}
	}

	lr.Database.LastGenerated = time.Now().Unix()

	return true;
}

func (lr LogReader) ParseTime(inputDate string) time.Time {
	// Convert timestamp into unix timestamp
	// @todo make input format configurable
	lineTime, err := time.Parse("2006-01-02 15:04:05 -0700", inputDate)
	if err != nil {
		panic(err)
	}
	return lineTime
}

func (lr LogReader) IsUserIgnored(username string) bool {
	for _, n := range(lr.Config.Ignore) {
		if username == n {
			return true
		}
	}
	return false
}

func (lr *LogReader) ParseLine(line string, isAction bool) bool {
	var parsed[][]string
	var user User

	// timestamp = [0][1]
	// nick = [0][2]
	// message/action = [0][3]

	if isAction == true {
		parsed = lr.RegexParseAction.FindAllStringSubmatch(line, -1)
	} else {
		parsed = lr.RegexParseMessage.FindAllStringSubmatch(line, -1)
	}

	// Convert timestamp into unix timestamp
	lineTime := lr.ParseTime(parsed[0][1])
	if lineTime.Unix() <= lr.Database.Channel.Last {
		return false
	}

	// Parse nick and check against ignore list
	lineNick := strings.Trim(strings.ToLower(parsed[0][2]), " ")
	if (lr.IsUserIgnored(lineNick) == true) {
		return false
	}

	// Parse message
	lineMessage := strings.Trim(parsed[0][3], " ")

	// If this is an empty line lets ignore it
	if lineMessage == "" {
		return false
	}

	// Get user, if not found make a new user
	if lr.Database.HasUser(lineNick) == true {
		user, _ = lr.Database.GetUser(lineNick)
	} else {
		user = NewUser(lineNick, lineTime.Unix())
	}

	lineMessageCharCount := (strings.Count(lineMessage, "") - 1)
	lineMessageWords := strings.Split(strings.ToLower(lineMessage), " ")
	lineMessageWordCount := len(lineMessageWords)

	user.Words = append(user.Words, lineMessageWords...)

	lr.Database.Channel.LineCount++
	user.LineCount++

	user.WordCount += uint(lineMessageWordCount)
	lr.Database.Channel.WordCount += uint(lineMessageWordCount)
	user.CharCount += uint(lineMessageCharCount)

	// Increment Days
	user.IncrementDay(lineTime.Format("2006-02-01"))
	lr.Database.IncrementDay(lineTime.Format("2006-02-01"))

	// Increment Hours
	user.IncrementHour(uint(lineTime.Hour()))
	lr.Database.IncrementHour(uint(lineTime.Hour()))

	// Set first and last seen timestamps
	if user.FirstSeen > lineTime.Unix() {
		user.FirstSeen = lineTime.Unix()
	}

	if user.LastSeen < lineTime.Unix() {
		user.LastSeen = lineTime.Unix()
	}

	if lr.Database.Channel.First == 0 {
		lr.Database.Channel.First = lineTime.Unix()
	}

	if lr.Database.Channel.Last == 0 {
		lr.Database.Channel.Last = lineTime.Unix()
	}

	if lr.Database.Channel.First > lineTime.Unix() {
		lr.Database.Channel.First = lineTime.Unix()
	}

	if lr.Database.Channel.Last < lineTime.Unix() {
		lr.Database.Channel.Last = lineTime.Unix()
	}

	lr.Database.SetUser(lineNick, user)
	return true
}