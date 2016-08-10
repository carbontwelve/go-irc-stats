package main

type ViewData struct {
	PageTitle string
	PageDescription string
	Channel Channel
	ActiveUsers map[string]User
}

func (d ViewData) TotalDays() int {
	return 123
}