package files

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/joiner"
)

const path1 = "bbb/ccc"

var fileData1 = []byte("fileData1")

const path2 = "aaa"

var fileData2 = []byte("fileData2")

func FilesTestScenario(t *testing.T, joinerOp joiner.Operator, interfaceKey joiner.InterfaceKey, bucketID BucketID) {
	filesOp, _ := joinerOp.Interface(interfaceKey).(Operator)
	require.NotNil(t, filesOp)

	path1Saved := saveTest(t, filesOp, bucketID, path1, fileData1)
	require.NotEmpty(t, path1Saved)

	path2Saved := saveTest(t, filesOp, bucketID, path2, fileData2)
	require.NotEmpty(t, path2Saved)
}

func saveTest(t *testing.T, filesOp Operator, bucketID BucketID, path string, data []byte) (pathCorrected string) {

	// check original path info ---------------------------------------------

	fi, err := filesOp.Stat(bucketID, filepath.Dir(path))
	require.NoError(t, err)
	require.NotNil(t, fi)
	require.True(t, fi.IsDir)
	size0 := fi.Size

	// save file ------------------------------------------------------------

	pathSaved, err := filesOp.Save(bucketID, path, "", data)
	require.NoError(t, err)
	require.NotEmpty(t, pathSaved)

	// check .Read(), .List(), .Stat() --------------------------------------

	dataReaded, err := filesOp.Read(bucketID, pathSaved)
	require.NoError(t, err)
	require.Equal(t, data, dataReaded)

	fis, err := filesOp.List(bucketID, filepath.Dir(pathSaved), 0)
	require.NoError(t, err)

	found := false
	for _, fi := range fis {
		if fi.Name == filepath.Base(pathSaved) {
			found = true
			require.Equalf(t, len(data), fi.Size, "%#v", fi)
		}
	}
	require.Truef(t, found, "%#v", fis)

	fi, err = filesOp.Stat(bucketID, filepath.Dir(pathSaved))
	require.NoError(t, err)
	require.NotNil(t, fi)
	require.True(t, fi.IsDir)
	require.Equalf(t, size0+int64(len(data)), fi.Size, "%#v", fi)

	// remove file ----------------------------------------------------------

	err = filesOp.Remove(bucketID, pathSaved)
	require.NoError(t, err)

	// check .Read(), .List(), .Stat() --------------------------------------

	dataReaded, err = filesOp.Read(bucketID, pathSaved)
	require.Error(t, err)
	require.Nil(t, dataReaded)

	fis, err = filesOp.List(bucketID, filepath.Dir(pathSaved), 0)
	require.NoError(t, err)

	found = false
	for _, fi := range fis {
		if fi.Name == filepath.Base(pathSaved) {
			found = true
			require.FailNowf(t, "this file should be removed", "%#v", fi)
		}
	}

	fi, err = filesOp.Stat(bucketID, filepath.Dir(pathSaved))
	require.NoError(t, err)
	require.NotNil(t, fi)
	require.True(t, fi.IsDir)
	require.Equalf(t, size0, fi.Size, "%#v", fi)

	return pathSaved
}
