package xogame

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Game struct {
	gameboard   board
	playerXName string
	playerOName string
	turn        bool // true = x & false = o
}
type state int

const (
	XWIN       state = iota
	OWIN       state = iota
	DRAW       state = iota
	INPROGRESS state = iota
)

type board struct {
	boardSlice []string
	gamestate  state
	xchar      string
	ochar      string
}

type move struct {
	player bool
	square string
}

// exported functions
// NewGame initializes the Game struct and returns a pointer to that new game struct.
// It also sets up the initial text and welcomes the player
func NewGame() *Game {
	game := Game{
		gameboard: board{
			boardSlice: []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
			gamestate:  INPROGRESS,
			xchar:      "X",
			ochar:      "O",
		},
		playerXName: "playerx",
		playerOName: "playero",
		turn:        true,
	}
	game.setup()
	return &game
}

func (game *Game) Play() (state, error) {
	for {
		if game.checkStatus() != INPROGRESS {
			return game.checkStatus(), nil
		}
		move, err := game.getMove()
		if err != nil {
			return INPROGRESS, err
		}
		game.MakeMove(move)
	}
}

func Congrat(game *Game, status state) {
	switch status {
	case XWIN:
		fmt.Printf("%s has won!\n", game.playerXName)
	case OWIN:
		fmt.Printf("%s has won!\n", game.playerOName)
	case DRAW:
		fmt.Print("It's a draw :(\n")
	}
}

// setup sets up a new game by initializing player names and priting initial help info
func (game *Game) setup() {
	fmt.Println(`Welcome to TicTacGo!
Our rules are simple:
	- Each square on the board is numbered from 1 to 9.
	- When it is your turn, you enter the number of the square you want to play at
	- That's it! You expected more! Silly! Haha.
	`)
	var playerx, playero string
	fmt.Print("Enter playerX name: ")
	fmt.Scanln(&playerx)
	game.playerXName = playerx
	fmt.Print("Enter PlayerO name: ")
	fmt.Scanln(&playero)
	game.playerOName = playero
	drawBoard(&game.gameboard)
}

// MakeMove validates player moves and makes them
func (game *Game) MakeMove(move move) error {
	// check if it is the correct turn for the player making the move
	if err := move.correctTurnForPlayer(game); err != nil {
		return err
	}
	// check if the move if valid
	if valid := move.isValidMove(); !valid {
		msg := fmt.Sprintf("%v is not a valid move. Valid moves are from 0 to 9", move.square)
		return errors.New(msg)
	}
	// parse move
	actualMove, err := move.getActualMove()
	if err != nil {
		msg := fmt.Sprintf("Your move: %v is not convertable to a board move. Double check your move.", move.square)
		return errors.New(msg)
	}
	if game.turn { // make the move
		game.gameboard.boardSlice[actualMove] = game.gameboard.xchar
	} else {
		game.gameboard.boardSlice[actualMove] = game.gameboard.ochar
	}
	game.turn = !game.turn // switch players

	game.adjustState()
	drawBoard(&game.gameboard) // draw the board
	return nil
}

func (game *Game) checkStatus() state {
	return game.gameboard.gamestate
}

func (game *Game) getMove() (move, error) {
	prompt := fmt.Sprintf("%s>> ", game.getPlayer())
	fmt.Print(prompt)
	var sqr string
	_, err := fmt.Scanln(&sqr)
	if err != nil {
		return move{}, errors.New("Didn't get move")
	}
	return move{player: game.turn, square: sqr}, nil
}

func (game *Game) getPlayer() string {
	if game.turn {
		return game.playerXName
	} else {
		return game.playerOName
	}
}

func (game *Game) adjustState() {
	switch {
	case checkWin(&game.gameboard, game.gameboard.xchar) == true:
		game.gameboard.gamestate = XWIN
	case checkWin(&game.gameboard, game.gameboard.ochar) == true:
		game.gameboard.gamestate = OWIN
	case checkDraw(&game.gameboard) == true:
		game.gameboard.gamestate = DRAW
	default:
		game.gameboard.gamestate = INPROGRESS
	}
}

func drawBoard(board *board) {
	gboard := board.boardSlice
	prettyBoard := fmt.Sprintf("%s | %s | %s\n%s | %s | %s\n%s | %s | %s", gboard[0], gboard[1], gboard[2], gboard[3], gboard[4], gboard[5], gboard[6], gboard[7], gboard[8])
	fmt.Println(prettyBoard)
}

func checkWin(board *board, char string) bool {
	// hori wins
	gboard := board.boardSlice
	if (gboard[0] == char && gboard[1] == char && gboard[2] == char) ||
		(gboard[3] == char && gboard[4] == char && gboard[5] == char) ||
		(gboard[6] == char && gboard[7] == char && gboard[8] == char) ||
		// vert wins
		(gboard[0] == char && gboard[3] == char && gboard[6] == char) ||
		(gboard[1] == char && gboard[4] == char && gboard[7] == char) ||
		(gboard[2] == char && gboard[5] == char && gboard[8] == char) ||
		// diag wins
		(gboard[0] == char && gboard[4] == char && gboard[8] == char) ||
		(gboard[2] == char && gboard[4] == char && gboard[6] == char) {
		return true
	}
	return false
}

func checkDraw(gboard *board) bool {
	return numberOfLegalMoves(gboard) == 0
}

func numberOfLegalMoves(board *board) int {
	num := 0
	re := regexp.MustCompile(`[0-9]`)
	for _, val := range board.boardSlice {
		if re.MatchString(val) {
			num++
		}
	}
	return num
}

func (move move) getActualMove() (int, error) {
	intmove, err := strconv.Atoi(move.square)
	intmove--
	return intmove, err
}

func (move move) isValidMove() bool {
	re := regexp.MustCompile(`[1-9]`)
	if len(move.square) == 1 && re.MatchString(move.square) {
		return true
	} else {
		return false
	}
}

func (move move) correctTurnForPlayer(game *Game) error {
	if game.turn { // player x turn
		if !move.player { // player o trying to make a move
			msg := fmt.Sprintf("Player %s's turn, not yours", game.playerXName)
			return errors.New(msg)
		}
	}
	if !game.turn { // player o turn
		if move.player { // player x trying to make a move
			msg := fmt.Sprintf("Player %s's turn, not yours", game.playerOName)
			return errors.New(msg)
		}
	}
	return nil
}
