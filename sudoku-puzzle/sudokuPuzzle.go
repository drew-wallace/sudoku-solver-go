package sudokuPuzzle

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// SudokuSolver ...
type SudokuSolver struct {
	grid [9][9][10]int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// ReadInts reads whitespace-separated ints from r. If there's an error, it returns the ints successfully read so far as well as the error value.
func ReadInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}

// SudokuPuzzle ...
func SudokuPuzzle(inputGridString string) SudokuSolver {
	ints, err := ReadInts(strings.NewReader(inputGridString))
	check(err)

	grid := [9][9][10]int{}

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			grid[r][c][0] = ints[r+c]
			for d := 1; d < 10; d++ {
				if grid[r][c][0] != 0 {
					grid[r][c][d] = 0 //sets all possible for that cell to 0 if a value was given
				} else {
					grid[r][c][d] = d //fills possible if no value was given
				}
			}
		}
	}

	return SudokuSolver{grid}
}

// Solve ...
func (puzzle SudokuSolver) Solve() int {
	fmt.Println("Solve", strconv.FormatBool(puzzle.isCorrect()), strconv.FormatBool(puzzle.isSolved()))

	if puzzle.isCorrect() && !puzzle.isSolved() {
		puzzle.setCellLoop()
	}
	if !puzzle.isCorrect() {
		return 1
	}
	if puzzle.isSolved() {
		return 0
	}
	// x and y are the coord. and n is the number possible
	x, y, n := puzzle.findLeastPoss()
	fmt.Println("Solve", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(n))
	for d := 1; d < 10; d++ {
		if puzzle.grid[x][y][d] != 0 {
			puzzle.grid[x][y][0] = puzzle.grid[x][y][d] //recursively goes through and checks each possibility. if the puzzle becomes solved it bails out. if not it check the next possibility
			fmt.Println("Solve", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(d))
			fmt.Println("Solve", strconv.Itoa(puzzle.grid[x][y][0]))
			rValue := puzzle.Solve()
			if rValue == 0 {
				return 0
			}
		}
	}

	return 1
}

// Output ...
func (puzzle SudokuSolver) Output(pretty bool) {
	f, err := os.Create("solved.txt")
	check(err)
	f.Sync()
	w := bufio.NewWriter(f)

	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			w.WriteString(strconv.Itoa(puzzle.grid[r][c][0]) + " ")
			if pretty && (c+1)%3 == 0 {
				w.WriteString("| ")
			}
		}
		w.WriteString("\n")
		if pretty && (r+1)%3 == 0 {
			w.WriteString("----------------------\n")
		}
	}
	w.Flush()
	defer f.Close()
}

func (puzzle SudokuSolver) zoneCheck(v int, cr int, cc int) bool { //value, current row, current column
	//zone row, zone column
	zr := 0
	zc := 0
	//determines what zone the value "v" is in. zr and zc are the coordinates of the starting element of the zone
	if cr >= 0 && cr <= 2 {
		zr = 0 //top row
		if cc >= 0 && cc <= 2 {
			zc = 0
		} else if cc >= 3 && cc <= 5 {
			zc = 3
		} else if cc >= 6 && cc <= 8 {
			zc = 6
		}
	} else if cr >= 3 && cr <= 5 {
		zr = 3 //middle row
		if cc >= 0 && cc <= 2 {
			zc = 0
		} else if cc >= 3 && cc <= 5 {
			zc = 3
		} else if cc >= 6 && cc <= 8 {
			zc = 6
		}
	} else if cr >= 6 && cr <= 8 {
		zr = 6 //bottom row
		if cc >= 0 && cc <= 2 {
			zc = 0
		} else if cc >= 3 && cc <= 5 {
			zc = 3
		} else if cc >= 6 && cc <= 8 {
			zc = 6
		}
	}
	//only searches the 3x3 zone. returns false if the zone already contains the value "v"
	for r := zr; r < (zr + 3); r++ {
		for c := zc; c < (zc + 3); c++ {
			if puzzle.grid[r][c][0] == v {
				return false
			}
		}
	}
	return true
}

func (puzzle SudokuSolver) zoneSet(v int, cr int, cc int) {
	zr := 0
	zc := 0
	//determines what zone the value "v" is in. zr and zc are the coordinates of the starting element of the zone
	if cr >= 0 && cr <= 2 {
		zr = 0
		if cc >= 0 && cc <= 2 {
			zc = 0
		} else if cc >= 3 && cc <= 5 {
			zc = 3
		} else if cc >= 6 && cc <= 8 {
			zc = 6
		}
	} else if cr >= 3 && cr <= 5 {
		zr = 3
		if cc >= 0 && cc <= 2 {
			zc = 0
		} else if cc >= 3 && cc <= 5 {
			zc = 3
		} else if cc >= 6 && cc <= 8 {
			zc = 6
		}
	} else if cr >= 6 && cr <= 8 {
		zr = 6
		if cc >= 0 && cc <= 2 {
			zc = 0
		} else if cc >= 3 && cc <= 5 {
			zc = 3
		} else if cc >= 6 && cc <= 8 {
			zc = 6
		}
	}
	for r := zr; r < (zr + 3); r++ {
		for c := zc; c < (zc + 3); c++ {
			//if the possible value is available then set it to 0
			if puzzle.grid[r][c][v] == v {
				puzzle.grid[r][c][v] = 0
			}
		}
	}
}

func (puzzle SudokuSolver) rowCheck(v int, cr int) bool {
	//searches the given row. returns false if the row already contains the value "v"
	for c := 0; c < 9; c++ {
		if puzzle.grid[cr][c][0] == v {
			return false
		}
	}
	return true
}

func (puzzle SudokuSolver) rowSet(v int, cr int) {
	for c := 0; c < 9; c++ {
		//if the possible value is available then set it to 0
		if puzzle.grid[cr][c][v] == v {
			puzzle.grid[cr][c][v] = 0
		}
	}
}

func (puzzle SudokuSolver) colCheck(v int, cc int) bool {
	//searches the given column. returns false if the column already contains the value "v"
	for r := 0; r < 9; r++ {
		if puzzle.grid[r][cc][0] == v {
			return false
		}
	}
	return true
}

func (puzzle SudokuSolver) colSet(v int, cc int) {
	for r := 0; r < 9; r++ {
		//if the possible value is available then set it to 0
		if puzzle.grid[r][cc][v] == v {
			puzzle.grid[r][cc][v] = 0
		}
	}
}

//checks a single cell for the 3 cases: row, column, and zone
func (puzzle SudokuSolver) cellCheck(v int, r int, c int) bool {
	return puzzle.zoneCheck(v, r, c) && puzzle.rowCheck(v, r) && puzzle.colCheck(v, c)
}

func (puzzle SudokuSolver) setPossible() {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			//runs through the whole grid and sets the possible values for each cell
			if puzzle.grid[r][c][0] != 0 {
				v := puzzle.grid[r][c][0]
				puzzle.colSet(v, c)
				puzzle.rowSet(v, r)
				puzzle.zoneSet(v, r, c)
			}
		}
	}
}

func (puzzle SudokuSolver) setCellLoop() {
	count := 0 //counts possible
	sd := 0    //singles possibility var.
	changed := true
	tmp := [9]int{} //array of possible values for the cells in a zone/row/column
	zr := 0
	zc := 0 //zone row, zone column
	for {
		changed = false
		for r := 0; r < 9; r++ {
			for c := 0; c < 9; c++ {
				puzzle.setPossible()
				if puzzle.grid[r][c][0] == 0 {
					//set singles
					for d := 1; d < 10; d++ {
						//if a number is found in the empty cell (cell of value 0), then search through the possibilities and count for every possiblility. also store the coord. of the current possibility found
						if puzzle.grid[r][c][d] != 0 {
							count++
							sd = d
						}
					}
					//if only one possibility found, use the stored coord. to change the cell value to the possibility
					if count == 1 {
						fmt.Println(strconv.Itoa(r), strconv.Itoa(c), strconv.Itoa(sd))
						puzzle.grid[r][c][0] = sd
						changed = true
						goto eocl
					} else if count == 0 {
						goto notCorrect
					}

					count = 0

					//for unique
					for i := 0; i < 9; i++ {
						tmp[i] = 0
					}

					//for zone
					//determines what zone the cell is in. zr and zc are the coordinates of the starting element of the zone
					if r >= 0 && r <= 2 {
						zr = 0 //top row
						if c >= 0 && c <= 2 {
							zc = 0
						} else if c >= 3 && c <= 5 {
							zc = 3
						} else if c >= 6 && c <= 8 {
							zc = 6
						}
					} else if r >= 3 && r <= 5 {
						zr = 3 //middle row
						if c >= 0 && c <= 2 {
							zc = 0
						} else if c >= 3 && c <= 5 {
							zc = 3
						} else if c >= 6 && c <= 8 {
							zc = 6
						}
					} else if r >= 6 && r <= 8 {
						zr = 6 //bottom row
						if c >= 0 && c <= 2 {
							zc = 0
						} else if c >= 3 && c <= 5 {
							zc = 3
						} else if c >= 6 && c <= 8 {
							zc = 6
						}
					}
					for cr := zr; cr < (zr + 3); cr++ {
						for cc := zc; cc < (zc + 3); cc++ {
							if cr == r && cc != c && puzzle.grid[cr][cc][0] == 0 {
								for cd := 1; cd < 10; cd++ {
									if puzzle.grid[cr][cc][cd] != 0 {
										//stores value of possibility into index of tmp. this prevents multiple of same value
										tmp[puzzle.grid[cr][cc][cd]-1] = puzzle.grid[cr][cc][cd]
									}
								}
							} else if cr != r && puzzle.grid[cr][cc][0] == 0 {
								for cd := 1; cd < 10; cd++ {
									if puzzle.grid[cr][cc][cd] != 0 {
										//stores value of possibility into index of tmp. this prevents multiple of same value
										tmp[puzzle.grid[cr][cc][cd]-1] = puzzle.grid[cr][cc][cd]
									}
								}
							}
						}
					}
					for d := 1; d < 10; d++ {
						//if a possible value of the cell is unique to the zone, then set the cell to that value
						if puzzle.grid[r][c][d] != tmp[puzzle.grid[r][c][d]-1] && puzzle.grid[r][c][d] != 0 {
							fmt.Println(strconv.Itoa(r), strconv.Itoa(c), strconv.Itoa(sd))
							puzzle.grid[r][c][0] = puzzle.grid[r][c][d]
							changed = true
							goto eocl
						}
					}

					for i := 0; i < 9; i++ {
						tmp[i] = 0
					}

					//for row
					for cc := 0; cc < 9; cc++ {
						if cc != c && puzzle.grid[r][cc][0] == 0 {
							for cd := 1; cd < 10; cd++ {
								if puzzle.grid[r][cc][cd] != 0 {
									tmp[puzzle.grid[r][cc][cd]-1] = puzzle.grid[r][cc][cd]
								}
							}
						}
					}
					for d := 1; d < 10; d++ {
						//if a possible value of the cell is unique to the row, then set the cell to that value
						if puzzle.grid[r][c][d] != tmp[puzzle.grid[r][c][d]-1] && puzzle.grid[r][c][d] != 0 {
							fmt.Println(strconv.Itoa(r), strconv.Itoa(c), strconv.Itoa(sd))
							puzzle.grid[r][c][0] = puzzle.grid[r][c][d]
							changed = true
							goto eocl
						}
					}

					for i := 0; i < 9; i++ {
						tmp[i] = 0
					}

					//for col
					for cr := 0; cr < 9; cr++ {
						if cr != r && puzzle.grid[cr][c][0] == 0 {
							for cd := 1; cd < 10; cd++ {
								if puzzle.grid[cr][c][cd] != 0 {
									tmp[puzzle.grid[cr][c][cd]-1] = puzzle.grid[cr][c][cd]
								}
							}
						}
					}

					for d := 1; d < 10; d++ {
						//if a possible value of the cell is unique to the column, then set the cell to that value
						if puzzle.grid[r][c][d] != tmp[puzzle.grid[r][c][d]-1] && puzzle.grid[r][c][d] != 0 {
							fmt.Println(strconv.Itoa(r), strconv.Itoa(c), strconv.Itoa(sd))
							puzzle.grid[r][c][0] = puzzle.grid[r][c][d]
							changed = true
							goto eocl
						}
					}
				eocl:
				}
			}
		}
		if !changed {
			break
		}
	}
notCorrect:
}

func (puzzle SudokuSolver) isSolved() bool {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle.grid[r][c][0] == 0 {
				return false
			}
		}
	}
	return true
}

func (puzzle SudokuSolver) isCorrect() bool {
	count := 0
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle.grid[r][c][0] == 0 {
				for d := 1; d < 10; d++ {
					if puzzle.grid[r][c][d] == 0 {
						count++
					}
				}
				if count == 9 {
					return false
				}
				count = 0
			}
		}
	}
	return true
}

func (puzzle SudokuSolver) findLeastPoss() (int, int, int) {
	count := 0
	x := 0
	y := 0
	n := 9
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if puzzle.grid[r][c][0] == 0 {
				for d := 1; d < 10; d++ {
					if puzzle.grid[r][c][d] != 0 {
						count++ //counts the number possible
					}
				}
				if count < n {
					n = count //if the number possible is less than the global possible save it and the coord.
					x = r
					y = c
				} else {
					count = 0
				}
			}
		}
	}

	fmt.Println("FLP", strconv.Itoa(x), strconv.Itoa(y), strconv.Itoa(n))
	return x, y, n
}
