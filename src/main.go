package main

import (
    "fmt"
    "bufio"
    "regexp"
    "os"
)

type LogReader struct{
    Users []User
    TotalLines uint
    FirstSeen string
    LastSeen string
    RegexAction *regexp.Regexp
    RegexMessage *regexp.Regexp
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
        lr.ParseLine(line)
        lr.TotalLines++
    }
    return true;
}

func (lr *LogReader) ParseLine(line string) {

    fmt.Println("Parsing Line: [" + line + "]")
    fmt.Println("Is action: ", lr.RegexAction.MatchString(line))
    fmt.Println("Is message: ", lr.RegexMessage.MatchString(line))    
    fmt.Println("")
}

func main() {
    lr := LogReader{RegexAction: regexp.MustCompile(`^\[(.+)\] \* (.+)$`), RegexMessage: regexp.MustCompile(`^\[(.+)\] <(.+)> (.+)$`)}
    lr.LoadFile("irctest.log") 

    u := NewUser("Simon", "Today")
    fmt.Printf("%v\n", lr)

    u.CalculateTotals()
    fmt.Printf("%v\n", u)
}
