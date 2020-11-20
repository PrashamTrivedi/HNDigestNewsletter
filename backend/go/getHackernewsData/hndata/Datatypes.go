package hndata

//Stories represent stories list from firebase api
type Stories []string

//APIResponse Response of API
type APIResponse struct {
	errorData error
	data      interface{}
}

//Story represent a story data
type Story struct {
	URL        string
	Title      string
	By         string
	Text       string
	Time       int64
	TimeString string
	ID         int
	Kids       []int
	KidData    []Story
}
