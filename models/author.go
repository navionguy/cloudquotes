package models

import (
	"encoding/json"
	"errors"
	"fmt"
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
	Name      string    `json:"name" db:"name"`
}

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
