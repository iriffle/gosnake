package main

import (
	"encoding/json"
)

type Score struct {
	Mode   string
	Player string
	Score  int
	Pos    int
}

func NewScore(mode, player string, score, pos int) *Score {
	s := Score{
		Mode:   mode,
		Player: player,
		Score:  score,
		Pos:    pos,
	}
	return &s
}

func (s *Score) Encode() ([]byte, error) {
	var jsonData []byte
	jsonData, err := json.Marshal(s)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func (s *Score) Decode(jsonData []byte) error {
	err := json.Unmarshal(jsonData, s)
	if err != nil {
		return err
	}
	return nil
}

func CheckScore(scores []*Score, score int) int {
	var pos int = -1
	for i, _ := range scores {
		if score > scores[i].Score {
			pos = i
		}
	}
	return pos
}

func (s *Score) AddScore(scores []*Score, pos int) []*Score {
	scores[pos] = s
	return scores
}

func ShiftScores(scores []*Score, pos, max int) []*Score {
	for i := pos; i < len(scores); i++ {
		if i < max {
			scores[i].Pos += 1
		} else {
			scores[i] = scores[len(scores)-1]
			scores[len(scores)-1] = nil
			scores = scores[:len(scores)-1]
		}
	}
	return scores
}

func UpdateScores(scores []*Score, mode, player string, score, max int) []*Score {
	if score > 0 {
		if len(scores) <= max {
			pos := CheckScore(scores, score)
			if pos != -1 {
				s := NewScore(mode, player, score, pos)
				scores = ShiftScores(scores, pos, max)
				scores = s.AddScore(scores, pos)
			} else if len(scores) < max {
				s := NewScore(mode, player, score, pos)
				scores = s.AddScore(scores, pos)
			}
		}
	}
	return scores
}
