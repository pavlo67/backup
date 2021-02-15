package notebook_http_tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/tools/apps/nb/nb_settings"
)

func TestNotebookREST(t *testing.T) {
	_, cfgService, l := apps.PrepareTests(t, "test_service", "../../../apps/", "test", "records_sqlite")
	require.NotNil(t, cfgService)

	var cfg config.Access
	err := cfgService.Value("files_fs", &cfg)
	require.NoErrorf(t, err, "%#v", cfgService)

	components, err := nb_settings.ClientComponents()
	require.NoError(t, err)

	joinerOp, err := starter.Run(components, cfgService, "HTTP/CLI BUILD FOR TESTS", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	// records.OperatorTestScenarioNoRBAC(t, joinerOp, l)
}
