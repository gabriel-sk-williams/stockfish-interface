package main

import (
	"encoding/json"
	"os"
	"stockfish/model"
)

func loadPositions(filename string) ([]model.ChessPosition, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var positions []model.ChessPosition
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

func loadBestMoveEvals(filename string) ([]model.ChessPositionBestMove, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var positions []model.ChessPositionBestMove
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

func loadTopMovesEvals(filename string) ([]model.ChessPositionTopMoves, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var positions []model.ChessPositionTopMoves
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}
