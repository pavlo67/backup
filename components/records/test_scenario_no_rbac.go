package records

import (
	"os"
	"testing"

	"github.com/pavlo67/data_exchange/components/ns"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/selectors"
)

// TODO: test .History
// TODO: test .List() with selectors

func OperatorTestScenarioNoRBAC(t *testing.T, joinerOp joiner.Operator, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	recordsOp, _ := joinerOp.Interface(InterfaceKey).(Operator)
	require.NotNil(t, recordsOp)

	cleanerOp, _ := joinerOp.Interface(InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, cleanerOp)

	// clear ------------------------------------------------------------------------------

	err := cleanerOp.Clean(nil)
	require.NoError(t, err)

	// prepare records to test  -----------------------------------------------------------

	// prepare records & auth.Identity -----------------------------------------

	identity := auth.Identity{ID: authID1}

	//// save record without identity: error ------------------------------------
	//
	//item01Saved, err := recordsOp.Save(item01, &options0)
	//require.Error(t, err)
	//require.Empty(t, item01Saved)

	// save record without OwnerNSS: added automatically, ok -------------------

	item01 := item11
	item01.OwnerNSS = ""
	item01SavedID, err := recordsOp.Save(item01, &identity)
	require.NoError(t, err)
	require.NotEmpty(t, item01SavedID)
	require.Equal(t, item01SavedID, authID1)

	item01Saved := item01
	item01Saved.ID = item01SavedID

	readOkTest(t, recordsOp, item01Saved, identity)

	// ------------------------------------------------------------------------

	item22Saved := dbTestNoRBAC(t, recordsOp, item11, item12, item22, identity)

	// check .Remove(), .Read(), .List(), -------------------------------------

	err = recordsOp.Remove(item22Saved.ID, &identity)
	require.NoError(t, err)

	readFailTest(t, recordsOp, item22Saved.ID, identity)
	readOkTest(t, recordsOp, item01Saved, identity)

}

func dbTestNoRBAC(t *testing.T, recordsOp Operator, itemToSave, itemToUpdate, itemToUpdateAgain Item, identity auth.Identity) Item {

	// prepare selector tagged ----------------------------

	selectorTagged := selectors.Term{
		Key:    HasTag,
		Values: []string{testNewTag},
	}

	// insert ---------------------------------------------

	itemToSave.OwnerNSS = ns.NSS(authID1)

	itemSaved1ID, err := recordsOp.Save(itemToSave, &identity)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved1ID)

	// TODO!!!
	// require.Equal(t, itemToSave.Content, itemSaved1.Content)

	itemToSave.ID = itemSaved1ID

	// prepare selector parent ----------------------------

	selectorParent := selectors.Term{
		Key:    HasParent,
		Values: []string{string(itemSaved1ID)},
	}

	// check inserted -------------------------------------

	readOkTest(t, recordsOp, itemToSave, identity)

	items, err := recordsOp.List(&selectorTagged, &identity)
	require.NoError(t, err)
	require.Equal(t, 0, len(items))

	items, err = recordsOp.List(&selectorParent, &identity)
	require.NoError(t, err)
	require.Equal(t, 0, len(items))

	// update ---------------------------------------------

	itemToUpdate.ID = itemToSave.ID
	itemToUpdate.Tags, err = recordsOp.AddParent(append(itemToUpdate.Tags, testNewTag), itemToSave.ID)
	require.NoError(t, err)

	itemSaved2ID, err := recordsOp.Save(itemToUpdate, &identity)
	require.NoError(t, err)
	require.Equal(t, itemToUpdate.ID, itemSaved2ID)

	// TODO!!!
	// require.Equal(t, itemToUpdate.Content, itemSaved2.Content)

	readOkTest(t, recordsOp, itemToUpdate, identity)

	items, err = recordsOp.List(&selectorTagged, &identity)
	require.NoError(t, err)
	require.Equal(t, 1, len(items))
	require.Equal(t, items[0].ID, itemToSave.ID)

	items, err = recordsOp.List(&selectorParent, &identity)
	require.NoError(t, err)
	require.Equal(t, 1, len(items))
	require.Equal(t, items[0].ID, itemToSave.ID)

	// prepare item to update again ------------------------------------------

	itemToUpdateAgain.ID = itemToSave.ID

	itemSaved3ID, err := recordsOp.Save(itemToUpdateAgain, &identity)
	require.NoError(t, err)
	require.Equal(t, itemToUpdateAgain.ID, itemSaved3ID)

	// TODO!!!
	//require.Equal(t, itemToUpdateAgain.Content, itemSaved3.Content)

	readOkTest(t, recordsOp, itemToUpdateAgain, identity)

	return itemToSave
}
