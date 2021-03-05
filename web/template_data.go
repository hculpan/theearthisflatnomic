package web

import (
	"fmt"
	"time"

	"github.com/hculpan/theearthisflatnomic/entity"
)

// TemplateData is the minimum data structure
// to send to all pages
type TemplateData struct {
	Cookies         map[string]string
	ErrorText       string
	Data            interface{}
	MessageText     string
	DestinationURL  string
	UserDisplayName string
	IsTurn          bool
	EndOfTurn       time.Time
	Action          entity.NextTurnAction

	OriginalURL     string
	OriginalURLName string
}

// IsSubmitPhase return true if player's phase
// is to submit a proprosal
func (t TemplateData) IsSubmitPhase() bool {
	return t.Action == entity.Submit
}

// IsVotingPhase returns true if player's phase
// is voting on their proposal
func (t TemplateData) IsVotingPhase() bool {
	return t.Action == entity.Voting
}

// DisplayEndOfTurn returns the EndOfTurn as a
// readable string
func (t TemplateData) DisplayEndOfTurn() string {
	return t.EndOfTurn.Format("Monday, Jan 2, 2006 3:04pm EST")
}

// DisplayActionMessage returns the action as a string
func (t TemplateData) DisplayActionMessage() string {
	switch t.Action {
	case entity.Submit:
		return fmt.Sprintf("It is your turn.  You have until %s to propose a rule-change.", t.DisplayEndOfTurn())
	case entity.Voting:
		return fmt.Sprintf("Voting on your proposal will end on %s.", t.DisplayEndOfTurn())
	default:
		return "** Uknown action **"
	}
}

// DisplaySecondaryActionMessage returns a sub-message for the action
func (t TemplateData) DisplaySecondaryActionMessage() string {
	switch t.Action {
	case entity.Submit:
		return "If you do not submit a proprosal by then, your turn will end automatically."
	case entity.Voting:
		return ""
	default:
		return "** Uknown action **"
	}
}
