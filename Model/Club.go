package Model

type Club struct {
	name   string
	id     int
}

func (c *Club) GetName() string {
	return c.name
}

func (c *Club) GetID() int {
	return c.id
}