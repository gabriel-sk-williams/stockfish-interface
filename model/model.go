package model

type ChessPosition struct {
	PositionNumber int64   `json:"positionNumber"`
	FEN            string  `json:"fen"`
	Eval           float64 `json:"eval"`
}
