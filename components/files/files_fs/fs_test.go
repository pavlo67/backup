package files_fs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/apps"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/tools/components/files"
)

func TestFilesFS(t *testing.T) {
	_, cfgService, l := apps.PrepareTests(t, "../../../apps/_environments/", "test", "files_fs.log")
	require.NotNil(t, cfgService)

	var cfg config.Access
	err := cfgService.Value("files_fs", &cfg)
	require.NoErrorf(t, err, "%#v", cfgService)

	bucketID := files.BucketID("test_bucket")
	components := []starter.Starter{
		{Starter(), common.Map{"buckets": files.Buckets{bucketID: cfg.Path}}},
	}

	joinerOp, err := starter.Run(components, cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	files.FilesTestScenario(t, joinerOp, files.InterfaceKey, files.InterfaceKeyCleaner, bucketID)
}
