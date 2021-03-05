package entity

import (
	"time"

	"gorm.io/gorm"
)

// Proposal is the entity for proposed rule-changes
type Proposal struct {
	gorm.Model

	Number      int
	Text        string
	SubmittedBy string // username of submitter

	RelatedToRule    int
	VotesInFavor     int
	VotesAgainst     int
	VotingInProgress bool
	VotingStarted    *time.Time
	VotingEnded      *time.Time
}

// IsVotingStarted determines if voting has
// started on this proposal
func (p Proposal) IsVotingStarted() bool {
	return p.VotingInProgress && p.VotingStarted != nil
}

// IsVotingEnded determines if voting has
// ended on this proposal
func (p Proposal) IsVotingEnded() bool {
	return p.VotingEnded != nil && !p.VotingInProgress
}

// StartVoting starts voting on this
// proposal
func (p *Proposal) StartVoting() {
	if p.IsVotingStarted() || p.IsVotingEnded() {
		return
	}

	//	loc, _ := time.LoadLocation("EST")
	p.VotingInProgress = true
	//	p.VotingStarted
}
