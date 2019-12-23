package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

// Conversation Common element of one or more quotes
type Conversation struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	OccurredOn time.Time `json:"occurredon" db:"occurredon"`
	Publish    bool      `json:"publish" db:"publish"`

	// Relationships
	Quotes Quotes `has_many:"quotes" orderby:"sequence" db:"-"`
}

// String is not required by pop and may be deleted
func (c Conversation) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Conversations is not required by pop and may be deleted
type Conversations []Conversation

// String is not required by pop and may be deleted
func (c Conversations) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Conversation) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.TimeIsPresent{Field: c.OccurredOn, Name: "SaidOn"},
		&validators.TimeIsBeforeTime{FirstTime: c.OccurredOn, SecondTime: time.Now().AddDate(0, 0, 1), FirstName: "Said on", SecondName: "Tomorrow"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Conversation) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Conversation) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

const tempError string = "NoErr"

// Create creates a new conversation.
func (c *Conversation) Create() (*validate.Errors, error) {
	var verrs *validate.Errors

	// start a transaction for the whole conversation
	err := DB.Transaction(func(db *pop.Connection) error {
		var err error

		// create the conversation record
		verrs, err = db.ValidateAndCreate(c)

		if err != nil {
			return err
		}

		if verrs.HasAny() {
			return errors.New(tempError) // force rollback of the transaction
		}

		// loop through all the quotes and add them
		for i, quote := range c.Quotes {
			quote.Sequence = i

			verrs, err = quote.Create(db, c.ID)
			if err != nil {
				return err
			}

			if verrs.HasAny() {
				return errors.New(tempError) // this is just to get pop to rollback the transaction
			}
		}

		return nil
	})

	if err != nil {
		if strings.Compare(tempError, err.Error()) == 0 {
			return verrs, nil
		}
	}

	if err != nil {
		return nil, err
	}

	// if nothing went south, all is good
	return verrs, nil
}

// Update re-saves an already created conversation
func (c *Conversation) Update() (*validate.Errors, error) {
	var verrs *validate.Errors

	// start a transaction for the whole conversation
	err := DB.Transaction(func(db *pop.Connection) error {
		var err error

		// update the conversation record
		verrs, err = db.ValidateAndUpdate(c)

		if err != nil {
			return err
		}

		if verrs.HasAny() {
			return errors.New(tempError) // force rollback of the transaction
		}

		// loop through all the quotes and add them
		for i, quote := range c.Quotes {
			quote.Sequence = i

			if c.ID.String() == "" {
				verrs, err = quote.Create(db, c.ID)
			} else {
				verrs, err = quote.Create(db, c.ID)
			}
			if err != nil {
				return err
			}

			if verrs.HasAny() {
				return errors.New(tempError) // this is just to get pop to rollback the transaction
			}
		}

		return nil
	})

	if err != nil {
		if strings.Compare(tempError, err.Error()) == 0 {
			return verrs, nil
		}
	}

	if err != nil {
		return nil, err
	}

	// if nothing went south, all is good
	return verrs, nil
}

// MarshalConversation the passed conversation
// I convert it to JSON and Set into the form context
func (c *Conversation) MarshalConversation() (string, error) {

	cvjson, err := json.Marshal(c)

	if err != nil {
		return "", err
	}

	cvsjson := string(cvjson)
	fmt.Println(cvsjson)
	fmt.Println("starting")

	//var ccv bytes.Buffer
	//err = json.Compact(&ccv, cvjson)

	uscv := url.PathEscape(string(cvjson))
	fmt.Println(uscv)

	return uscv, nil
}

// UnmarshalConversation  pull quote out of the form
func (c *Conversation) UnmarshalConversation(cvjson string) error {
	fmt.Println("raw")
	fmt.Println(cvjson)
	fmt.Println("data")

	ccv, err := url.PathUnescape(cvjson)

	if err != nil {
		return err
	}

	ccv2, _ := url.QueryUnescape(cvjson)

	fmt.Println("results")
	fmt.Println(ccv)
	fmt.Printf("alternate\n%s\n", ccv2)
	err = json.Unmarshal([]byte(ccv), c)
	if err != nil {
		fmt.Printf("json unmarshall error %s/n", err.Error())
		return err
	}
	fmt.Println(c)
	fmt.Println(len(c.Quotes))
	fmt.Println("final")

	return nil
}
