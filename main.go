package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
)

// StockfishEngine represents a connection to Stockfish
type StockfishEngine struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Scanner
	ready  bool
}

// NewStockfishEngine creates and initializes a Stockfish engine
func NewStockfishEngine() (*StockfishEngine, error) {
	cmd := exec.Command("/mnt/ssd/env/stockfish/src/stockfish")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start Stockfish: %v", err)
	}

	engine := &StockfishEngine{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewScanner(stdout),
		ready:  false,
	}

	// Initialize UCI mode
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
func (engine *StockfishEngine) getBestMove(fen string, depth int) (string, float64, error) {
	// Set the position using the FEN string
	engine.sendCommand(fmt.Sprintf("position fen %s", fen))

	// Start the analysis
	engine.sendCommand(fmt.Sprintf("go depth %d", depth))

	bestMove := ""
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
				if scoreMatches[1] == "cp" {
					// Convert centipawn score to pawns
					fmt.Sscanf(scoreMatches[2], "%f", &score)
					score = score / 100.0
				} else if scoreMatches[1] == "mate" {
					// Handle mate scores
					var mateIn int
					fmt.Sscanf(scoreMatches[2], "%d", &mateIn)
					if mateIn > 0 {
						score = 999.0 // Positive mate
					} else {
						score = -999.0 // Negative mate
					}
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
		return "", 0, fmt.Errorf("failed to get best move")
	}

	return bestMove, score, nil
}

func main() {
	// Example Chess960 FEN
	// This is position 518 (RNKQBBNR)
	fen := "rnkqbbnr/pppppppp/8/8/8/8/PPPPPPPP/RNKQBBNR w KQkq - 0 1"

	// Initialize Stockfish
	engine, err := NewStockfishEngine()
	if err != nil {
		fmt.Println("Error initializing Stockfish:", err)
		return
	}
	defer engine.Close()

	// Get best move
	bestMove, score, err := engine.getBestMove(fen, 20) // Analyze to depth 20
	if err != nil {
		fmt.Println("Error getting best move:", err)
		return
	}

	fmt.Printf("Best move: %s (Score: %.2f)\n", bestMove, score)
}
