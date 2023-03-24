package main

import (
	"Math/rand"
	"log"
	"fmt"
	"strconv"
	"strings"
)

type Cell struct {
	isMine        bool
	isFlagged     bool
	isOpen        bool
	adjacentMines int
	xpos int
	ypos int
}

type Board struct {
	rows  int
	cols  int
	cells [][]Cell
	mines int
	won   bool
	lost	bool
}

type Game struct {
	board Board
}

func (b *Board) adjacentCells(row int, col int) []Cell {
	var cells []Cell
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if r >= 0 && r < b.rows && c >= 0 && c < b.cols {
				b.cells[r][c].xpos = r
				b.cells[r][c].ypos = c
				cells = append(cells, b.cells[r][c])
			}
		}
	}
	return cells
}

func (b *Board) init() {
	mineCount := 0
	b.cells = make([][]Cell, b.rows, b.cols)
	for i := 0; i < b.rows; i++ {
		b.cells[i] = make([]Cell, b.rows)
	}
	for mineCount < b.mines {
		row := rand.Intn(b.rows)
		cols := rand.Intn(b.cols)
		if !b.cells[row][cols].isMine {
			b.cells[row][cols].isMine = true
			mineCount++
		}
	}
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			if !b.cells[row][col].isMine {
				for _, cell := range b.adjacentCells(row, col) {
					if cell.isMine {
						b.cells[row][col].adjacentMines++
					}
				}
			}
		}
	}
	b.display()
}

func (b *Board) symbol (c Cell) string {
	if c.isOpen {
		if c.isMine {
			b.lost = true
			return "\U0001F4A3"
		}
		return strconv.Itoa(c.adjacentMines)
	}
	if c.isFlagged {
		return "🚩"
	}
	return "◼"
}

func (b *Board) display() {
	for row := 0; row < b.rows; row++ {
		val := []string{}
		for col := 0; col < b.cols; col++ {
			val = append(val, b.symbol(b.cells[row][col]))
		}
		log.Printf("%s", val)
	}
}

func (b *Board) action(a string, row int, col int) {
  if a == "o" {
		b.cells[row][col].isOpen = true
		for _, cell := range b.adjacentCells(row, col) {
			if !cell.isMine {
				b.action("f", cell.xpos, cell.ypos)
				return
			}
		}
	} else {
		b.cells[row][col].isFlagged = true
	}
}

func (b *Board) move(c string) {
	str := strings.Split(c, ",")

	row, err := strconv.Atoi(str[0])
	if err != nil {
		log.Println("err", err)
	}

	col, err := strconv.Atoi(str[1])
	if err != nil {
		log.Println("err", err)
	}

	b.action(str[2], row, col)
	b.display()
}

func main() {
	board := new(Board)
	board.rows = 6
	board.cols = 6 
	board.mines = 5 
	board.won = false
	board.lost = false
	board.init()
	for {
		fmt.Print("Enter the row,col: ")
		var input string
		fmt.Scanln(&input)
		board.move(input)
		if board.won {
			log.Println("You won!")
			break
		}
		if board.lost {
			log.Println("oops!")
			break;
		}
	}
}
