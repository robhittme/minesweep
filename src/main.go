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
	lost  bool
	count int
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

	log.Println("[x 0 1 2]")
	for row := 0; row < b.rows; row++ {
		val := []string{strconv.Itoa(row)}
		for col := 0; col < b.cols; col++ {
			val = append(val, b.symbol(b.cells[row][col]))
		}
		log.Printf("%s", val)
	}
}

func (b *Board) openAdjacent(row int, col int) {
	for _, cell := range b.adjacentCells(row, col) {
		if !cell.isMine && !cell.isOpen {
			b.cells[cell.xpos][cell.ypos].isOpen = true
			if cell.adjacentMines == 0 {
				b.openAdjacent(cell.xpos, cell.ypos)
			}
		}
	}
}

func (b *Board) action(a string, row int, col int) {
	if a == "o" {
		b.cells[row][col].isOpen = true
		if b.cells[row][col].adjacentMines == 0 {
			for _, cell := range b.adjacentCells(row, col) {
				if !cell.isMine {
					b.cells[cell.xpos][cell.ypos].isOpen = true
					b.openAdjacent(cell.xpos, cell.ypos)
				}
			}
		}
	} else {
		b.cells[row][col].isFlagged = true
	}
}

func (b *Board) evaluate() {
	opened := 0
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			if b.cells[row][col].isOpen {
				opened++
			}
		}
	}
	if opened == b.count {
		b.won = true
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
	b.evaluate()
	b.display()
}

func main() {
	b := new(Board)
	b = &Board{
		rows: 3,
		cols: 3,
		mines: 3,
		won: false,
		lost: false,
	}
	b.count = b.rows * b.cols - b.mines
	b.init()
	for {
		fmt.Print("Enter the row,col,command(o,f): ")
		var input string
		fmt.Scanln(&input)
		b.move(input)
		if b.won {
			log.Println("You won!😎")
			break
		}
		if b.lost {
			log.Println("boom! 😵")
			break;
		}
	}
}
