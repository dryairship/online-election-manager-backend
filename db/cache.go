package db

type Cache struct {
	IsSet         bool
	Data          interface{}
	ResetFunction func() (interface{}, error)
}

func (cache *Cache) Reset() error {
	var err error
	cache.Data, err = cache.ResetFunction()
	return err
}
