package main

type Stats struct {
    Hours [23]uint          // 24 hours
    Days  map[string]uint   // total days seen
}

// Increment Hours
func (s *Stats) IncrementHour (hour uint) {
    s.Hours[hour]++    
}

// Increment Days
func (s *Stats) IncrementDay (date string) {
    if _, ok := s.Days[date]; ok {
        s.Days[date]++
    }else{
        s.Days[date] = 1
    }
}
