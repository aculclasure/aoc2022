package rps

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Game struct {
	responses         map[string]int
	opponentWinRules  map[string]string
	opponentDrawRules map[string]string
	opponentLossRules map[string]string
	matchLossPoints   int
	matchDrawPoints   int
	matchWinPoints    int
}

func NewGame() Game {
	return Game{
		responses: map[string]int{
			"X": 1, // rock
			"Y": 2, // paper
			"Z": 3, // scissors
		},
		opponentWinRules: map[string]string{
			"A": "Z", // rock beats scissors
			"B": "X", // paper beats rock
			"C": "Y", // scissors beats paper
		},
		opponentDrawRules: map[string]string{
			"A": "X", // rock draws with rock
			"B": "Y", // paper draws with paper
			"C": "Z", // scissors draws with scissors
		},
		opponentLossRules: map[string]string{
			"A": "Y", // rock loses to paper
			"B": "Z", // paper loses to scissors
			"C": "X", // scissors loses to rock
		},
		matchLossPoints: 0,
		matchDrawPoints: 3,
		matchWinPoints:  6,
	}
}

func (g Game) MatchOutcome(opponentPlay, responsePlay string) (int, error) {
	_, ok := g.opponentWinRules[opponentPlay]
	if !ok {
		return 0, fmt.Errorf("opponent play must be one of A, B, C (got %s)", opponentPlay)
	}
	_, ok = g.responses[responsePlay]
	if !ok {
		return 0, fmt.Errorf("response play must be one of X, Y, Z (got %s)", responsePlay)
	}

	switch {
	case responsePlay == g.opponentWinRules[opponentPlay]:
		return g.matchLossPoints + g.responses[responsePlay], nil
	case responsePlay == g.opponentDrawRules[opponentPlay]:
		return g.matchDrawPoints + g.responses[responsePlay], nil
	default:
		return g.matchWinPoints + g.responses[responsePlay], nil
	}
}

func (g Game) CheatMatchOutcome(opponentPlay, cheatMatchResult string) (int, error) {
	_, ok := g.opponentWinRules[opponentPlay]
	if !ok {
		return 0, fmt.Errorf("opponent play must be one of A, B, C (got %s)", opponentPlay)
	}
	_, ok = g.responses[cheatMatchResult]
	if !ok {
		return 0, fmt.Errorf("cheat match result must be one of X, Y, Z (got %s)", cheatMatchResult)
	}

	const (
		loseMatch = "X"
		drawMatch = "Y"
		winMatch  = "Z"
	)
	var (
		responsePlay string
		score        int
		err          error
	)
	switch {
	case cheatMatchResult == loseMatch:
		responsePlay = g.opponentWinRules[opponentPlay]
		score, err = g.MatchOutcome(opponentPlay, responsePlay)
		if err != nil {
			return 0, err
		}
		return score, nil
	case cheatMatchResult == drawMatch:
		responsePlay = g.opponentDrawRules[opponentPlay]
		score, err = g.MatchOutcome(opponentPlay, responsePlay)
		if err != nil {
			return 0, err
		}
		return score, nil
	default:
		responsePlay = g.opponentLossRules[opponentPlay]
		score, err = g.MatchOutcome(opponentPlay, responsePlay)
		if err != nil {
			return 0, err
		}
		return score, nil
	}
}

func ComputeStrategyScore(strategy io.Reader) (int, error) {
	if strategy == nil {
		return 0, errors.New("strategy must point to a non-nil strategy source")
	}

	game := NewGame()
	sc := bufio.NewScanner(strategy)
	score := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		opponentPlay, responsePlay := fields[0], fields[1]
		matchScore, err := game.MatchOutcome(opponentPlay, responsePlay)
		if err != nil {
			return 0, err
		}
		score += matchScore
	}

	return score, nil
}

func ComputeCheatStrategyScore(strategy io.Reader) (int, error) {
	if strategy == nil {
		return 0, errors.New("strategy must point to a non-nil strategy source")
	}

	game := NewGame()
	sc := bufio.NewScanner(strategy)
	score := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		opponentPlay, responsePlay := fields[0], fields[1]
		matchScore, err := game.CheatMatchOutcome(opponentPlay, responsePlay)
		if err != nil {
			return 0, err
		}
		score += matchScore
	}

	return score, nil
}
