package lru

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestPutAndGetSuccess(t *testing.T) {

    cache := New(10)
    cache.Put("key1", "value1")

    assert.Equal(t, "value1", cache.Get("key1").(string))

}

func TestPutAndGetSuccessSizeLimited(t *testing.T) {

    cache := New(2)
    cache.Put("key1", "value1")
    cache.Put("key2", "value2")
    cache.Put("key3", "value3")

    assert.Nil(t, cache.Get("key1"))
    assert.Equal(t, "value2", cache.Get("key2").(string))
    assert.Equal(t, "value3", cache.Get("key3").(string))

}

func TestPutAndGeRemovingLeastUsed(t *testing.T) {

    cache := New(3)

    cache.Put("key1", "value1")
    cache.Put("key2", "value2")
    cache.Put("key3", "value3")

    cache.Get("key2")
    cache.Get("key1")

    cache.Put("key4", "value4")

    assert.Equal(t, "value1", cache.Get("key1"))
    assert.Equal(t, "value2", cache.Get("key2").(string))
    assert.Nil(t, cache.Get("key3"))
    assert.Equal(t, "value4", cache.Get("key4").(string))

}

func TestPurge(t *testing.T) {

    cache := New(3)

    cache.Put("key1", "value1")
    cache.Put("key2", "value2")

    assert.Equal(t, "value1", cache.Get("key1"))
    assert.Equal(t, "value2", cache.Get("key2"))

    cache.Clear()

    assert.False(t, cache.ContainsKey("key1"))
    assert.False(t, cache.ContainsKey("key2"))

}

func TestRemove(t *testing.T) {

    cache := New(3)

    cache.Put("key1", "value1")
    cache.Put("key2", "value2")

    assert.Equal(t, "value1", cache.Get("key1"))
    assert.Equal(t, "value2", cache.Get("key2"))

    cache.Remove("key2")

    assert.True(t, cache.ContainsKey("key1"))
    assert.False(t, cache.ContainsKey("key2"))

}

func TestKeySet(t *testing.T) {

    cache := New(3)

    cache.Put("key1", "value1")
    cache.Put("key2", "value2")

    assert.Equal(t, []interface{}{"key1", "key2"}, cache.KeySet())

}

func TestSize(t *testing.T) {

    cache := New(3)

    cache.Put("key1", "value1")
    cache.Put("key2", "value2")

    assert.Equal(t, 2, cache.Size())

}
