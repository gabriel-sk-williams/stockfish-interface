package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"stockfish/model"
	"stockfish/positions"
)

var (
	txtInput   = "positions/txt/960positions.txt"
	jsonOutput = "positions/json/960positions.json"
	evalOutput = "output/evalDepth20.json"
)

type analysis struct {
	leftEdgeCases  map[string]int
	centerCases    map[string]int
	rightEdgeCases map[string]int
}

type fileAnalysis struct {
	bFile map[string]int
	cFile map[string]int
	dFile map[string]int
	eFile map[string]int
	fFile map[string]int
	gFile map[string]int
}

func main() {

	evalPositions, err := loadPositions(evalOutput)
	check(err)

	p1Moves := analysis{
		leftEdgeCases:  make(map[string]int),
		centerCases:    make(map[string]int),
		rightEdgeCases: make(map[string]int),
	}
	p2Moves := analysis{
		leftEdgeCases:  make(map[string]int),
		centerCases:    make(map[string]int),
		rightEdgeCases: make(map[string]int),
	}
	knightMoves := analysis{
		leftEdgeCases:  make(map[string]int),
		centerCases:    make(map[string]int),
		rightEdgeCases: make(map[string]int),
	}
	p2CenterMoves := fileAnalysis{
		bFile: make(map[string]int),
		cFile: make(map[string]int),
		dFile: make(map[string]int),
		eFile: make(map[string]int),
		fFile: make(map[string]int),
		gFile: make(map[string]int),
	}
	p2CenterPieces := fileAnalysis{
		bFile: make(map[string]int),
		cFile: make(map[string]int),
		dFile: make(map[string]int),
		eFile: make(map[string]int),
		fFile: make(map[string]int),
		gFile: make(map[string]int),
	}

	for _, position := range evalPositions {

		layout := position.FEN[:8]

		for _, move := range position.TopMoves {
			code, dtc, base, sf, bp := getStats(move, layout)

			if code == "p1" {
				findCase(p1Moves, base, sf)
			}
			if code == "p2" {
				findCase(p2Moves, base, sf)
			}
			if code == "knight" {
				findCase(knightMoves, base, sf)
			}

			if code == "p2" && dtc < 3 {
				findFile(p2CenterMoves, base, sf)
				findFile(p2CenterPieces, bp, sf)
			}

			if bp == "n" && (sf == "d" || sf == "e") {
				fmt.Println(base)
			}
		}
	}

	//sortMap(knightMoves.leftEdgeCases)
	//sortMap(knightMoves.centerCases)
	//sortMap(knightMoves.rightEdgeCases)

	//sortMap(p1Moves.leftEdgeCases)
	//sortMap(p1Moves.centerCases)
	//sortMap(p1Moves.rightEdgeCases)

	//sortMap(p2Moves.leftEdgeCases)
	//sortMap(p2Moves.centerCases)
	//sortMap(p2Moves.rightEdgeCases)

	//fmt.Println("b")
	//sortMap(p2CenterMoves.bFile)
	//fmt.Println("c")
	//sortMap(p2CenterMoves.cFile)
	//fmt.Println("d")
	//sortMap(p2CenterMoves.dFile)
	//fmt.Println("e")
	//sortMap(p2CenterMoves.eFile)
	//fmt.Println("f")
	//sortMap(p2CenterMoves.fFile)
	//fmt.Println("g")

	fmt.Println("b")
	sortMap(p2CenterPieces.bFile)
	fmt.Println("c")
	sortMap(p2CenterPieces.cFile)
	fmt.Println("d")
	sortMap(p2CenterPieces.dFile)
	fmt.Println("e")
	sortMap(p2CenterPieces.eFile)
	fmt.Println("f")
	sortMap(p2CenterPieces.fFile)
	fmt.Println("g")
}

func findCase(m analysis, base string, startingFile string) {
	if startingFile == "a" {
		incrementToMap(m.leftEdgeCases, base)
	} else if startingFile == "h" {
		incrementToMap(m.rightEdgeCases, base)
	} else {
		incrementToMap(m.centerCases, base)
	}
}

func findFile(m fileAnalysis, base string, startingFile string) {
	if startingFile == "b" {
		incrementToMap(m.bFile, base)
	} else if startingFile == "c" {
		incrementToMap(m.cFile, base)
	} else if startingFile == "d" {
		incrementToMap(m.dFile, base)
	} else if startingFile == "e" {
		incrementToMap(m.eFile, base)
	} else if startingFile == "f" {
		incrementToMap(m.fFile, base)
	} else if startingFile == "g" {
		incrementToMap(m.gFile, base)
	}
}

func incrementToMap(m map[string]int, key string) {
	_, exists := m[key]
	if exists {
		m[key] = m[key] + 1
	} else {
		m[key] = 1
	}
}

func sortMap(m map[string]int) { // map[string]int {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	fmt.Println(ss)
	fmt.Println()
}

func parsePositions() {
	// Create json with FEN positions
	positions.ParseFEN(txtInput, jsonOutput)
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

	var positionsWithEval []model.ChessPosition

	for index, position := range fischerRandomPositions {
		//if index < 3 {
		bestMoves, eval, err := engine.getTopMoves(position.FEN, 2) // Analyze to depth 20
		check(err)

		position.TopMoves = bestMoves
		position.Eval = eval
		positionsWithEval = append(positionsWithEval, position)

		fmt.Println(index, position.FEN[:8])
		fmt.Printf("Best moves: %s Eval: %.2f\n\n", bestMoves, eval)
		//}
	}

	positions.WriteOutput(evalOutput, positionsWithEval)
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
