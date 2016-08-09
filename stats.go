package main

type Stats struct {
	Hours [23]uint        // 24 hours
	Days  map[string]uint // total days seen
}

// Increment Hours
func (s *Stats) IncrementHour(hour uint) {
	s.Hours[hour]++
}

// Increment Days
func (s *Stats) IncrementDay(date string, increment uint) {
	if _, ok := s.Days[date]; ok {
		s.Days[date] += increment
	} else {
		s.Days[date] = 1
	}
}

func (s Stats) FindPeakDay() (date string, total uint) {
	for d, t := range (s.Days) {
		if (t > total) {
			date = d
			total = t
		}
	}
	return
}

func (s *Stats) InitiateStats() {
	for i := 0; i < 23; i++ {
		s.Hours[i] = 0
	}
	s.Days = make(map[string]uint)
}
