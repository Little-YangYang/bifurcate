package bifurcate

type Context struct {
	Data map[string]string

	Err error
}

func (c *Context) SetData(key string, value string) {
	c.Data[key] = value
}

func (c *Context) GetData(key string) string {
	return c.Data[key]
}

func NewContext() *Context {
	return &Context{
		Data: make(map[string]string),
	}
}
