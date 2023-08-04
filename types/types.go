package types

type Outcome int

const (
	Checkmate Outcome = iota
	Resignation
	Timeout
	Stalemate
	Draw
	Unknown
)

type Game struct {
	White   string
	Black   string
	Winner  *string
	Outcome Outcome
	Moves   []string
}

type ChessApi interface {
	GetGames(user string) ([]Game, error)
}

type Profile struct {
	User                string
	OpeningsWhite       []string
	OpeningsBlack       []string
	MateWinRatio        float64
	ResignLossRatio     float64
	DurationPercentiles map[int]int
}
