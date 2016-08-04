package main

import (
    "fmt"
    "bufio"
    "os"
)

type LogReader struct{
    Users []User
    TotalLines uint
    FirstSeen string
    LastSeen string
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
        fmt.Println(line)
        lr.TotalLines++
    }
    return true;
}

func main() {

   lr := LogReader{}
   lr.LoadFile("irctest.log")

   u := NewUser("Simon", "Today")
   fmt.Printf("%v\n", lr)

   u.CalculateTotals()
   fmt.Printf("%v\n", u)
}
