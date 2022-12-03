package rps_test

import (
	"strings"
	"testing"

	"github.com/aculclasure/aoc2022/rps"
)

func TestMatchOutcomeErrorCases(t *testing.T) {
	t.Parallel()
	const invalidPlay = "G"
	testCases := map[string]struct {
		opponentPlay string
		responsePlay string
	}{
		"Invalid opponent play returns error": {
			opponentPlay: invalidPlay,
			responsePlay: "X",
		},
		"Invalid response play returns error": {
			opponentPlay: "A",
			responsePlay: invalidPlay,
		},
	}
	game := rps.NewGame()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := game.MatchOutcome(tc.opponentPlay, tc.responsePlay)
			if err == nil {
				t.Error("expected an error but did not receive one")
			}
		})
	}
}

func TestMatchOutcomeSuccessCases(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		opponentPlay string
		responsePlay string
		want         int
	}{
		"Opponent scissors against paper response returns expected match points": {
			opponentPlay: "C",
			responsePlay: "Y",
			want:         2,
		},
		"Opponent scissors against rock response returns expected match points": {
			opponentPlay: "C",
			responsePlay: "X",
			want:         7,
		},
		"Opponent scissors against scissors response returns expected match points": {
			opponentPlay: "C",
			responsePlay: "Z",
			want:         6,
		},
		"Opponent paper against paper response returns expected match points": {
			opponentPlay: "B",
			responsePlay: "Y",
			want:         5,
		},
		"Opponent paper against rock response returns expected match points": {
			opponentPlay: "B",
			responsePlay: "X",
			want:         1,
		},
		"Opponent paper against scissors response returns expected match points": {
			opponentPlay: "B",
			responsePlay: "Z",
			want:         9,
		},
		"Opponent rock against paper response returns expected match points": {
			opponentPlay: "A",
			responsePlay: "Y",
			want:         8,
		},
		"Opponent rock against rock response returns expected match points": {
			opponentPlay: "A",
			responsePlay: "X",
			want:         4,
		},
		"Opponent rock against scissors response returns expected match points": {
			opponentPlay: "A",
			responsePlay: "Z",
			want:         3,
		},
	}
	game := rps.NewGame()

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := game.MatchOutcome(tc.opponentPlay, tc.responsePlay)
			if err != nil {
				t.Fatal("did not expect an error but got one: ", err)
			}
			if tc.want != got {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}

func TestComputeStrategyScore(t *testing.T) {
	t.Parallel()
	data := strings.NewReader(`A Y
B X
C Z
`)
	want := 15
	got, err := rps.ComputeStrategyScore(data)
	if err != nil {
		t.Fatal("got an unexpected error: ", err)
	}

	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
