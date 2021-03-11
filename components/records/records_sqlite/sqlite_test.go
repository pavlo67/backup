package records_sqlite

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/connect/connect_sqlite"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/records"
)

func TestCRUD(t *testing.T) {
	_, cfgService, l := apps.PrepareTests(t, "../../../apps/", "test", "records_sqlite")
	require.NotNil(t, cfgService)

	var cfg config.Access
	err := cfgService.Value("files_fs", &cfg)
	require.NoErrorf(t, err, "%#v", cfgService)

	components := []starter.Starter{
		{connect_sqlite.Starter(), nil},
		{Starter(), nil},
	}

	joinerOp, err := starter.Run(components, cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	records.OperatorTestScenarioNoRBAC(t, joinerOp, l)
}
