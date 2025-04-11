package positions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"stockfish/model"
	"strconv"
	"strings"
	"text/scanner"
)

type parser struct {
	s         *scanner.Scanner
	index     int64
	positions []model.ChessPosition
}

func ParseFEN(txtInput string, jsonOutput string) {
	file, err := os.Open(txtInput)
	check(err)

	stat, err := file.Stat()
	check(err)

	fileName := stat.Name()

	r := bufio.NewReader(file)

	var s scanner.Scanner
	p := parser{s: s.Init(r)}
	p.parse(fileName)

	WriteOutput(jsonOutput, p.positions)
}

func WriteOutput(jsonOutput string, object any) {
	// Write positions to JSON file
	outputFile, err := os.Create(jsonOutput)
	check(err)
	defer outputFile.Close()

	// Create JSON encoder with pretty printing
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")

	// Encode and write positions
	if err := encoder.Encode(object); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Printf("Successfully processed %s positions\n", jsonOutput)
}

// parse walks the document and renders elements to a cell buffer document.
func (p *parser) parse(fileName string) (err error) {

	p.s.Filename = fileName

	for token := p.s.Scan(); token != scanner.EOF; token = p.s.Scan() {

		switch token {
		case 38: // & ampersand
		case 40: // ( paren open
		case 41: // ) paren close
		case 44: // , comma
		case 45: // - hyphen
		case 46: // . period
		case 47: // / slash
		case 59: // ; semicolon
		case 63: // ? question mark
		case 91: // [ bracket open
		case 93: // [ bracket close
		case 94: // ^ carat (dagger)
		case -2: // word
			sequence := p.s.TokenText()
			fen := createFEN(sequence)
			p.positions = append(p.positions, model.ChessPosition{
				PositionNumber: p.index,
				FEN:            fen,
			})
		case -3: // number
			index := p.s.TokenText()
			positionNumber, err := strconv.ParseInt(index, 10, 64)
			check(err)
			p.index = positionNumber
		//case -3: // large number eg 141
		//case -4: // chapter-period eg 4.
		//case -5: // apostrophe ERROR currently escapes
		default:
			fmt.Println("not recognized: ", token)
		}
	}

	return nil
}

func createFEN(upper string) string {
	// FEN Format:
	// bbqnnrkr/pppppppp/8/8/8/8/PPPPPPPP/BBQNNRKR w KQkq - 0 1
	lower := strings.ToLower(upper)
	sequence := fmt.Sprintf("%s/pppppppp/8/8/8/8/PPPPPPPP/%s w KQkq - 0 1", lower, upper)
	return sequence
}

func check(err error) {
	if err != nil {
		fmt.Println("error:", err)
	}
}
