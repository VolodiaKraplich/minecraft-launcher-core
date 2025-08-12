package client

type Handshake struct {
	V        string `json:"v"`
	ClientId string `json:"client_id"`
}

type Frame struct {
	Cmd   string `json:"cmd"`
	Args  Args   `json:"args"`
	Nonce string `json:"nonce"`
}

type Args struct {
	Pid      int              `json:"pid"`
	Activity *PayloadActivity `json:"activity"`
}

// ActivityType namespace for activity types
var ActivityType = struct {
	Playing   int
	Listening int
	Watching  int
	Custom    int
	Competing int
}{
	Playing:   0,
	Listening: 2,
	Watching:  3,
	Custom:    4,
	Competing: 5,
}

type PayloadActivity struct {
	Type       int                `json:"type,omitempty"`
	Details    string             `json:"details,omitempty"`
	State      string             `json:"state,omitempty"`
	Assets     PayloadAssets      `json:"assets,omitempty"`
	Party      *PayloadParty      `json:"party,omitempty"`
	Timestamps *PayloadTimestamps `json:"timestamps,omitempty"`
	Secrets    *PayloadSecrets    `json:"secrets,omitempty"`
	Buttons    []*PayloadButton   `json:"buttons,omitempty"`
}

type PayloadAssets struct {
	LargeImage string `json:"large_image,omitempty"`
	LargeText  string `json:"large_text,omitempty"`
	SmallImage string `json:"small_image,omitempty"`
	SmallText  string `json:"small_text,omitempty"`
}

type PayloadParty struct {
	ID   string `json:"id,omitempty"`
	Size [2]int `json:"size,omitempty"`
}

type PayloadTimestamps struct {
	Start *uint64 `json:"start,omitempty"`
	End   *uint64 `json:"end,omitempty"`
}

type PayloadSecrets struct {
	Match    string `json:"match,omitempty"`
	Join     string `json:"join,omitempty"`
	Spectate string `json:"spectate,omitempty"`
}

type PayloadButton struct {
	Label string `json:"label,omitempty"`
	Url   string `json:"url,omitempty"`
}
