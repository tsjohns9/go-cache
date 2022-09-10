package cache

import (
	"fmt"
	"strings"

	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/tsjohns9/go-cache/pkg/sets"
)

type IndexFunc func(obj interface{}) ([]string, error)
type Indexers map[string]IndexFunc
type Index map[string]sets.String
type Indices map[string]Index
type KeyFunc func(obj interface{}) (string, error)

const ScopeIndex = "scope"

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

func NewAPIServiceIndexer() Indexer {
	return NewIndexer(ScopeAndNameKeyFunc, map[string]IndexFunc{
		ScopeIndex: ScopeIndexFunc,
	})
}

func ScopeAndNameKeyFunc(obj interface{}) (string, error) {
	meta, err := AccessorMeta(obj)
	if err != nil {
		return "", fmt.Errorf("cannot create key for %T", obj)
	}

	if meta == nil {
		return "", fmt.Errorf("cannot create key for nil value of %T", obj)
	}

	scopeName := meta.GetMetadata().Scope.Name
	if len(scopeName) > 0 {
		return scopeName + "/" + meta.GetName(), nil
	}

	return meta.GetName(), nil
}

func IDKeyFunc(obj interface{}) (string, error) {
	meta, err := AccessorMeta(obj)
	if err != nil {
		return "", fmt.Errorf("cannot create key for %T", obj)
	}

	if meta == nil {
		return "", fmt.Errorf("cannot create key for nil value of %T", obj)
	}

	return meta.GetMetadata().ID, nil
}

func ScopeIndexFunc(obj interface{}) ([]string, error) {
	ri, err := AccessorMeta(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}

	if ri.GetMetadata().Scope.Name == "" {
		return []string{}, nil
	}

	return []string{ri.GetMetadata().Scope.Name}, nil
}

func AccessorInterface(obj interface{}) (v1.Interface, error) {
	switch t := obj.(type) {
	case v1.Interface:
		return t, nil
	default:
		return nil, fmt.Errorf("expected %T to be of type v1.Interface", obj)
	}
}

func AccessorMeta(obj interface{}) (v1.Meta, error) {
	switch t := obj.(type) {
	case v1.Meta:
		return t, nil
	default:
		return nil, fmt.Errorf("expected %T to be of type v1.Meta", obj)
	}
}

func AccessorXAgentDetails(obj interface{}) (map[string]interface{}, error) {
	meta, err := AccessorMeta(obj)
	if meta == nil {
		return nil, fmt.Errorf("expected %T to be of type v1.Meta: %s", obj, err)
	}

	return util.GetAgentDetails(meta), nil
}

func SplitScopeAndNameKey(key string) (scope, name string, err error) {
	parts := strings.Split(key, "/")

	switch len(parts) {
	case 1:
		// name only, no namespace
		return "", parts[0], nil
	case 2:
		// namespace and name
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected key format: %q", key)
}
