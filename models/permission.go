package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Permission holds a permission granted to a user
type Permission struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Name      string    `json:"name" db:"name"`

	// Relationships
	User User `belongs_to:"user" db:"-"`

	// Foreign keys
	UserID uuid.UUID `json:"user_id" db:"user_id"`
}

// String is not required by pop and may be deleted
func (a Permission) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Permissions is not required by pop and may be deleted
type Permissions []Permission

// String is not required by pop and may be deleted
func (a Permissions) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Permission) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Permission) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Permission) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// SelectValue returns the Permission ID value to a form SelectTag
func (a Permission) SelectValue() interface{} {
	return a.ID.String()
}

// SelectLabel allows Permissions to be in a form SelectTag
func (a Permission) SelectLabel() string {
	return a.Name
}

// FindByID pulls up the Permission record based on ID
func (a *Permission) FindByID() error {

	permRecs := []Permission{}
	query := DB.Where(fmt.Sprintf("id = '%s'", a.ID))
	err := query.All(&permRecs)

	if err != nil {
		return err
	}

	if len(permRecs) == 0 {
		return errors.New("Permission ID not found in db")
	}

	*a = permRecs[0]

	return nil
}
