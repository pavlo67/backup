package files

import (
	"os"
	"time"
)

type BucketID string

type FileInfo struct {
	Path string
	// Name      string
	IsDir     bool
	Size      int64
	CreatedAt time.Time
}

type FilesInfo []FileInfo

func (fis FilesInfo) Append(basePath string, info os.FileInfo) (FilesInfo, error) {
	path := info.Name()

	//if len(path) <= len(basePath) {
	//	return nil, errors.Errorf("wrong path (%s) on basePath = '%s'", path, basePath)
	//}

	if info.IsDir() {
		fis = append(fis, FileInfo{
			Path: path,
			// Path:      path[len(basePath):],
			IsDir:     true,
			CreatedAt: info.ModTime(),
		})
	} else {
		fis = append(fis, FileInfo{
			Path: path,
			// Path:      path[len(basePath):],
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return fis, nil
}

type Operator interface {
	Save(bucketID BucketID, path, newFilePattern string, data []byte) (string, error)
	Read(bucketID BucketID, path string) ([]byte, error)
	Remove(bucketID BucketID, path string) error
	List(bucketID BucketID, path string, depth int) (FilesInfo, error)
	Stat(bucketID BucketID, path string, depth int) (*FileInfo, error)
}
