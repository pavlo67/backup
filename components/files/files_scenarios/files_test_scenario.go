package files_scenarios

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/tools/components/files"

)

const path1 = "aaa"
var fileData1 = []byte("fileData1")

const path2 = "bbb/ccc"
var fileData2 = []byte("fileData2")

func FilesTestScenario(t *testing.T, joinerOp joiner.Operator, interfaceKey joiner.InterfaceKey, bucketID files.BucketID) {
	filesOp, _ := joinerOp.Interface(interfaceKey).(files.Operator)
	require.NotNil(t, filesOp)

	path1Saved := saveTest(t, filesOp, bucketID, path1, fileData1)
	require.NotEmpty(t, path1Saved)

	path2Saved := saveTest(t, filesOp, bucketID, path2, fileData2)
	require.NotEmpty(t, path2Saved)
}

func saveTest(t *testing.T, filesOp files.Operator, bucketID files.BucketID, path string, data []byte) (pathCorrected string) {
	pathSaved, err := filesOp.Save(bucketID, path, "", data)
	require.NoError(t, err)
	require.NotEmpty(t, pathSaved)

	dataReaded, err := filesOp.Read(bucketID, pathSaved)
	require.NoError(t, err)
	require.Equal(t, data, dataReaded)

	return pathSaved
}
