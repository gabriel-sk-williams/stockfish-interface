package main

import (
	"fmt"
	"strings"
)

// FEN manipulation

// FEN Fields
const (
	PiecePlacement = iota
	SideToMove
	CastlingAbility
	EnPassantTargetSquare
	HalfmoveClock
	FullmoveCounter
)

const (
	NormalStartingPosition = `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`
	MateinOneForWhite      = `r1bqkbnr/ppp2ppp/2np4/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 4`
	MateinOneForBlack      = `rnbqkbnr/pppp1ppp/8/4p3/6P1/5P2/PPPPP2P/RNBQKBNR b KQkq - 0 2`
	BlackWinning           = `rnbqkb1r/pppp1ppp/8/4p3/4P1n1/3P4/PPP2PPP/RNB1KBNR w KQkq - 0 4`
	ActiveWhiteLosing      = `8/8/8/8/6q1/5k2/8/7K w - - 8 5`
	ActiveBlackLosing      = `8/8/8/8/6Q1/5K2/8/7k b - - 0 1`
)

func getField(fen string, index int) (string, error) {
	fields := strings.Fields(fen)

	if len(fields) > index {
		return fields[index], nil
	} else {
		return "", fmt.Errorf("failed to get field")
	}
}
