package model

import (
	"fmt"
	"sort"
)

//
// Positions + Eval
//

type ChessPosition struct {
	PositionNumber int64  `json:"positionNumber"`
	FEN            string `json:"fen"`
}

type ChessPositionBestMove struct {
	ChessPosition
	BestMove string   `json:"bestMove"`
	Eval     float64  `json:"eval"` // string allows for #3 checkmate
	Line     []string `json:"line"`
}

type ChessPositionTopMoves struct {
	ChessPosition
	TopMoves []string  `json:"topMoves"`
	Evals    []float64 `json:"evals"`
}

//
// Analysis, with loose grouping into edge and center moves
//

type Analysis struct {
	leftEdgeCases  map[string]int
	centerCases    map[string]int
	rightEdgeCases map[string]int
}

func CreateAnalysis() Analysis {
	return Analysis{
		leftEdgeCases:  make(map[string]int),
		centerCases:    make(map[string]int),
		rightEdgeCases: make(map[string]int),
	}
}

func (m Analysis) FindCase(base string, startingFile string) {
	if startingFile == "a" {
		incrementToMap(m.leftEdgeCases, base)
	} else if startingFile == "h" {
		incrementToMap(m.rightEdgeCases, base)
	} else {
		incrementToMap(m.centerCases, base)
	}
}

func (m Analysis) ShowAnalysis() {
	sortMap(m.leftEdgeCases)
	sortMap(m.centerCases)
	sortMap(m.rightEdgeCases)
}

//
// CenterFileAnalysis, for analyzing frequency of non-edge files.
//

type CenterFileAnalysis struct {
	bFile map[string]int
	cFile map[string]int
	dFile map[string]int
	eFile map[string]int
	fFile map[string]int
	gFile map[string]int
}

func CreateCenterFileAnalaysis() CenterFileAnalysis {
	return CenterFileAnalysis{
		bFile: make(map[string]int),
		cFile: make(map[string]int),
		dFile: make(map[string]int),
		eFile: make(map[string]int),
		fFile: make(map[string]int),
		gFile: make(map[string]int),
	}
}

func (m CenterFileAnalysis) ShowCenterFileAnalysis() {
	fmt.Println("b")
	sortMap(m.bFile)
	fmt.Println("c")
	sortMap(m.cFile)
	fmt.Println("d")
	sortMap(m.dFile)
	fmt.Println("e")
	sortMap(m.eFile)
	fmt.Println("f")
	sortMap(m.fFile)
	fmt.Println("g")
}

func (m CenterFileAnalysis) FindFile(base string, startingFile string) {
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

//
// Utility Functions
//

func incrementToMap(m map[string]int, key string) {
	_, exists := m[key]
	if exists {
		m[key] = m[key] + 1
	} else {
		m[key] = 1
	}
}

func sortMap(m map[string]int) {
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
