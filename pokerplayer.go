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

const RAINMAN_URL = "http://rainman.leanpoker.org/rank"

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

	fmt.Println(state)
	u, _ := url.Parse(RAINMAN_URL)
	q := u.Query()

	var cards []Card = state.CommunityCards
	var stack int = 1000

	for _, player := range state.Players {
		if player.ID == 2 {
			stack = player.Stack
			cards = append(cards, player.HoleCards...)
		}
	}

	fmt.Println(fmt.Sprintf("Number of Cards: %d", len(cards)))

	if len(cards) == 2 {
		if AnalyseFirstTwoCards(cards[0], cards[1]) {
			return state.CurrentBuyIn
		}
		return ReturnDefaultBet()
	}

	j, err := json.Marshal(cards)

	if err != nil {
		return ReturnDefaultBet()
	}

	q.Set("cards", string(j))
	u.RawQuery = q.Encode()

	fmt.Println(u.String())
	res, err := http.Get(u.String())

	var rr = RainmanResponse{}
	err = json.NewDecoder(res.Body).Decode(&rr)

	if err != nil {
		return ReturnDefaultBet()
	}
	fmt.Println(rr.Rank, len(cards))
	switch rr.Rank {
	case 0:
		if len(cards) >= 6 {
			return 0
		}
		return state.CurrentBuyIn
	case 1:
		return 1 / 5 * stack
	case 2:
		return 1 / 4 * stack
	case 3:
		return 1 / 3 * stack
	case 4:
		return stack
	case 5:
		return stack
	case 6:
		return stack
	case 7:
		return stack
	case 8:
		return stack
	default:
		return ReturnDefaultBet()
	}
}

func AnalyseFirstTwoCards(card1 Card, card2 Card) bool {
	if card1.Rank == card2.Rank {
		return true
	}

	var highCards = []string{"J", "Q", "K", "A"}

	if contains(card1.Rank, highCards) || contains(card2.Rank, highCards) {
		return true
	}
	return false
}

func contains(s string, a []string) bool {
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func ReturnDefaultBet() int {
	return 101
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
