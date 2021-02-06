package files_fs

import (
	"testing"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/files"
)

func TestFilesFS(t *testing.T) {
	_, cfgService := apps.PrepareTests(t, "test_service", "../../../apps/", "test")
	require.NotNil(t, cfgService)

	var cfg config.Access
	err := cfgService.Value("files_fs", &cfg)
	require.NoErrorf(t, err, "%#v", cfgService)

	bucketID := files.BucketID("test_bucket")
	components := []starter.Starter{
		{Starter(), common.Map{"buckets": Buckets{bucketID: cfg.Path}}},
	}

	joinerOp, err := starter.Run(components, cfgService, "CLI BUILD FOR TEST")
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	files.FilesTestScenario(t, joinerOp, files.InterfaceKey, files.InterfaceKeyCleaner, bucketID)
}
