package cache

import (
	"fmt"

	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	"github.com/Axway/agent-sdk/pkg/util"
)

// GenericLister is a lister on top of an Indexer
type GenericLister interface {
	// List will return all objects across all scopes
	List(selector Selector) (ret []v1.Interface, err error)
	// Get will attempt to retrieve assuming that name==key
	Get(name string) (v1.Interface, error)
	// ByScope will give you a GenericScopeLister for one resource scope
	ByScope(scope string) GenericScopeLister
}

// GenericScopeLister is a lister skin on a generic Indexer
type GenericScopeLister interface {
	// List will return all objects in this scope
	List(selector Selector) (ret []v1.Interface, err error)
	// Get will attempt to retrieve by scope and name
	Get(name string) (v1.Interface, error)
}

type genericLister struct {
	indexer  Indexer
	resource v1.GroupVersionKind
}

// NewGenericLister creates a new instance for the genericLister.
func NewGenericLister(indexer Indexer, resource v1.GroupVersionKind) GenericLister {
	return &genericLister{indexer: indexer, resource: resource}
}

func (g *genericLister) List(selector Selector) ([]v1.Interface, error) {
	return ListAll(g.indexer, selector)
}

func (g *genericLister) Get(name string) (v1.Interface, error) {
	obj, exists, err := g.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("item %s of kind %s was not found in the store", name, g.resource.Kind)
	}
	return obj.(v1.Interface), nil
}

func (g *genericLister) ByScope(scope string) GenericScopeLister {
	return &genericScopeLister{
		indexer:  g.indexer,
		scope:    scope,
		resource: v1.GroupVersionKind{},
	}
}

type genericScopeLister struct {
	indexer  Indexer
	scope    string
	resource v1.GroupVersionKind
}

func (s *genericScopeLister) List(selector Selector) (ret []v1.Interface, err error) {
	return ListAllByScope(s.indexer, s.scope, selector)
}

func (s *genericScopeLister) Get(name string) (v1.Interface, error) {
	obj, exists, err := s.indexer.GetByKey(s.scope + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("item %s of kind %s was not found in the store", name, s.resource.Kind)
	}
	return obj.(v1.Interface), nil
}

// ListAll calls appendFn with each value retrieved from store which matches the selector.
func ListAll(store Store, selector Selector) ([]v1.Interface, error) {
	if selector == nil {
		selector = NewSelector()
	}

	selectAll := selector.Empty()
	var list []v1.Interface

	for _, m := range store.List() {
		if selectAll {
			list = append(list, m.(v1.Interface))
			continue
		}
		metadata, err := AccessorMeta(m)
		if err != nil {
			return list, err
		}
		if selector.Matches(AttrMap(metadata.GetAttributes())) {
			list = append(list, m.(v1.Interface))
		}
	}

	return list, nil
}

func ListAllByScope(indexer Indexer, scope string, selector Selector) ([]v1.Interface, error) {
	selectAll := selector.Empty()
	var list []v1.Interface

	if scope == "" {
		return ListAll(indexer, selector)
	}

	rm := &v1.ResourceMeta{Metadata: v1.Metadata{Scope: v1.MetadataScope{Name: scope}}}
	items, err := indexer.Index(ScopeIndex, rm)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve list of objects using index %s", ScopeIndex)
	}

	for _, m := range items {
		if selectAll {
			list = append(list, m.(v1.Interface))
			continue
		}

		riMeta, err := AccessorMeta(m)
		if err != nil {
			return list, err
		}

		if matchXAgentDetails(selector, riMeta) {
			list = append(list, m.(v1.Interface))
			continue
		}

		if matchAttributes(selector, riMeta) {
			list = append(list, m.(v1.Interface))
			continue
		}

		if matchTags(selector, riMeta) {
			list = append(list, m.(v1.Interface))
			continue
		}
	}

	return list, nil
}

func matchAttributes(sel Selector, metadata v1.Meta) bool {
	return sel.Matches(AttrMap(metadata.GetAttributes()))
}

func matchTags(sel Selector, metadata v1.Meta) bool {
	return sel.Matches(Tags(metadata.GetTags()))
}

func matchXAgentDetails(sel Selector, metadata v1.Meta) bool {
	d := util.GetAgentDetails(metadata)
	if d == nil {
		return false
	}

	details := util.MapStringInterfaceToStringString(d)
	return sel.Matches(AttrMap(details))
}
