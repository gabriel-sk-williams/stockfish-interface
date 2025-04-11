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
			bestMove, eval, err := engine.getTopMoves(position.FEN, 2) // Analyze to depth 20
			check(err)

			fmt.Println(position.FEN[:8])
			fmt.Printf("Best moves: %s Eval: %.2f\n\n", bestMove, eval)
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
