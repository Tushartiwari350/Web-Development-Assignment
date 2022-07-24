package HandlerFunc

// Structs
type User struct {
	Name string `json:"name"`
	Msgs string `json"msgs"`
}

type Resultstring struct {
	Result string `json:"result"`
}

type FollowInfo struct {
	Username  string `json:"username"`
	Following string `json:"following"`
}

type PrintPost struct {
	Name string   `json:"Username"`
	Post []string `json:"posts"`
}
