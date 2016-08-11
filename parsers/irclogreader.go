package parsers

import (
	"regexp"
	"os"
	"bufio"
	"https://github.com/carbontwelve/go-irc-stats/data"
)

//
// This IRC Log Reader parses through each line of the input file and pushes it to Database
//
type IrcLogReader struct {
	RegexAction       *regexp.Regexp
	RegexMessage      *regexp.Regexp
	RegexParseAction  *regexp.Regexp
	RegexParseMessage *regexp.Regexp
}

func (lr *IrcLogReader) Load(filename string, db *data.Database) (err error) {
	// Open the input filename
	file, err := os.Open(filename)
	if (err != nil) {
		return
	}

	// Create Buffered Scanner
	scanner := bufio.NewScanner(file)


}

func (lr *IrcLogReader) parseLine(line string, action bool) {

}