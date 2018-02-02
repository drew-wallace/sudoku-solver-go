package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/drew-wallace/sudoku-solver-go/sudoku-puzzle"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("./SudokuPuzzle (your input file)")
		os.Exit(0)
	}

	file1, err := ioutil.ReadFile(os.Args[1])
	check(err)
	file1String := string(file1)
	puzzle := sudokuPuzzle.SudokuPuzzle(file1String)

	then := time.Now()
	check := puzzle.Solve()
	duration := time.Since(then)

	if check == 0 {
		fmt.Printf("Puzzle solved in %v seconds!\nSolved puzzle stored in solved.txt\n", duration.Seconds())
	}

	puzzle.Output(false)
}
