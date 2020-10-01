package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type face string

const (
	faceBrain   face = "Brain"
	faceShotgun face = "Shotgun"
	faceRunner  face = "Runner"
)

type color string

const (
	colorGreen  color = "Green"
	colorYellow color = "Yellow"
	colorRed    color = "Red"
)

type die struct {
	color color
	faces [6]face
}

var dieGreen = die{color: colorGreen, faces: [6]face{faceBrain, faceBrain, faceBrain, faceShotgun, faceRunner, faceRunner}}
var dieYellow = die{color: colorYellow, faces: [6]face{faceBrain, faceBrain, faceShotgun, faceShotgun, faceRunner, faceRunner}}
var dieRed = die{color: colorRed, faces: [6]face{faceBrain, faceShotgun, faceShotgun, faceShotgun, faceRunner, faceRunner}}

// 6 green, 4 yellow, 3 red
var initDice = []die{dieGreen, dieGreen, dieGreen, dieGreen, dieGreen, dieGreen,
	dieYellow, dieYellow, dieYellow, dieYellow, dieRed, dieRed, dieRed}

/* Gameplay
The player has to shake a cup containing 13 dice and randomly select 3 of them without
looking into the cup and then roll them. The faces of each die represent brains, shotgun
blasts or "runners" with different colours containing a different distribution of faces
(the 6 green dice have 3 brains, 1 shotgun and 2 runners, the 4 yellow dice have 2 of each
and the 3 red dice have 1 brain, 3 shotguns and 2 runners). The object of the game is
to roll 13 brains. If a player rolls 3 shotgun blasts their turn ends and they lose
the brains they have accumulated so far that turn. It is possible for a player to roll
3 blasts in a single roll, but if only one or two blasts have been rolled the player will
have to decide whether it is worth it to risk rolling again or "bank" the brains acquired
so far and pass play to the next player. A "runner" is represented by feet and rolling a
runner means that the player can roll that same dice if they choose to press their luck.
A winner is determined if a player rolls 13 brains and all other players have taken at least
one more turn without reaching 13 brains.
*/

var playerScores = []int{0, 0}

const (
	player1 = 1
	player2 = 2
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	selectionBox := initDice

	var handBrainDice []die
	var handShotgunDice []die
	var handRunnerDice []die
	var diceAtHand []die
	var player int

	player = nextPlayer(player)

	turn := 1
	for {
		// Reached 13 or more brains means the game has finished
		if player == player1 && (playerScores[0] >= 13 || playerScores[1] >= 13) {
			winner := player1
			if playerScores[1] > playerScores[0] {
				winner = player2
			}
			fmt.Printf("Player %d is winner after %d turns\n", winner, turn)
			fmt.Printf("Final score: Player 1: %d - Player 2: %d\n", playerScores[0], playerScores[1])
			return
		}

		fmt.Printf("===> Turn: %d Player %d's turn\n", turn, player)
		fmt.Printf("Score: Player 1: %d - Player 2: %d\n", playerScores[0], playerScores[1])
		fmt.Printf("Available number of dice: %d\n", len(selectionBox))

		for {
			if len(diceAtHand) == 3 || len(selectionBox) == 0 {
				break
			}
			n := rand.Intn(len(selectionBox))
			diceAtHand = append(diceAtHand, selectionBox[n])
			selectionBox = append(selectionBox[:n], selectionBox[n+1:]...)
		}

		for i := 0; i < len(diceAtHand); i++ {
			fmt.Printf("Selected die %d: %s\n", i+1, diceAtHand[i].color)
		}
		faces := rollDice(diceAtHand)

		for i := 0; i < len(diceAtHand); i++ {
			fmt.Printf("Die %d: %s %s\n", i+1, diceAtHand[i].color, faces[i])
			if faces[i] == faceBrain {
				handBrainDice = append(handBrainDice, diceAtHand[i])
			}
			if faces[i] == faceShotgun {
				handShotgunDice = append(handShotgunDice, diceAtHand[i])
			}
			if faces[i] == faceRunner {
				handRunnerDice = append(handRunnerDice, diceAtHand[i])
			}
		}

		fmt.Printf("At hand: %d Brain - %d Runner - %d Shotgun\n", len(handBrainDice), len(handRunnerDice), len(handShotgunDice))

		// 3 shotgun is you are dead
		if len(handShotgunDice) >= 3 {
			fmt.Printf("3 or more shotgun. Busted for this turn\n")
			selectionBox = append(selectionBox, handBrainDice...)
			handBrainDice = []die{}
			selectionBox = append(selectionBox, handShotgunDice...)
			handShotgunDice = []die{}
			selectionBox = append(selectionBox, handRunnerDice...)
			handRunnerDice = []die{}
			diceAtHand = []die{}
			player = nextPlayer(player)
			if player == player2 {
				turn++
			}

			fmt.Print("Press enter to continue...")
			fmt.Scanln()
			continue
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Continue? ... y/n ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("y", text) == 0 {
			diceAtHand = handRunnerDice
			handRunnerDice = []die{}
		}
		if strings.Compare("n", text) == 0 {
			playerScores[player-1] += len(handBrainDice)
			selectionBox = append(selectionBox, handBrainDice...)
			handBrainDice = []die{}
			selectionBox = append(selectionBox, handShotgunDice...)
			handShotgunDice = []die{}
			selectionBox = append(selectionBox, handRunnerDice...)
			handRunnerDice = []die{}
			diceAtHand = []die{}
			player = nextPlayer(player)
			if player == player2 {
				turn++
			}
		}
	}
}

func rollDice(dice []die) []face {
	fmt.Printf("Rolling %d dice\n", len(dice))

	var faces []face
	for i := 0; i < len(dice); i++ {
		n := rand.Intn(6)
		faces = append(faces, dice[i].faces[n])
	}
	return faces
}

func nextPlayer(player int) int {
	if player == player1 {
		return player2
	}

	return player1
}
