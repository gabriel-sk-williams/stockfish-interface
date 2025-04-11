package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	stockfishLocation = "/mnt/ssd/env/stockfish/src/stockfish"
	depth             = 20
	timeMs            = 8000
	threads           = 14
	RAM               = 8192
)

type StockfishEngine struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
	ready  bool
}

func NewStockfishEngine() (*StockfishEngine, error) {
	cmd := exec.Command(stockfishLocation)
	stdin, err := cmd.StdinPipe()
	check(err)

	stdout, err := cmd.StdoutPipe()
	check(err)

	err = cmd.Start()
	check(err)

	engine := &StockfishEngine{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewScanner(stdout),
		ready:  false,
	}

	// Initialize Universal Chess Interface mode
	engine.sendCommand("uci")

	// Wait for "uciok" response
	for engine.stdout.Scan() {
		text := engine.stdout.Text()
		if strings.Contains(text, "uciok") {
			engine.ready = true
			break
		}
	}

	if !engine.ready {
		return nil, fmt.Errorf("failed to initialize UCI mode")
	}

	// Set Chess960 option
	engine.sendCommand("setoption name UCI_Chess960 value true")

	// Set utilizable Threads and RAM
	engine.sendCommand(fmt.Sprintf("setoption name Threads value %d", threads)) // 14/16 CPU cores
	engine.sendCommand(fmt.Sprintf("setoption name Hash value %d", RAM))        // 8192 -> 8GB RAM

	return engine, nil
}

// sendCommand sends a command to Stockfish
func (engine *StockfishEngine) sendCommand(command string) error {
	_, err := io.WriteString(engine.stdin, command+"\n")
	return err
}

// Close closes the Stockfish engine
func (engine *StockfishEngine) Close() error {
	engine.sendCommand("quit")
	return engine.cmd.Wait()
}

// getBestMove gets the best move for a given position from Stockfish
func (engine *StockfishEngine) getBestMove(fen string) (string, string, error) {
	// Set the position using the FEN string
	engine.sendCommand(fmt.Sprintf("position fen %s", fen))

	// Start the analysis
	//engine.sendCommand(fmt.Sprintf("go depth %d", depth))
	engine.sendCommand(fmt.Sprintf("go depth %d movetime %d", depth, timeMs))

	bestMove := ""
	var eval string
	var score float64

	// Wait for "bestmove" response
	for engine.stdout.Scan() {
		text := engine.stdout.Text()

		// Look for info depth <depth> ... score cp <score> ...
		if strings.Contains(text, fmt.Sprintf("depth %d", depth)) && strings.Contains(text, "score") {
			// Parse the score
			scoreRegex := regexp.MustCompile(`score (cp|mate) (-?\d+)`)
			scoreMatches := scoreRegex.FindStringSubmatch(text)

			if len(scoreMatches) >= 3 {
				if scoreMatches[1] == "cp" { // Convert centipawn score to pawns
					fmt.Sscanf(scoreMatches[2], "%f", &score)
					score = score / 100.0
					fmt.Println(score)
					eval = strconv.FormatFloat(score, 'f', -1, 64)
				} else if scoreMatches[1] == "mate" {
					var mateIn int
					fmt.Sscanf(scoreMatches[2], "%d", &mateIn)
					eval = fmt.Sprintf("#%d", mateIn)
				}
			}
		}

		// Look for bestmove command
		if strings.HasPrefix(text, "bestmove") {
			parts := strings.Fields(text)
			if len(parts) >= 2 {
				bestMove = parts[1]
				break
			}
		}
	}

	if bestMove == "" {
		return "", "", fmt.Errorf("failed to get best move")
	}

	return bestMove, eval, nil
}

// Updated function to return multiple moves
func (engine *StockfishEngine) getTopMoves(fen string, numMoves int) ([]string, []float64, error) {
	// Set MultiPV option
	engine.sendCommand(fmt.Sprintf("setoption name MultiPV value %d", numMoves))

	// Set the position using the FEN string
	engine.sendCommand(fmt.Sprintf("position fen %s", fen))

	// Start the analysis
	// engine.sendCommand(fmt.Sprintf("go depth %d", depth))
	engine.sendCommand(fmt.Sprintf("go depth %d movetime %d", depth, timeMs))

	// Arrays to store top moves and scores
	moves := make([]string, numMoves)
	scores := make([]float64, numMoves)

	// Track which PVs we've updated at the final depth
	pvsFound := 0

	// Wait for "bestmove" response
	for engine.stdout.Scan() {
		text := engine.stdout.Text()

		// Look for info depth <depth> multipv <n> ... pv <move>
		if strings.Contains(text, fmt.Sprintf("depth %d", depth)) && strings.Contains(text, "multipv") {
			// Parse the MultiPV line
			pvRegex := regexp.MustCompile(`multipv (\d+)`)
			pvMatches := pvRegex.FindStringSubmatch(text)

			if len(pvMatches) >= 2 {
				var pvNum int
				fmt.Sscanf(pvMatches[1], "%d", &pvNum)

				// PV numbers are 1-based, convert to 0-based index
				pvIndex := pvNum - 1

				if pvIndex < numMoves {
					// Parse the score
					scoreRegex := regexp.MustCompile(`score (cp|mate) (-?\d+)`)
					scoreMatches := scoreRegex.FindStringSubmatch(text)

					if len(scoreMatches) >= 3 {
						if scoreMatches[1] == "cp" {
							var cp float64
							fmt.Sscanf(scoreMatches[2], "%f", &cp)
							scores[pvIndex] = cp / 100.0 // Convert centipawns to pawns
						} else if scoreMatches[1] == "mate" {
							var mateIn int
							fmt.Sscanf(scoreMatches[2], "%d", &mateIn)
							if mateIn > 0 {
								scores[pvIndex] = 999.0 // Positive mate
							} else {
								scores[pvIndex] = -999.0 // Negative mate
							}
						}
					}

					// Parse the move
					pvRegex := regexp.MustCompile(`pv ([a-h][1-8][a-h][1-8][qrbn]?)`)
					pvMatches := pvRegex.FindStringSubmatch(text)
					if len(pvMatches) >= 2 {
						moves[pvIndex] = pvMatches[1]
						pvsFound++
					}
				}
			}
		}

		// Look for bestmove command
		if strings.HasPrefix(text, "bestmove") {
			// Make sure we have at least the first move
			if moves[0] != "" {
				break
			}

			// If we don't have the move yet, get it from bestmove
			parts := strings.Fields(text)
			if len(parts) >= 2 {
				moves[0] = parts[1]
				break
			}
		}
	}

	if moves[0] == "" {
		return nil, nil, fmt.Errorf("failed to get best move")
	}

	return moves, scores, nil
}

// Add this to your StockfishEngine struct methods
func (engine *StockfishEngine) getFENAfterMove(initialFEN string, move string) (string, error) {
	// Set the initial position
	engine.sendCommand(fmt.Sprintf("position fen %s", initialFEN))

	// Make the move
	engine.sendCommand(fmt.Sprintf("position fen %s moves %s", initialFEN, move))

	// Request the position display
	engine.sendCommand("d")

	// Look for the FEN in the output
	var updatedFEN string
	for engine.stdout.Scan() {
		text := engine.stdout.Text()

		// Look for "Fen: " in the output
		if strings.HasPrefix(text, "Fen:") {
			updatedFEN = strings.TrimSpace(strings.TrimPrefix(text, "Fen:"))
			break
		}
	}

	if updatedFEN == "" {
		return "", fmt.Errorf("failed to get updated FEN")
	}

	return updatedFEN, nil
}
