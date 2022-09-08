package cache

type Store interface {
	Add(obj interface{}) error
	Delete(obj interface{}) error
	Get(obj interface{}) (item interface{}, exists bool, err error)
	GetByKey(key string) (item interface{}, exists bool, err error)
	List() []interface{}
	ListKeys() []string
	Update(obj interface{}) error
}

func NewStore(keyFunc KeyFunc) Store {
	return &store{
		keyFunc:      keyFunc,
		cacheStorage: NewThreadSafeStore(Indexers{}, Indices{}),
	}
}

type store struct {
	keyFunc      KeyFunc
	cacheStorage ThreadSafeStore
}

func (c *store) Add(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return err
	}
	c.cacheStorage.Add(key, obj)
	return nil
}

func (c *store) Update(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return err
	}
	c.cacheStorage.Update(key, obj)
	return nil
}

func (c *store) Delete(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return err
	}
	c.cacheStorage.Delete(key)
	return nil
}

func (c *store) List() []interface{} {
	return c.cacheStorage.List()
}

func (c *store) ListKeys() []string {
	return c.cacheStorage.ListKeys()
}

func (c *store) Get(obj interface{}) (item interface{}, exists bool, err error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, false, err
	}
	return c.GetByKey(key)
}

func (c *store) GetByKey(key string) (item interface{}, exists bool, err error) {
	item, exists = c.cacheStorage.Get(key)
	return item, exists, nil
}

func (c *store) Index(indexName string, obj interface{}) ([]interface{}, error) {
	return c.cacheStorage.Index(indexName, obj)
}

func (c *store) IndexKeys(indexName, indexKey string) ([]string, error) {
	return c.cacheStorage.IndexKeys(indexName, indexKey)
}

func (c *store) ListIndexFuncValues(indexName string) []string {
	return c.cacheStorage.ListIndexFuncValues(indexName)
}

func (c *store) ByIndex(indexName, indexKey string) ([]interface{}, error) {
	return c.cacheStorage.ByIndex(indexName, indexKey)
}

func (c *store) GetIndexers() Indexers {
	return c.cacheStorage.GetIndexers()
}

func (c *store) AddIndexers(newIndexers Indexers) error {
	return c.cacheStorage.AddIndexers(newIndexers)
}
