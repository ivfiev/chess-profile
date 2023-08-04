package apis

import (
	. "chess-profile/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Lichess struct {
	maxGames   int
	httpClient *http.Client
}

func NewLichess(client *http.Client) *Lichess {
	return &Lichess{20, client}
}

func (li *Lichess) GetGames(user string) ([]Game, error) {
	url := fmt.Sprintf("https://lichess.org/api/games/user/%s?max=%d&rated=true&perfType=rapid,classical", user, li.maxGames)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/x-ndjson")
	res, err := li.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	games := make([]Game, li.maxGames)
	decoder := json.NewDecoder(strings.NewReader(string(resp)))
	for i := 0; decoder.More(); i++ {
		var tmpGame game
		err = decoder.Decode(&tmpGame)
		if err != nil {
			return nil, err
		}
		games[i] = convert(tmpGame)
	}
	return games, nil
}

func convert(g game) Game {
	return Game{
		White:   g.Players.White.User.Name,
		Black:   g.Players.Black.User.Name,
		Winner:  getWinner(g),
		Outcome: getOutcome(g),
		Moves:   getMoves(g),
	}
}

func getOutcome(g game) Outcome {
	switch g.Status {
	case "mate":
		return Checkmate
	case "resign":
		return Resignation
	case "outoftime":
	case "timeout":
		return Timeout
	case "stalemate":
		return Stalemate
	case "draw":
		return Draw
	default:
		return Unknown
	}
	return Unknown
}

func getWinner(g game) *string {
	switch g.Winner {
	case "black":
		return &g.Players.Black.User.Name
	case "white":
		return &g.Players.White.User.Name
	default:
		return nil
	}
}

func getMoves(g game) []string {
	return strings.Split(g.Moves, " ")
}

type game struct {
	Status  string  `json:status`
	Players players `json:players`
	Winner  string  `json:winner`
	Moves   string  `json:moves`
}

type players struct {
	White player `json:white`
	Black player `json:black`
}

type player struct {
	User struct {
		Name string `json:name`
	} `json:user`
}
