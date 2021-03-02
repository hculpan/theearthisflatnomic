package entity

import "gorm.io/gorm"

// Rule represents a single rule
type Rule struct {
	gorm.Model
	RuleNumber int
	RuleText   string
	Initial    bool
	Active     bool
}

// FindAllActiveRules returns the list of
// active rules
func FindAllActiveRules() []Rule {
	var result []Rule = []Rule{}
	db.Where(&Rule{Active: true}).Find(&result)
	return result
}

// Insert inserts the record into the db
// This does not check to see if the records
// already exists; it assumes it does not.
func (r Rule) Insert() error {
	return db.Create(&r).Error
}

// Save will insert or update the record,
// based on whether it has a primary key or
// not
func (r Rule) Save() error {
	return db.Save(&r).Error
}
