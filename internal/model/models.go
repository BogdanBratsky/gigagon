package model

const (
	RoomWaiting = "waiting"
	RoomWriting = "writing"
	RoomVoting  = "voting"
	RoomResults = "results"
)

type Room struct {
	ID     string
	HostID string
	State  string

	Players map[string]*Player

	Proverbs []Proverb
	Index    int

	Rounds []Round
}

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type Proverb struct {
	ID   string
	Text string
}

type Round struct {
	Proverb Proverb

	Answers map[string]Answer
	Votes   map[string]string

	Answered map[string]bool
	Voted    map[string]bool

	Done bool
}

type Answer struct {
	ID       string
	PlayerId string
	Text     string
}
