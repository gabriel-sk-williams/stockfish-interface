package main

import (
	"encoding/json"
	"fmt"
	"os"
	"stockfish/model"
)

var (
	// txtInput   = "positions/txt/960positions.txt"
	jsonOutput = "positions/json/960positions.json"
)

func main() {

	// Create json with FEN positions
	// positions.ParseFEN(txtInput, jsonOutput)

	// Load positions
	fischerRandomPositions, err := loadPositions(jsonOutput)
	check(err)

	// Initialize Stockfish
	engine, err := NewStockfishEngine()
	if err != nil {
		fmt.Println("Error initializing Stockfish:", err)
		return
	}
	defer engine.Close()

	for index, position := range fischerRandomPositions {
		if index < 3 {
			fmt.Println(position.FEN)
			// Get best move
			bestMove, score, err := engine.getBestMove(position.FEN, 20) // Analyze to depth 20
			if err != nil {
				fmt.Println("Error getting best move:", err)
				return
			}
			fmt.Printf("Best move: %s (Score: %.2f)\n", bestMove, score)
		}
	}

}

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

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
