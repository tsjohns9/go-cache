package cache

import "github.com/tsjohns9/go-cache/pkg/sets"

type IndexFunc func(obj interface{}) ([]string, error)
type Indexers map[string]IndexFunc
type Index map[string]sets.String
type Indices map[string]Index
type KeyFunc func(obj interface{}) (string, error)

type Indexer interface {
	Store
	AddIndexers(newIndexers Indexers) error
	ByIndex(indexName, indexedValue string) ([]interface{}, error)
	GetIndexers() Indexers
	Index(indexName string, obj interface{}) ([]interface{}, error)
	IndexKeys(indexName, indexedValue string) ([]string, error)
	ListIndexFuncValues(indexName string) []string
}

func NewIndexer(keyFunc KeyFunc, indexers Indexers) Indexer {
	return &store{
		keyFunc:      keyFunc,
		cacheStorage: NewThreadSafeStore(indexers, Indices{}),
	}
}
