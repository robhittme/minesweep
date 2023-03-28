package minesweeper

import (
  "strconv"
  "log"
  "math/rand"
  "fmt"
  "strings"
  "time"
) 
/** 
--init game
--create board
--place mines randomly
--add cell values
add action
--display board
adjacent opening logic
evaluate won/lost
for loop for game.
**/

type GameStatus string
type CellState string

const (
  Won GameStatus = "Won" 
  Lost GameStatus = "Lost" 
  Active GameStatus = "Active" 
)


type Cell struct {
  Opened bool
  Flagged bool
  Questioned bool
  IsMine bool
  AdjacentMines int
  Row int
  Col int
}

type Game struct {
  Rows int
  Cols int
  Count int
  Mines int
  Cells [][]Cell
  Status GameStatus
}

func (c *Cell) Display() string {
  if c.Opened {
    if c.IsMine {
      return "[*]"
    }
    return "[" + strconv.Itoa(c.AdjacentMines) + "]"
  }
  if c.Flagged {
    return "[>]"
  }
  if c.Questioned {
    return "[?]" 
  }
  return "[ ]"
}

func (c *Cell) AdjacentCells(g Game) []Cell {
  cells := []Cell{}
  for row := c.Row - 1; row <= c.Row+1; row++ {
    for col := c.Col - 1; col <= c.Col+1; col++ {
      if row < g.Rows && row >= 0 && col < g.Cols && col >= 0 {
        cells = append(cells, g.Cells[row][col])
      }
    }
  }
  return cells 
}

func (c *Cell) OpenAdjacentCells(g Game) {
  for _, cell := range c.AdjacentCells(g) {
    if !cell.IsMine && !cell.Opened && c.AdjacentMines == 0 {
      g.Cells[cell.Row][cell.Col].Opened = true
      g.Cells[cell.Row][cell.Col].OpenAdjacentCells(g)
    }
  }
}

func (g *Game) PopulateCells() {
  for row := 0; row < g.Rows; row++ {
    for col := 0; col < g.Cols; col++ {
      for _, cell := range g.Cells[row][col].AdjacentCells(*g) {
        if cell.IsMine {
          g.Cells[row][col].AdjacentMines++
        }
      }
    }
  }
}

func (g *Game) DisperseMines() {
  mineCount := g.Mines
  rand.Seed(time.Now().UnixNano())
  for mineCount > 0 {
    row := rand.Intn(g.Rows)
    col := rand.Intn(g.Cols)
    if !g.Cells[row][col].IsMine {
      g.Cells[row][col].IsMine = true
      mineCount--
    }
  }
}

func (g *Game) Display() {
  guide := []string{" X "}
  for a := 0; a < g.Rows; a++ {
    guide = append(guide, " "+strconv.Itoa(a)+ " ")
  }
  log.Println(guide)
  for row := 0; row < g.Rows; row++ {
    val := []string{" " + strconv.Itoa(row) + " "}
    for col := 0; col < g.Cols; col++ {
      val = append(val, g.Cells[row][col].Display())  
    }
    log.Printf("%s", val)
  }
}
func (g *Game) Init() {
  g.Cells = make([][]Cell, g.Rows)  
  for i := 0; i < g.Rows; i++ {
    g.Cells[i] = make([]Cell, g.Cols)
    for j := 0; j < g.Cols; j++ {
      g.Cells[i][j] = Cell{
        Row: i,
        Col: j,
        Opened: true,
      }
    }
  }
}

func (g *Game) Evaluate() {
  count := 0
  for row := 0; row < g.Rows; row++ {
    for col := 0; col < g.Cols; col++ {
      if g.Cells[row][col].IsMine && g.Cells[row][col].Opened {
        g.Status = Lost
      }
      count++
    }
  }
  if count == g.Count {
    g.Status = Won
  }
}

func (g *Game) Action(input string) {
  splitStr := strings.Split(input, ",")
  row, err := strconv.Atoi(splitStr[0])
  if err != nil {
    log.Fatal("wrong input")
  }
  col, err := strconv.Atoi(splitStr[1])
  if err != nil {
    log.Fatal("wrong input")
  }
  command := splitStr[2]

  switch command {
  case "o":
    g.Cells[row][col].Opened = true
    if g.Cells[row][col].AdjacentMines == 0 {
      g.Cells[row][col].OpenAdjacentCells(*g)
    }

  case "f":
    g.Cells[row][col].Flagged = true 
  case "q":
    g.Cells[row][col].Questioned = true
}

}
func sweep() {
  g := &Game{
    Rows: 10,
    Cols: 10,
    Mines: 10, 
    Status: Active,
    Count: 10 * 10 - 10,
  }
  g.Init()
  g.DisperseMines()
  g.PopulateCells()
  g.Display()
  for {
    fmt.Println("row,col,command(o,f,q): ") 
    var input string
    fmt.Scanln(&input)
    g.Action(input)
    g.Display()
    g.Evaluate()
    switch g.Status {
    case Lost:
      fmt.Println("oops! you lost")
    case Won:
      fmt.Println("Whoo!! you Won")
  }
  }
}
