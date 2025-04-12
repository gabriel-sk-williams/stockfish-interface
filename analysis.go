package main

const (
	a = iota
	b
	c
	d
	e
	f
	g
	h
)

var (
	fileMap = map[string]int{
		"a": a,
		"b": b,
		"c": c,
		"d": d,
		"e": e,
		"f": f,
		"g": g,
		"h": h,
	}
)

func getStats(move string, layout string) (string, int, string, string, string) {

	startingFile := string(move[0])
	startingRank := string(move[1])
	endingRank := string(move[3])

	var code string // p1, p2, knight
	if startingRank == "1" {
		code = "knight"
	} else {
		if endingRank == "3" {
			code = "p1"
		} else {
			code = "p2"
		}
	}

	var dtc int // distance to center -> 0, 1, 2, 3 <- edge
	if startingFile == "d" || startingFile == "e" {
		dtc = 0
	}
	if startingFile == "c" || startingFile == "f" {
		dtc = 1
	}
	if startingFile == "b" || startingFile == "g" {
		dtc = 2
	}
	if startingFile == "a" || startingFile == "h" {
		dtc = 3
	}

	var base string // two or three letters
	integrus := fileMap[startingFile]
	backingPiece := string(layout[integrus])
	start, end := getBase(integrus)
	base = layout[start:end]

	// fmt.Println(base, layout, move)

	return code, dtc, base, startingFile, backingPiece
}

func getBase(integrus int) (int, int) {

	// return a-b corner pair
	if integrus == 0 {
		return 0, 2
	}

	// return g-h corner pair
	if integrus == 7 {
		return 6, 8
	}

	// return triplet centered on file
	return integrus - 1, integrus + 2
}
