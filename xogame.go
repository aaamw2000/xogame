package xogame

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"golang.org/x/term"
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
	square int
}

// exported functions
// NewGame initializes the Game struct and returns a pointer to that new game struct.
// It also sets up the initial text and welcomes the player
func NewGame() (*Game, error) {
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
	err := game.setup()
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (game *Game) Play() (state, error) {
	for {
		if game.checkStatus() != INPROGRESS {
			return game.checkStatus(), nil
		}
		game.MakeMove()
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
func (game *Game) setup() error {
	welcome := "Welcome to TicTacGo!\nOur rules are simple:"
	rule0 := "- Each square on the board is numbered from 1 to 9."
	rule1 := "- When it is your turn, you enter the number of the square you want to play at."
	rule2 := "- That's it! Those are all the rules. You expected more? Silly! Haha. Have fun playing!"

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return errors.New("Couldn't get terminal size.")
	}
	welcomeStyle := lipgloss.NewStyle().
		Bold(true).
		Width(width - 2).
		Height(height - 2000).
		Align(lipgloss.Center).
		Background(lipgloss.Color("#7D56F4")).
		Padding(2).
		BorderStyle(lipgloss.Border{
			Top:         "._.:*:",
			Bottom:      "._.:*:",
			Left:        "|*",
			Right:       "|*",
			TopLeft:     "*",
			TopRight:    "*",
			BottomLeft:  "*",
			BottomRight: "*",
		})
	welcomeMsg := lipgloss.JoinVertical(lipgloss.Left, welcome, rule0, rule1, rule2)

	fmt.Println(welcomeStyle.Render(welcomeMsg))
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4"))

	var playerx, playero string
	fmt.Print(style.Render("Enter PlayerX name:") + " ")
	fmt.Scanln(&playerx)
	game.playerXName = playerx
	fmt.Print(style.Render("Enter PlayerO name:") + " ")
	fmt.Scanln(&playero)
	game.playerOName = playero
	drawBoard(&game.gameboard)
	return nil
}

// MakeMove validates player moves and makes them
func (game *Game) MakeMove() error {
	actualMove, err := game.betterGetMove()
	if err != nil {
		return err
	}
	moveInt := actualMove.square - 1

	if game.turn { // make the move
		game.gameboard.boardSlice[moveInt] = game.gameboard.xchar
	} else {
		game.gameboard.boardSlice[moveInt] = game.gameboard.ochar
	}
	game.turn = !game.turn // switch players

	game.adjustState()
	drawBoard(&game.gameboard) // draw the board
	return nil
}

func (game *Game) checkStatus() state {
	return game.gameboard.gamestate
}

func (gameboard *board) isEmptySquare(sqr int) bool {
	return gameboard.boardSlice[sqr-1] != "X" && gameboard.boardSlice[sqr-1] != "O"
}

func (game *Game) setPrompt() string {
	if game.turn {
		return fmt.Sprintf("(%s) Move [1-9]", game.playerXName)
	} else {
		return fmt.Sprintf("(%s) Move [1-9]", game.playerOName)
	}
}

func (game *Game) betterGetMove() (move, error) {
	validate := func(input string) error {
		parsedInt, err := strconv.ParseInt(input, 10, 8)
		if err != nil {
			return errors.New("Invalid move character!")
		}
		if parsedInt < 1 || parsedInt > 9 {
			return errors.New("Invalid move number!")
		}
		if !game.gameboard.isEmptySquare(int(parsedInt)) {
			return errors.New("Invalid move; square occupied.")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    game.setPrompt(),
		Validate: validate,
	}

	moveStr, err := prompt.Run()
	if err != nil {
		err := errors.New("Failed to get input")
		return move{}, err
	}

	moveInt, _ := strconv.Atoi(moveStr)

	return move{player: game.turn, square: moveInt}, nil
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
