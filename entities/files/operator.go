package files

import (
	"os"
	"time"
)

type Operator interface {
	Save(bucketID BucketID, path, newFilePattern string, data []byte) (string, error)
	Read(bucketID BucketID, path string) ([]byte, error)
	Remove(bucketID BucketID, path string) error
	List(bucketID BucketID, path string, depth int) (Items, error)
	Stat(bucketID BucketID, path string, depth int) (*Item, error)
}

type BucketID string

type Buckets map[BucketID]string

type Item struct {
	Path string
	// Name      string
	IsDir     bool
	Size      int64
	CreatedAt time.Time
}

type Items []Item

func (fis Items) Append(basePath string, info os.FileInfo) (Items, error) {
	path := info.Name()

	//if len(path) <= len(basePath) {
	//	return nil, fmt.Errorf("wrong path (%s) on basePath = '%s'", path, basePath)
	//}

	if info.IsDir() {
		fis = append(fis, Item{
			Path: path,
			// Path:      path[len(basePath):],
			IsDir:     true,
			CreatedAt: info.ModTime(),
		})
	} else {
		fis = append(fis, Item{
			Path: path,
			// Path:      path[len(basePath):],
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
		})
	}

	return fis, nil
}
