package files_fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pkg/errors"

	"github.com/pavlo67/tools/components/files"
)

var _ files.Operator = &filesFS{}

type filesFS struct {
	buckets files.Buckets
}

const onNew = "on filesFS.New(): "

func New(buckets files.Buckets) (files.Operator, crud.Cleaner, error) {
	if len(buckets) < 1 {
		return nil, nil, errors.New(onNew + ": no buckets to process")
	}

	var err error
	for bucketID, basePath := range buckets {
		if buckets[bucketID], err = filelib.Dir(basePath); err != nil {
			return nil, nil, errors.Wrapf(err, onNew+": creating bucket '%s'", bucketID)
		}
	}

	filesOp := filesFS{
		buckets: buckets,
	}

	return &filesOp, &filesOp, nil
}

const onSave = "on filesFS.Save()"

func (filesOp *filesFS) Save(bucketID files.BucketID, path, newFilePattern string, data []byte, options *crud.Options) (string, error) {
	basePath := filesOp.buckets[bucketID]
	if basePath == "" {
		return "", errors.Errorf(onSave+": wrong bucket (%s)", bucketID)
	}

	path = basePath + path

	var err error
	var dirPath string
	var file *os.File

	// TODO!!! check if dirPath doesn't contain "/../"
	if newFilePattern != "" {
		if dirPath, err = filelib.Dir(path); err != nil {
			return "", errors.Wrapf(err, onSave+": wrong path (%s)", path)
		}
		if file, err = ioutil.TempFile(dirPath, newFilePattern); err != nil {
			return "", errors.Wrapf(err, onSave+": can't ioutil.TempFile(%s, %s)", dirPath, newFilePattern)
		}
	} else {
		var filename string
		dirPath, filename = filepath.Dir(path), filepath.Base(path)

		dirPath, err = filelib.Dir(dirPath)

		if file, err = os.OpenFile(dirPath+filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
			return "", errors.Wrapf(err, onSave+": can't os.OpenFile(%s, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644)", dirPath+filename)
		}
	}
	defer func() {
		if err := file.Close(); err != nil {
			l.Errorf(onSave+": on file.Close() got %s", err)
		}

	}()

	filename := strings.ReplaceAll(file.Name(), "/./", "/")

	if len(filename) <= len(basePath) {
		return "", errors.Errorf(onSave+": wrong filename (%s) on basePath = '%s'", filename, basePath)
	}

	if _, err = file.Write(data); err != nil {
		return "", errors.Wrapf(err, onSave+": can't file.Write(%s)", file.Name())
	}

	return filename[len(basePath):], nil
}

const onRead = "on filesFS.Read()"

func (filesOp *filesFS) Read(bucketID files.BucketID, path string, options *crud.Options) ([]byte, error) {
	basePath := filesOp.buckets[bucketID]
	if basePath == "" {
		return nil, errors.Errorf(onRead+": wrong bucket (%s)", bucketID)
	}
	filePath := basePath + path

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, onRead+": can't ioutil.ReadFile(%s)", filePath)
	}

	return data, nil
}

const onRemove = "on filesFS.Remove()"

func (filesOp *filesFS) Remove(bucketID files.BucketID, path string, options *crud.Options) error {
	basePath := filesOp.buckets[bucketID]
	if basePath == "" {
		return errors.Errorf(onRemove+": wrong bucket (%s)", bucketID)
	}
	filePath := basePath + path

	if err := os.Remove(filePath); err != nil {
		return errors.Wrapf(err, onRemove+": can't os.Remove(%s)", filePath)
	}

	return nil
}

const onList = "on filesFS.List()"

func (filesOp *filesFS) List(bucketID files.BucketID, path string, depth int, options *crud.Options) (files.FilesInfo, error) {
	basePath := filesOp.buckets[bucketID]
	if basePath == "" {
		return nil, errors.Errorf(onRead+": wrong bucket (%s)", bucketID)
	}
	filePath := basePath + path

	var filesInfo files.FilesInfo

	if depth == 0 {
		fis, err := ioutil.ReadDir(filePath)

		if err != nil {
			return nil, errors.Wrapf(err, onList+": can't ioutil.ReadDir(%s)", filePath)
		}

		for _, fi := range fis {
			filesInfo, err = filesInfo.Append("", fi) // basePath
			if err != nil {
				return nil, errors.Wrap(err, onList)
			}
		}

		return filesInfo, nil
	}

	// TODO: process depth > 0 more thoroughly here
	err := filepath.Walk(filePath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filesInfo, err = filesInfo.Append("", fi) // basePath
		if err != nil {
			return errors.Wrap(err, onList)
		}

		return nil
	})

	return filesInfo, err
}

const onStat = "on filesFS.Stat()"

func (filesOp *filesFS) Stat(bucketID files.BucketID, path string, depth int, options *crud.Options) (*files.FileInfo, error) {
	basePath := filesOp.buckets[bucketID]
	if basePath == "" {
		return nil, errors.Errorf(onStat+": wrong bucket (%s)", bucketID)
	}
	filePath := basePath + path

	fi, err := os.Stat(filePath)
	if err != nil {
		//if os.IsNotExist(err) {
		//	return nil, nil
		//}
		return nil, errors.Wrapf(err, onStat+": can't  os.Stat(%s)", filePath)
	}

	filesInfo, err := files.FilesInfo{}.Append("", fi) // basePath
	if err != nil || len(filesInfo) != 1 {
		return nil, errors.Errorf(onStat+": got %#v / %s", filesInfo, err)
	}

	fileInfo := filesInfo[0]

	if depth != 0 && fileInfo.IsDir {
		// TODO: process depth > 0 more thoroughly here
		err = filepath.Walk(filePath, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !fi.IsDir() {
				fileInfo.Size += fi.Size()
			}

			return nil
		})
	}

	return &fileInfo, err

}
