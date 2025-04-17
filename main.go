package main

import (
	"fmt"
	"stockfish/model"
	"stockfish/positions"
)

var (
	txtInput       = "positions/txt/960positions.txt"
	jsonOutput     = "positions/json/960positions.json"
	topMovesOutput = "output/topMovesDepth20.json"
	bestMoveOutput = "output/bestMoveDepth42.json"
)

func main() {

	// evalBestMove()
	// evalTopMoves()
	// AnalyzeTopMoves(topMovesOutput)
	AnalyzeBestMove(bestMoveOutput)
}

func evalBestMove() {
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

	var positionsWithEval []model.ChessPositionBestMove

	for index, position := range fischerRandomPositions {
		//if index < 3 {
		bestMove, eval, line, err := engine.getBestMoveAdvanced(position.FEN)
		check(err)

		cpbm := model.ChessPositionBestMove{
			ChessPosition: position,
			BestMove:      bestMove,
			Eval:          eval,
			Line:          line,
		}

		positionsWithEval = append(positionsWithEval, cpbm)

		fmt.Println(index, position.FEN[:8])
		fmt.Printf("Best move: %s Eval: %.2f\n\n", bestMove, eval)
		//}
	}

	positions.WriteOutput(bestMoveOutput, positionsWithEval)
}

func evalTopMoves() {
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

	var positionsWithEval []model.ChessPositionTopMoves

	for index, position := range fischerRandomPositions {
		//if index < 3 {
		topMoves, evals, err := engine.getTopMoves(position.FEN, 2)
		check(err)

		cptm := model.ChessPositionTopMoves{
			ChessPosition: position,
			TopMoves:      topMoves,
			Evals:         evals,
		}

		positionsWithEval = append(positionsWithEval, cptm)

		fmt.Println(index, position.FEN[:8])
		fmt.Printf("Top moves: %s Evals: %.2f\n\n", topMoves, evals)
		//}
	}

	positions.WriteOutput(topMovesOutput, positionsWithEval)
}

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
