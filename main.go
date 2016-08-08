package main

import (
    "fmt"
    "errors"
    "bufio"
    "regexp"
    "strings"
    "time"
    "os"
)

type LogReader struct{
    Database Database
    Channel Channel
    Users map[string]User
    LastGenerated int64
    RegexAction *regexp.Regexp
    RegexMessage *regexp.Regexp
    RegexParseAction *regexp.Regexp
    RegexParseMessage *regexp.Regexp
}

func (lr *LogReader) AddUser(u User) {
    if lr.HasUser(u.Username) == true {
        panic("Adding a user that already exists in database")
    }
    lr.SetUser(u.Username, u)
}

func (lr LogReader) HasUser(nick string) bool {
    if _, ok := lr.Users[nick]; ok {        
        return true
    }
    return false
}

func (lr *LogReader) SetUser(nick string, u User) {
    lr.Users[nick] = u
}

func (lr LogReader) GetUser(nick string) (user User, err error) {
    if lr.HasUser(nick) == false {
        err = errors.New("User doesn't exist")
        return
    }
    user = lr.Users[nick] 
    return
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
            fmt.Printf("error reading line [%i]\n", lr.Channel.LineCount)
        }
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
    return true
}

func (lr *LogReader) ParseLine(line string, isAction bool) bool {
    var parsed[][]string
    var user User

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
    lineNick := strings.Trim(strings.ToLower(parsed[0][2]), " ")
    // @todo check against ignore list

    // Parse message
    lineMessage := strings.Trim(parsed[0][3], " ")

    // If this is an empty line lets ignore it
    if lineMessage == "" {
        return false
    }

    // Get user, if not found make a new user
    if lr.HasUser(lineNick) == true {
        user,_ = lr.GetUser(lineNick)
    }else{
        user = NewUser(lineNick, lineTime.Unix())
    }

    lineMessageCharCount := (strings.Count(lineMessage, "") - 1)
    lineMessageWords := strings.Split(strings.ToLower(lineMessage), " ")
    lineMessageWordCount := len(lineMessageWords)

    user.Words = append(user.Words, lineMessageWords...)

    lr.Channel.LineCount++
    user.LineCount++

    user.WordCount += uint(lineMessageWordCount)
    lr.Channel.WordCount += uint(lineMessageWordCount)
    user.CharCount += uint(lineMessageCharCount)
    
    // Increment Days
    if _, ok := user.Days[lineTime.Format("2006-02-01")]; ok {
        user.Days[lineTime.Format("2006-01-02")]++
    }else{
        user.Days[lineTime.Format("2006-01-02")] = 1
    }

    // Incremember Hours
    user.Hours[uint(lineTime.Hour())]++

    // Set first and last seen timestamps
    if user.FirstSeen > lineTime.Unix() {
        user.FirstSeen = lineTime.Unix()
    }

    if user.LastSeen < lineTime.Unix() {
        user.LastSeen = lineTime.Unix()
    }

    if lr.Channel.First == 0 {
        lr.Channel.First = lineTime.Unix()
    }

    if lr.Channel.Last == 0 {
        lr.Channel.Last = lineTime.Unix()
    }

    if lr.Channel.First > lineTime.Unix() {
        lr.Channel.First = lineTime.Unix()
    }

    if lr.Channel.Last < lineTime.Unix() {
        lr.Channel.Last = lineTime.Unix()
    } 
        
    lr.SetUser(lineNick, user)
    return true
}

func main() {
    db := Database{}
    db.Load("db.bin")

    lr := LogReader{
        RegexAction: regexp.MustCompile(`^\[(.+)\] \* (.+)$`),
        RegexMessage: regexp.MustCompile(`^\[(.+)\] <(.+)> (.+)$`),
        RegexParseAction: regexp.MustCompile(`^\[(.+)\] \* (\S+) (.+)$`), 
        RegexParseMessage: regexp.MustCompile(`^\[(.+)\] <(\S+)> (.+)$`),
        Users: make(map[string]User),
	Database: db,
    }

    lr.LoadFile("irctest.log") 

    db.Save("db.bin")

    //u.CalculateTotals()
    fmt.Printf("%v\n", lr)
}
