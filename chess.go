package main

import (
	"fmt"
	"sort"
)

type (
	chess struct {
		board [][]int
	}

	aiitem struct {
		row, col    int // 行,列
		myvictory   int // 我方胜利度
		peervictory int // 对方胜利度
	}

	ailist []*aiitem
)

const (
	coin_cross = iota - 1
	coin_blank
	coin_zero
)

const (
	BoardSize   = 3	// 棋盘边长
	VictorySize = 3 // 几个连续的棋子指示胜利
	TotalSize   = BoardSize * BoardSize
)

func main() {

	var row, col, step int

	game := newChess()

	turn := coin_cross

	for step < TotalSize {

		game.output()
		fmt.Printf("\n%v turn: ", game.name(turn))

		if turn == coin_zero {
			fmt.Scanf("%d%d\n", &row, &col)
			game.input(row, col, turn)
		} else {
			game.ai(turn)
		}

		// 胜利或者换手
		if game.checkwin(game.board, turn) {
			game.output()
			fmt.Printf("\n%v win!\n", game.name(turn))
			return
		} else {
			turn = -turn
			step++
		}
	}

	game.output()
	fmt.Println("game end,no one win.")
}

func newChess() *chess {
	board := make([][]int, BoardSize)
	for idx := range board {
		board[idx] = make([]int, BoardSize)
	}
	return &chess{board: board}
}

func (this *chess) clone() *chess {
	game := newChess()
	for x := 0; x < BoardSize; x++ {
		for y := 0; y < BoardSize; y++ {
			game.board[x][y] = this.board[x][y]
		}
	}
	return game
}

func (this *chess) ai(turn int) {

	var list ailist

	// 计算胜利度
	for i := 0; i < TotalSize; i++ {

		row := i / BoardSize
		col := i % BoardSize

		if this.board[row][col] == coin_blank {

			board := this.clone().board
			board[row][col] = turn

			myvictory := this.victoryNess(board, turn)
			peervictory := this.victoryNess(board, -turn)

			item := &aiitem{
				row:         row,
				col:         col,
				myvictory:   myvictory,
				peervictory: peervictory,
			}

			list = append(list, item)
		}
	}

	// 排序
	sort.Sort(list)

	fmt.Printf("%v %v\n", list[0].row, list[0].col)
	this.input(list[0].row, list[0].col, turn)
}

func (this ailist) Len() int      { return len(this) }
func (this ailist) Swap(x, y int) { this[x], this[y] = this[y], this[x] }
func (this ailist) Less(x, y int) bool {
	// 对方胜利度不等时,将对方胜利度较小的条目排在前面,优先阻止对方胜利,然后考虑我方胜算
	if this[x].peervictory < this[y].peervictory {
		return true
	} else if this[x].peervictory == this[y].peervictory {
		return this[x].myvictory >= this[y].myvictory
	} else {
		return false
	}
}

// 胜算几何
func (this *chess) victoryNess(board [][]int, coin int) (victory int) {

	// 下标: 连续棋子数  值: 权重
	weight := make([]int, BoardSize+1)
	for i := 1; i <= BoardSize; i++ {
		if i == 1 {
			weight[i] = 1
		} else {
			weight[i] = weight[i-1] * TotalSize
		}
	}

	// 行/列/斜线
	row := func(x, y int) (a, b int) { a, b = x, y; return }
	col := func(x, y int) (a, b int) { a, b = y, x; return }
	cross1 := func(x, y int) (a, b int) { a, b = y, y; return }
	cross2 := func(x, y int) (a, b int) { a, b = y, BoardSize-1-y; return }

	/*
		固定一个坐标,检查所在行/列/斜线(由fn决定)上连续的我方棋子个数
	*/
	check := func(x int, fn func(int, int) (int, int)) (cnt int) {
		for y := 0; y < BoardSize; y++ {
			a, b := fn(x, y)	// 确定棋子坐标
			// 是我方棋子则连续棋子数增加
			// 不是空白位置,则是对方棋子挡路,我方不可能在这条线上取得胜利了
			if state := board[a][b]; state == coin {
				cnt++
			} else if state != coin_blank {
				cnt = 0
				break
			}
		}
		return
	}

	// 计算每行/每列的权重(胜利度)之和
	for i := 0; i < BoardSize; i++ {
		victory += weight[check(i, row)]
		victory += weight[check(i, col)]
	}
	
	// 计算两条斜线的权重(胜利度)之和
	victory += weight[check(0, cross1)]
	victory += weight[check(0, cross2)]

	return
}

func (this *chess) input(row, col, coin int) {
	this.board[row][col] = coin
}

func (this *chess) checkwin(board [][]int, coin int) bool {

	var win bool

	row := func(x, y int) (a, b int) { a, b = x, y; return }
	col := func(x, y int) (a, b int) { a, b = y, x; return }
	cross1 := func(x, y int) (a, b int) { a, b = y, y; return }
	cross2 := func(x, y int) (a, b int) { a, b = y, BoardSize-1-y; return }
	
	/*
		固定一个坐标,检查所在行/列/斜线(由fn决定)上连续的我方棋子个数
		如果全部为我方棋子,则胜利;如果有某个位置上不是我方棋子(对方或者空白)则没有胜利
	*/
	check := func(x int, fn func(int, int) (int, int)) (win bool) {
		win = true
		for y := 0; y < BoardSize; y++ {
			a, b := fn(x, y)
			if state := board[a][b]; state != coin {
				win = false
				break
			}
		}
		return
	}

	// 检查每行/每列上是否取得胜利
	for i := 0; i < BoardSize; i++ {
		win = win || check(i, row)
		win = win || check(i, col)
		if win {
			return win
		}
	}

	// 检查两条斜线上是否取得胜利
	win = win || check(0, cross1)
	win = win || check(0, cross2)

	return win
}

func (this *chess) output() {
	for _, row := range this.board {
		for _, coin := range row {
			fmt.Printf("%3v", this.name(coin))
		}
		fmt.Println("")
	}
}

func (this *chess) name(coin int) string {
	switch coin {
	case coin_blank:
		return "_"
	case coin_cross:
		return "X"
	case coin_zero:
		return "O"
	default:
		return "?"
	}
}
