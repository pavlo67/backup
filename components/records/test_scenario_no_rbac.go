package records

import (
	"os"
	"testing"

	"github.com/pavlo67/common/common/joiner"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/logger"
)

// TODO: test .History
// TODO: test .List() with selectors

func OperatorTestScenarioNoRBAC(t *testing.T, joinerOp joiner.Operator, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	recordsOp, _ := joinerOp.Interface(InterfaceKey).(Operator)
	require.NotNil(t, recordsOp)

	cleanerOp, _ := joinerOp.Interface(InterfaceCleanerKey).(crud.Cleaner)
	require.NotNil(t, cleanerOp)

	// clear ------------------------------------------------------------------------------

	err := cleanerOp.Clean(nil)
	require.NoError(t, err)

	// prepare records to test  -----------------------------------------------------------

	// prepare records & crud.Options -----------------------------------------

	options := crud.Options{Identity: &auth.Identity{ID: authID1}}

	//// save record without identity: error ------------------------------------
	//
	//item01Saved, err := recordsOp.Save(item01, &options0)
	//require.Error(t, err)
	//require.Empty(t, item01Saved)

	// save record without ownerID: added automatically, ok -------------------

	item01 := item11
	item01.OwnerID = ""
	item01Saved, err := recordsOp.Save(item01, &options)
	require.NoError(t, err)
	require.NotEmpty(t, item01Saved)
	require.Equal(t, item01Saved.OwnerID, authID1) // options.Identity.ID

	readOkTest(t, recordsOp, *item01Saved, options)

	// ------------------------------------------------------------------------

	item22Saved := saveTestNoRBAC(t, recordsOp, item11, item12, item22, options)

	// check .Remove(), .Read(), .List(), -------------------------------------

	err = recordsOp.Remove(item22Saved.ID, &options)
	require.NoError(t, err)

	readFailTest(t, recordsOp, item22Saved.ID, options)
	readOkTest(t, recordsOp, *item01Saved, options)

}

func saveTestNoRBAC(t *testing.T, recordsOp Operator, itemToSave, itemToUpdate, itemToUpdateAgain Item, options crud.Options) Item {

	// insert ---------------------------------------------

	itemToSave.OwnerID = authID1

	itemSaved1, err := recordsOp.Save(itemToSave, &options)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved1)
	require.Equal(t, itemToSave.Content, itemSaved1.Content)

	itemToSave.ID = itemSaved1.ID

	readOkTest(t, recordsOp, itemToSave, options)

	// update ---------------------------------------------

	itemToUpdate.ID = itemToSave.ID
	itemSaved2, err := recordsOp.Save(itemToUpdate, &options)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved2)
	require.Equal(t, itemToUpdate.ID, itemSaved2.ID)
	require.Equal(t, itemToUpdate.Content, itemSaved2.Content)

	readOkTest(t, recordsOp, itemToUpdate, options)

	// prepare item to update again ------------------------------------------

	itemToUpdateAgain.ID = itemToSave.ID

	itemSaved3, err := recordsOp.Save(itemToUpdateAgain, &options)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved3)
	require.Equal(t, itemToUpdateAgain.ID, itemSaved3.ID)
	require.Equal(t, itemToUpdateAgain.Content, itemSaved3.Content)

	readOkTest(t, recordsOp, itemToUpdateAgain, options)

	return itemToSave
}
