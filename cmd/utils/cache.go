package utils

type Cache struct {
	Data map[string]any
}

func NewCache() Cache {
	data := make(map[string]any)
	return Cache{data}
}

func (c *Cache) PutValue(field string, value any) {
	c.Data[field] = value
}

func (c *Cache) GetValue(field string) any {
	return c.Data[field]
}
