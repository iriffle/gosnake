package game

import (
	"encoding/json"
	"github.com/google/logger"
	"io/ioutil"
	"os"
	"sort"
)

// Score is a struct that keeps track of a specific high score
// and its corresponding player and mode.
type Score struct {
	Name  string `json:"name"`
	Mode  int    `json:"mode"`
	Score int    `json:"score"`
}

// NewScore creates a new Score struct with provided values.
func NewScore(name string, mode, score int) *Score {
	s := Score{
		Name:  name,
		Mode:  mode,
		Score: score,
	}
	return &s
}

// AddScore appends a score to a slice of Scores.
func AddScore(scores []*Score, newScore *Score) []*Score {
	scores = append(scores, newScore)
	return scores
}

// RemoveScores removes all Scores from a slice that are past the
// max size of the slice. Usually this value is taken from the
// MaxHighScores variable to determine the maximum number of high
// scores that should be kept.
func RemoveScores(scores []*Score, max int) []*Score {
	for i := (len(scores) - 1); len(scores) > max; i-- {
		scores[i] = nil
		scores = scores[:len(scores)-1]
	}
	return scores
}

// SortScores takes a slice of Scores and sorts them based on the
// actual score int value of each Score. Sorted high to low.
func SortScores(scores []*Score) []*Score {
	sort.Slice(scores, func(i, j int) bool { return scores[i].Score > scores[j].Score })
	return scores
}

// CheckScore checks to see if a particular score is higher than
// any of the scores in a slice of Scores.
func CheckScore(scores []*Score, pScore int) bool {
	for i := range scores {
		if pScore > scores[i].Score {
			return true
		}
	}
	return false
}

// UpdateScores is used to determine if a new score is in fact a high score
// that should be placed in a high scores slice. If it is greater than any
// score currently in the slice, or if the slice is not at max size, it will
// add the score, resort the slice, and then remove any extra scores that make
// the slice bigger than the max size.
func UpdateScores(scores []*Score, pName string, pScore, mode, max int) ([]*Score, bool) {
	scoreChange := false

	// Used to debug scoring issues.
	listScores := GetScores(scores)
	logger.Infof("Current scores: %v", listScores)

	// Ignore the score if it is 0
	if pScore > 0 {
		if scores != nil {

			// Check to see if the score is greater than any current scores
			scoreChange = CheckScore(scores, pScore)

			// If the score is greater or the current number of high scores is
			// less than the max number then add the score to the list.
			if scoreChange || len(scores) < max {

				// This is set again in case a score was not greater than any current score
				// but the length of the scores list is smaller than the max.
				scoreChange = true

				logger.Infof("Scores list changed: Name: %v - Score: %v - Mode: %v", pName, pScore, mode)

				// Create the new Score and append it to the scores list
				s := NewScore(pName, mode, pScore)
				scores = AddScore(scores, s)

				// Clean up the scores list if longer than max
				if len(scores) > max {
					scores = RemoveScores(scores, max)
				}

				// Sort the scores so that they are in order highest to lowest
				listScores = GetScores(scores)
				logger.Infof("New scores: %v", listScores)
			}
			// If the current list of high scores is empty then add the score skipping
			// any score checking.
		} else {
			s := NewScore(pName, mode, pScore)
			scores = AddScore(scores, s)
			scoreChange = true
			listScores = GetScores(scores)
			logger.Infof("Scores Empty. New scores: %v", listScores)
		}
	}
	return scores, scoreChange
}

// DecodeScores takes a byteValue from a JSON file and converts it into
// Score structs and assigns to proper slices
func DecodeScores(byteValue []byte) ([]*Score, []*Score) {
	var scores []*Score
	var newScores1 []*Score
	var newScores2 []*Score

	// Read the JSON byteData into Score structs
	json.Unmarshal(byteValue, &scores)
	listScores := GetScores(scores)
	logger.Infof("Decode scores: %v", listScores)

	// Separate Scores into the proper slices
	if scores != nil {
		for i := range scores {
			if scores[i].Mode == Player1 {
				newScores1 = append(newScores1, scores[i])
			} else if scores[i].Mode == Player2 {
				newScores2 = append(newScores2, scores[i])
				listScores = GetScores(newScores2)
				logger.Infof("newScores2 scores: %v", listScores)
			}
		}
	}
	return newScores1, newScores2
}

// EncodeScores takes the two score slices, combines them into one
// slice to be marshalled back into byteData and returned to be
// written to a file.
func EncodeScores(newScores1, newScores2 []*Score) []byte {
	var scores []*Score

	newScores1 = SortScores(newScores1)
	newScores2 = SortScores(newScores2)
	for _, score := range newScores1 {
		scores = append(scores, score)
	}
	for _, score := range newScores2 {
		scores = append(scores, score)
	}
	file, _ := json.MarshalIndent(scores, "", " ")
	return file
}

// WriteScores opens a JSON file and writes scores to it.
func WriteScores(newScores1, newScores2 []*Score, file string) {
	f, err := os.OpenFile(file, os.O_CREATE, 0660)
	if err != nil {
		logger.Errorf("Error creating new file: %v", err)
	}

	if err = f.Close(); err != nil {
		logger.Errorf("Error closing file: %v", err)
	}

	data := EncodeScores(newScores1, newScores2)
	_ = ioutil.WriteFile(file, data, 0644)
}

// GetScores is used to make a list of all the current scores to view
// for debugging purposes
func GetScores(scores []*Score) []int {
	var listScores []int

	for i := range scores {
		listScores = append(listScores, scores[i].Score)
	}
	return listScores
}
