package main

import (
    "fmt"
    "bufio"
    "regexp"
    "strings"
    "time"
    "os"
)

type LogReader struct{
    Users []User
    TotalLines uint
    LastGenerated int64
    FirstSeen string
    LastSeen string
    RegexAction *regexp.Regexp
    RegexMessage *regexp.Regexp
    RegexParseAction *regexp.Regexp
    RegexParseMessage *regexp.Regexp
}

func (lr *LogReader) LoadFile(path string) bool {
    
    // Open the file
    f, err := os.Open(path)

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
            fmt.Printf("error reading line [%i]\n", lr.TotalLines)
        }
        lr.TotalLines++
    }
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
    // ...
}

func (lr *LogReader) ParseLine(line string, isAction bool) bool {

    var parsed[][]string

    fmt.Println("Parsing Line: [" + line + "]\n\n")
    // timestamp = [0][1]
    // nick = [0][2]
    // message/action = [0][3]
    if isAction == true {
        parsed = lr.RegexParseAction.FindAllStringSubmatch(line, -1)
    }else{
        parsed = lr.RegexParseMessage.FindAllStringSubmatch(line, -1)
    }

    // Convert timestamp into unix timestamp
    lineTime := lr.ParseTime(parsed[0][1])
    if lineTime.Unix() < lr.LastGenerated {
        return false
    }

    // Parse nick and check against ignore list
    lineNick := strings.Trim(parsed[0][2], " ")
    // @todo check against ignore list

    return true
}

func main() {
    lr := LogReader{
        RegexAction: regexp.MustCompile(`^\[(.+)\] \* (.+)$`),
        RegexMessage: regexp.MustCompile(`^\[(.+)\] <(.+)> (.+)$`),
        RegexParseAction: regexp.MustCompile(`^\[(.+)\] \* (\S+) (.+)$`), 
        RegexParseMessage: regexp.MustCompile(`^\[(.+)\] <(\S+)> (.+)$`),
    }

    lr.LoadFile("irctest.log") 

    u := NewUser("Simon", "Today")
    //fmt.Printf("%v\n", lr)

    u.CalculateTotals()
    //fmt.Printf("%v\n", u)
}
