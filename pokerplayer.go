package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// VERSION provides a short description of the player's current version
// The string specified here will be shown for live.leanpoker.org games
const VERSION = "Default Go folding player"

// PokerPlayer is a struct to organize player methods
type PokerPlayer struct{}

// NewPokerPlayer creates a new instance of *PokerPlayer
func NewPokerPlayer() *PokerPlayer {
	return &PokerPlayer{}
}

// BetRequest handles the main betting logic. The return value of this
// function will be used to decide whether the player want to fold,
// call, raise or do an all-in; more information about this behaviour
// can be found here: http://leanpoker.org/player-api
func (p *PokerPlayer) BetRequest(state *Game) int {

	u, _ := url.Parse("http://rainman.leanpoker.org/rank")
	q := u.Query()

	j, err := json.Marshal(state.CommunityCards)

	if err != nil {
		return 15
	}

	q.Set("cards", string(j))
	u.RawQuery = q.Encode()

	fmt.Println(u.String())
	res, err := http.Get(u.String())

	var rr = RainmanResponse{}
	err = json.NewDecoder(res.Body).Decode(&rr)

	fmt.Println(fmt.Sprintf("Rank: %d", rr.Rank))
	if err != nil {
		return 15
	}

	return 15
}

// Showdown is called at the end of every round making it possible to
// e.g. collect statistics or log end results of the games
func (p *PokerPlayer) Showdown(state *Game) {

}

// Version returns the version string that will be shown on the UI at
// live.leanpoker.org
func (p *PokerPlayer) Version() string {
	return VERSION
}
