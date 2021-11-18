package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLruCache(t *testing.T) {
	c := NewLruCache(3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)
	c.Put(4, 4)
	assert.Equal(t, 3, c.Len())
	assert.Nil(t, c.Get(1))
	assert.Equal(t, 2, c.Get(2))
	assert.Equal(t, 3, c.Get(3))
	assert.Equal(t, 4, c.Get(4))
}
