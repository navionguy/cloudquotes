package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Author holds the name of somebody who authored a quote
type Author struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name" form:"name"`
}

// AuthorCredit allows me to find out how many quotes each author has
type AuthorCredit struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Name  string    `json:"name" db:"name"`
	Count int       `json:"count" db:"count"`
}

// AuthorCredits holds all the authors
type AuthorCredits []AuthorCredit

// String is not required by pop and may be deleted
func (a Author) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Authors is not required by pop and may be deleted
type Authors []Author

// String is not required by pop and may be deleted
func (a Authors) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Author) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
		&validators.StringLengthInRange{Field: a.Name, Name: "Name", Min: 1, Max: 255, Message: "length must be 1-255"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Author) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Author) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SelectValue returns the author ID value to a form SelectTag
func (a Author) SelectValue() interface{} {
	return a.ID.String()
}

// SelectLabel allows authors to be in a form SelectTag
func (a Author) SelectLabel() string {
	return a.Name
}

// FindByID pulls up the author record based on ID
func (a *Author) FindByID() error {

	authRecs := []Author{}
	query := DB.Where(fmt.Sprintf("id = '%s'", a.ID))
	err := query.All(&authRecs)

	if err != nil {
		return err
	}

	if len(authRecs) == 0 {
		return errors.New("author ID not found in db")
	}

	*a = authRecs[0]

	return nil
}

// FindByName check for an author by name
func (a *Author) FindByName() error {
	// name can't be empty
	if len(a.Name) == 0 {
		return errors.New("author name can't be blank")
	}

	// Break the passed name down into pieces
	parts := strings.Split(a.Name, " ")

	authRecs := []Author{}
	query := DB.Where(fmt.Sprintf("name ILIKE '%%%s%%' AND name ILIKE '%%%s%%'", parts[0], parts[len(parts)-1]))
	err := query.All(&authRecs)

	if err != nil {
		return err
	}

	if len(authRecs) == 0 {
		return errors.New("author name not found in db")
	}

	*a = authRecs[0]

	return nil
}

// Create adds a new speaker to the authors table
func (a *Author) Create() (*validate.Errors, error) {

	verrs, err := DB.ValidateAndCreate(a)

	return verrs, err
}
