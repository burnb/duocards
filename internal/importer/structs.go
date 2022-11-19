package importer

//go:generate easyjson -all structs.go

type Card struct {
	Front string
	Back  string
	Hint  string
}

type Request struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type CreateCardVariables struct {
	DeckId      string `json:"deckId"`
	Front       string `json:"front"`
	Back        string `json:"back"`
	LangBack    string `json:"langBack"`
	Hint        string `json:"hint"`
	FromSharing bool   `json:"fromSharing"`
}
