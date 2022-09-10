package cache

import (
	"fmt"
	"testing"

	v1 "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/api/v1"
	management "github.com/Axway/agent-sdk/pkg/apic/apiserver/models/management/v1alpha1"
	"github.com/Axway/agent-sdk/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestGenericLister(t *testing.T) {
	indexer := NewAPIServiceIndexer()
	svc1 := management.NewAPIService("svc1", "one")
	svc1.Attributes = map[string]string{
		"key":  "val",
		"key2": "val3",
		"key4": "val4",
	}
	svc1.Tags = []string{"abc", "123", "456"}
	util.SetAgentDetails(svc1, map[string]interface{}{
		"externalAPIID": "123",
		// "externalAPIName": "dumb-one",
	})

	svc2 := management.NewAPIService("svc2", "one")
	svc2.Attributes = map[string]string{
		"abc": "123",
	}
	util.SetAgentDetails(svc1, map[string]interface{}{
		"externalAPIID":   "321",
		"externalAPIName": "dumb-two",
	})

	services := []v1.Interface{
		svc1,
		svc2,
		management.NewAPIService("svc3", "one"),
		management.NewAPIService("svc1", "two"),
	}

	for _, s := range services {
		err := indexer.Add(s)
		assert.Nil(t, err)
	}

	lister := NewGenericLister(indexer, management.NewAPIService("a", "a").GroupVersionKind)
	sel := NewSelector()

	r1 := NewAttrRule("key", DoubleEquals, []string{"val", "val3"})
	r2 := NewAttrRule("key2", DoubleEquals, []string{"val3"})
	r3 := NewTagRule("abc")

	sel.Add(r1)
	sel.Add(r2)
	sel.Add(r3)

	attr := AttrMap(svc1.Attributes)
	ok := sel.Matches(attr)
	assert.True(t, ok)

	ret, err := lister.ByScope("one").List(sel)
	assert.Nil(t, err)
	fmt.Println("ret", ret)
}
