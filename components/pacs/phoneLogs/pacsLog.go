package phoneLogs

type PACSLogResponse struct {
	Records [][]interface{}
	Last    bool
}

type PACSLog struct {
	Number uint
	Date   string
	Phone  string
}
