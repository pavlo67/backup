package records

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/logger"
)

// TODO: test .History
// TODO: test .List() with selectors

func OperatorTestScenario(t *testing.T, recordsOp Operator, cleanerOp db.Cleaner, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	// clear ------------------------------------------------------------------------------

	err := cleanerOp.Clean(nil)
	require.NoError(t, err)

	// prepare records to test  -----------------------------------------------------------

	// prepare records & auth.Identity -----------------------------------------

	item01 := item11
	item01.OwnerNSS = ""

	identity0 := auth.Identity{}
	identity1 := auth.Identity{ID: authID1}

	// save record without identity1: error ------------------------------------

	item01SavedID, err := recordsOp.Save(item01, &identity0)
	require.Error(t, err)
	require.Empty(t, item01SavedID)

	// save record without OwnerNSS: added automatically, ok -------------------

	require.Empty(t, item01.OwnerNSS)
	item01SavedID, err = recordsOp.Save(item01, &identity1)
	require.NoError(t, err)

	// TODO!!!
	// require.Equal(t, item01SavedID.OwnerNSS, identity1.Identity.ID)

	saveTest(t, recordsOp, item11, item12, item22) // item22Saved :=
	saveTest(t, recordsOp, item11, item12, item22) // item22SavedAgain :=

	// check .Remove(), .Read(), .List(), -------------------------------------

	//owner22Identity := auth.Identity{ID: item22Saved.OwnerNSS}
	//owner22ViewerIdentity := auth.Identity{ID: item22Saved.ViewerNSS}
	//
	//err = recordsOp.Remove(item22Saved.ID, &owner22Identity)
	//require.NoError(t, err)
	//
	//readFailTest(t, recordsOp, item22Saved.ID, owner22Identity)
	//readFailTest(t, recordsOp, item22Saved.ID, owner22ViewerIdentity)
	//readOkTest(t, recordsOp, item22SavedAgain, owner22Identity)
	//readOkTest(t, recordsOp, item22SavedAgain, owner22ViewerIdentity)

}

func saveTest(t *testing.T, recordsOp Operator, itemToSave, itemToUpdate, itemToUpdateAgain Item) Item {
	//
	//identity1 := auth.Identity{ID: authID1}
	//identity2 := auth.Identity{ID: authID2}
	//identity3 := auth.Identity{ID: authID3}
	//
	//// prepare item to save --------------------------------------------------
	//
	//itemToSave.OwnerNSS = authID1
	//itemToSave.ViewerNSS = authID1
	//
	//// check .Save() with other identity: error -------------------------------
	//
	//itemSaved, err := recordsOp.Save(itemToSave, &identity2)
	//require.Error(t, err)
	//require.Empty(t, itemSaved)
	//
	//// check .Save() with owner identity: ok ----------------------------------
	//
	//itemSaved, err = recordsOp.Save(itemToSave, &identity1)
	//require.NoError(t, err)
	//require.NotEmpty(t, itemSaved)
	//require.Equal(t, itemToSave.Content, itemSaved.Content)
	//require.Equal(t, itemToSave.OwnerNSS, itemSaved.OwnerNSS)
	//require.Equal(t, itemToSave.ViewerNSS, itemSaved.ViewerNSS)
	//
	//itemToSave.ID = itemSaved.ID
	//
	//// check .Read(), .List() with owner/viewer identity ----------------------
	//
	//readOkTest(t, recordsOp, itemToSave, identity1)
	//
	//// check .Read(), .List() with other identity -----------------------------
	//
	//readFailTest(t, recordsOp, itemToSave.ID, identity2)
	//
	//// prepare item to update ------------------------------------------------
	//
	//itemToUpdate.ID = itemToSave.ID
	//itemToUpdate.OwnerNSS = authID1
	//itemToUpdate.ViewerNSS = authID2
	//
	//// check updating .Save() with other identity: error ----------------------
	//
	//itemSaved, err = recordsOp.Save(itemToUpdate, &identity2)
	//require.Error(t, err)
	//require.Empty(t, itemSaved)
	//
	//// check updating .Save() with owner identity: ok -------------------------
	//
	//itemSaved, err = recordsOp.Save(itemToUpdate, &identity1)
	//require.NoError(t, err)
	//require.NotEmpty(t, itemSaved)
	//require.Equal(t, itemToUpdate.ID, itemSaved.ID)
	//require.Equal(t, itemToUpdate.Content, itemSaved.Content)
	//require.Equal(t, itemToUpdate.OwnerNSS, itemSaved.OwnerNSS)
	//require.Equal(t, itemToUpdate.ViewerNSS, itemSaved.ViewerNSS)
	//
	//// check .Read(), .List() with owner identity -----------------------------
	//
	//readOkTest(t, recordsOp, itemToUpdate, identity1)
	//
	//// check .Read(), .List() with viewer identity ----------------------------
	//
	//readOkTest(t, recordsOp, itemToUpdate, identity2)
	//
	//// check .Read(), .List() with other identity -----------------------------
	//
	//readFailTest(t, recordsOp, itemToUpdate.ID, identity3)
	//
	//// prepare item to update again ------------------------------------------
	//
	//itemToUpdateAgain.ID = itemToSave.ID
	//itemToUpdateAgain.OwnerNSS = authID2
	//itemToUpdateAgain.ViewerNSS = authID2
	//
	//// check updating .Save() with other identity: error ----------------------
	//
	//itemSaved, err = recordsOp.Save(itemToUpdateAgain, &identity2)
	//require.Error(t, err)
	//require.Empty(t, itemSaved)
	//
	//// check updating .Save() with owner identity: ok -------------------------
	//
	//itemSaved, err = recordsOp.Save(itemToUpdateAgain, &identity1)
	//require.NoError(t, err)
	//require.NotEmpty(t, itemSaved)
	//require.Equal(t, itemToUpdateAgain.ID, itemSaved.ID)
	//require.Equal(t, itemToUpdateAgain.Content, itemSaved.Content)
	//require.Equal(t, itemToUpdateAgain.OwnerNSS, itemSaved.OwnerNSS)
	//require.Equal(t, itemToUpdateAgain.ViewerNSS, itemSaved.ViewerNSS)
	//
	//// check .Read(), .List() with owner/viewer identity ----------------------
	//
	//readOkTest(t, recordsOp, itemToUpdateAgain, identity2)
	//
	//// check .Read(), .List() with other identity -----------------------------
	//
	//readFailTest(t, recordsOp, itemToUpdateAgain.ID, identity1)
	//

	return itemToSave
}
