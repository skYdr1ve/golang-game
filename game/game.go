package game

import "fmt"

type State byte

var FieldSizeInBytes = 9

const (
	GOINGON State = iota
	DRAW
	PLAYER1WON
	PLAYER2WON
	DISCONNECTED
)

type GameState struct {
	PlayingField []byte
	State        State
}

func (game *GameState) CheckState() {
	if game.PlayingField[0] == game.PlayingField[1] && game.PlayingField[1] == game.PlayingField[2] && game.PlayingField[0] != 0 {
		if game.PlayingField[0] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[0] == game.PlayingField[4] && game.PlayingField[4] == game.PlayingField[8] && game.PlayingField[0] != 0 {
		if game.PlayingField[0] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[0] == game.PlayingField[3] && game.PlayingField[3] == game.PlayingField[6] && game.PlayingField[0] != 0 {
		if game.PlayingField[0] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[3] == game.PlayingField[4] && game.PlayingField[4] == game.PlayingField[5] && game.PlayingField[3] != 0 {
		if game.PlayingField[3] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[6] == game.PlayingField[7] && game.PlayingField[7] == game.PlayingField[8] && game.PlayingField[6] != 0 {
		if game.PlayingField[6] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[1] == game.PlayingField[4] && game.PlayingField[4] == game.PlayingField[7] && game.PlayingField[1] != 0 {
		if game.PlayingField[1] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[2] == game.PlayingField[5] && game.PlayingField[5] == game.PlayingField[8] && game.PlayingField[2] != 0 {
		if game.PlayingField[2] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	if game.PlayingField[2] == game.PlayingField[4] && game.PlayingField[4] == game.PlayingField[6] && game.PlayingField[2] != 0 {
		if game.PlayingField[2] == 1 {
			game.State = PLAYER1WON
		} else {
			game.State = PLAYER2WON
		}
		return
	}
	isDraw := true
	for i := range game.PlayingField {
		if game.PlayingField[i] == 0 {
			isDraw = false
		}
	}
	if isDraw {
		game.State = DRAW
	} else {
		game.State = GOINGON
	}
}

func (game *GameState) ResetGame() {
	for i := range game.PlayingField {
		game.PlayingField[i] = 0
	}
	game.State = GOINGON
}

func New() GameState {
	newGame := GameState{
		PlayingField: make([]byte, 9),
		State:        GOINGON,
	}
	return newGame
}

//Function checks if the user made the move correctly
func Check(bytes []byte, x int, y int) bool {
	if x < 1 || x > 3 {
		return false
	}
	if y < 1 || y > 3 {
		return false
	}
	if bytes[(x-1)*3+y-1] != 0 {
		return false
	}
	return true
}

//Function draws a tic tac toe map
func DrawMap(bytes []byte) {
	fmt.Printf("\n-----------\n")
	for i := 0; i < 9; i++ {
		if i == 0 || i == 3 || i == 6 {
			fmt.Printf(" " + drawObj(int(bytes[i])))
		} else if i == 1 || i == 4 || i == 7 {
			fmt.Printf(" | " + drawObj(int(bytes[i])))
		} else {
			fmt.Printf(" | " + drawObj(int(bytes[i])) + " ")
			fmt.Printf("\n-----------\n")
		}
	}
	fmt.Println()
}

//Function draws elements on the map
func drawObj(obj int) string {
	if obj == 0 {
		return " "
	} else if obj == 1 {
		return "X"
	} else {
		return "O"
	}
}
