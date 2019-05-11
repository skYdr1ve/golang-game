package game

type State byte

const (
	GOINGON State = iota
	DRAW
	PLAYER1WON
	PLAYER2WON
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
