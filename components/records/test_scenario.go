package records

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/logger"
)

// TODO: test .History
// TODO: test .List() with selectors

func OperatorTestScenario(t *testing.T, recordsOp Operator, cleanerOp crud.Cleaner, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	// clear ------------------------------------------------------------------------------

	err := cleanerOp.Clean(nil)
	require.NoError(t, err)

	// prepare records to test  -----------------------------------------------------------

	// prepare records & crud.Options -----------------------------------------

	item01 := item11
	item01.OwnerID = ""

	options0 := crud.Options{}
	options1 := crud.Options{Identity: &auth.Identity{ID: authID1}}

	// save record without identity: error ------------------------------------

	item01Saved, err := recordsOp.Save(item01, &options0)
	require.Error(t, err)
	require.Empty(t, item01Saved)

	// save record without ownerID: added automatically, ok -------------------

	require.Empty(t, item01.OwnerID)
	item01Saved, err = recordsOp.Save(item01, &options1)
	require.NoError(t, err)
	require.NotEmpty(t, item01Saved)
	require.Equal(t, item01Saved.OwnerID, options1.Identity.ID)

	item22Saved := saveTest(t, recordsOp, item11, item12, item22)
	item22SavedAgain := saveTest(t, recordsOp, item11, item12, item22)

	// check .Remove(), .Read(), .List(), -------------------------------------

	owner22OwnerOptions := crud.Options{Identity: &auth.Identity{ID: item22Saved.OwnerID}}
	owner22ViewerOptions := crud.Options{Identity: &auth.Identity{ID: item22Saved.ViewerID}}

	err = recordsOp.Remove(item22Saved.ID, &owner22OwnerOptions)
	require.NoError(t, err)

	readFailTest(t, recordsOp, item22Saved.ID, owner22OwnerOptions)
	readFailTest(t, recordsOp, item22Saved.ID, owner22ViewerOptions)
	readOkTest(t, recordsOp, item22SavedAgain, owner22OwnerOptions)
	readOkTest(t, recordsOp, item22SavedAgain, owner22ViewerOptions)

}

func saveTest(t *testing.T, recordsOp Operator, itemToSave, itemToUpdate, itemToUpdateAgain Item) Item {

	options1 := crud.Options{Identity: &auth.Identity{ID: authID1}}
	options2 := crud.Options{Identity: &auth.Identity{ID: authID2}}
	options3 := crud.Options{Identity: &auth.Identity{ID: authID3}}

	// prepare item to save --------------------------------------------------

	itemToSave.OwnerID = authID1
	itemToSave.ViewerID = authID1

	// check .Save() with other options: error -------------------------------

	itemSaved, err := recordsOp.Save(itemToSave, &options2)
	require.Error(t, err)
	require.Empty(t, itemSaved)

	// check .Save() with owner options: ok ----------------------------------

	itemSaved, err = recordsOp.Save(itemToSave, &options1)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved)
	require.Equal(t, itemToSave.Content, itemSaved.Content)
	require.Equal(t, itemToSave.OwnerID, itemSaved.OwnerID)
	require.Equal(t, itemToSave.ViewerID, itemSaved.ViewerID)

	itemToSave.ID = itemSaved.ID

	// check .Read(), .List() with owner/viewer options ----------------------

	readOkTest(t, recordsOp, itemToSave, options1)

	// check .Read(), .List() with other options -----------------------------

	readFailTest(t, recordsOp, itemToSave.ID, options2)

	// prepare item to update ------------------------------------------------

	itemToUpdate.ID = itemToSave.ID
	itemToUpdate.OwnerID = authID1
	itemToUpdate.ViewerID = authID2

	// check updating .Save() with other options: error ----------------------

	itemSaved, err = recordsOp.Save(itemToUpdate, &options2)
	require.Error(t, err)
	require.Empty(t, itemSaved)

	// check updating .Save() with owner options: ok -------------------------

	itemSaved, err = recordsOp.Save(itemToUpdate, &options1)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved)
	require.Equal(t, itemToUpdate.ID, itemSaved.ID)
	require.Equal(t, itemToUpdate.Content, itemSaved.Content)
	require.Equal(t, itemToUpdate.OwnerID, itemSaved.OwnerID)
	require.Equal(t, itemToUpdate.ViewerID, itemSaved.ViewerID)

	// check .Read(), .List() with owner options -----------------------------

	readOkTest(t, recordsOp, itemToUpdate, options1)

	// check .Read(), .List() with viewer options ----------------------------

	readOkTest(t, recordsOp, itemToUpdate, options2)

	// check .Read(), .List() with other options -----------------------------

	readFailTest(t, recordsOp, itemToUpdate.ID, options3)

	// prepare item to update again ------------------------------------------

	itemToUpdateAgain.ID = itemToSave.ID
	itemToUpdateAgain.OwnerID = authID2
	itemToUpdateAgain.ViewerID = authID2

	// check updating .Save() with other options: error ----------------------

	itemSaved, err = recordsOp.Save(itemToUpdateAgain, &options2)
	require.Error(t, err)
	require.Empty(t, itemSaved)

	// check updating .Save() with owner options: ok -------------------------

	itemSaved, err = recordsOp.Save(itemToUpdateAgain, &options1)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved)
	require.Equal(t, itemToUpdateAgain.ID, itemSaved.ID)
	require.Equal(t, itemToUpdateAgain.Content, itemSaved.Content)
	require.Equal(t, itemToUpdateAgain.OwnerID, itemSaved.OwnerID)
	require.Equal(t, itemToUpdateAgain.ViewerID, itemSaved.ViewerID)

	// check .Read(), .List() with owner/viewer options ----------------------

	readOkTest(t, recordsOp, itemToUpdateAgain, options2)

	// check .Read(), .List() with other options -----------------------------

	readFailTest(t, recordsOp, itemToUpdateAgain.ID, options1)

	return itemToSave
}
