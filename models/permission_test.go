package models_test

import (
	"strings"
	"testing"

	"github.com/gobuffalo/uuid"
	"github.com/navionguy/cloudquotes/models"
	"github.com/stretchr/testify/require"
)

func Test_Permission(t *testing.T) {

	var fields = []struct {
		fn  string
		msg string
	}{
		{"id", "Permission id field not found"},
		{"name", "Permission name field not found"},
		{"created_at", "created_at field not found"},
		{"updated_at", "updated_at field not found"},
	}

	a := models.Permission{
		Name: validName,
	}

	// convert Permission to json

	js := a.String()

	// if you get nothing, that's a problem

	if len(js) == 0 {
		t.Error("unable to marshal Permission")
		t.FailNow()
	}

	// make sure expected fields are there

	rq := require.New(t)

	for _, fld := range fields {
		rq.Containsf(js, fld.fn, fld.msg)
	}

	var ar models.Permissions

	ar = append(ar, a)

	js = ar.String()

	if len(js) == 0 {
		t.Error("Unable to marshal array of Permissions")
		t.Fail()
	}
}

const invalidUUID = "563cd207-ab16-4a46-b44e-7317b96c6ba9"
const validUUID = "b39300f0-6760-4feb-bc32-4b8682b0175d" // matches entry in testPermissions.toml
const validName = "RuleTheWorld"

// Test for finding an existing Permission
func (ms *ModelSuite) Test_Permission_FindByID() {
	ms.LoadFixture("test Permissions")

	id, err := uuid.FromString(validUUID)

	if err != nil {
		ms.Fail("uuid.FromString failed", err.Error())
	}

	auth := models.Permission{
		ID: id,
	}

	pauth, err := auth.FindByID()

	if err != nil {
		ms.Fail("FindByID failed", err.Error())
	}

	if pauth == nil {
		ms.Fail("FindByID failed", "validUUID was not found in database")
	}

	if strings.Compare(pauth.Name, validName) != 0 {
		ms.Fail("FindByID didn't find expected Permission", pauth.Name)
	}

	// as long as I have a valid Permission, check some other functions

	if strings.Compare(pauth.SelectLabel(), validName) != 0 {
		ms.Fail("unexpected SelectLabel", pauth.SelectLabel())
	}

	v := pauth.SelectValue()
	s, ok := v.(string)

	if ok {
		if strings.Compare(s, validUUID) != 0 {
			ms.Fail("unexpected SelectValue", s)
		}
	} else {
		ms.Fail("unexpected SelectValue", "not a string")
	}
}

// test that FindByID correctly handles NOT finding the Permission
func (ms *ModelSuite) Test_Permission_FindByID_BadID() {
	ms.LoadFixture("test Permissions")

	id, err := uuid.FromString(invalidUUID)

	if err != nil {
		ms.Fail("uuid.FromString failed", err.Error())
	}

	auth := models.Permission{
		ID: id,
	}

	pauth, err := auth.FindByID()

	if err != nil {
		ms.Fail("FindByID failed", err.Error())
	}

	if pauth != nil {
		ms.Fail("FindByID succeeded with an invalid UUID", pauth.Name)
	}
}

func (ms *ModelSuite) Test_Permission_Create() {
	ms.LoadFixture("test Permissions")

	auth := models.Permission{
		Name: "Brand New Permission",
	}

	verrs, err := ms.DB.ValidateAndCreate(&auth)

	if err != nil {
		ms.Fail("failed to create Permission", err.Error())
	}

	if verrs.HasAny() {
		ms.Fail("unable to validate new Permission", verrs.String())
	}
}

func (ms *ModelSuite) Test_Permission_CreateInvalid() {
	ms.LoadFixture("test Permissions")

	auth := models.Permission{}

	verrs, err := ms.DB.ValidateAndCreate(&auth)

	if err != nil {
		ms.Fail("failed to create Permission", err.Error())
	}

	if !verrs.HasAny() {
		ms.Fail("invalid Permission validated", "no name supplied")
	}
}
