package cache

import (
	"time"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	c := New()
	for i := 0; i < 1000; i++ {
		item := rand.Int()
		key := strconv.Itoa(item)
		c.Add(key, item)
	}

	if length, expected := len(c.Keys()), 1000; length != expected {
		t.Errorf("Cache had %d keys, expected %d", length, expected)
	}
}

func TestAddBasic(t *testing.T) {
	c := New()
	c.Add("1", 1)

	if result, expected := c.Get("1"), 1; !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}
}

func TestClear(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Add(strconv.Itoa(i), i)
	}

	c.Clear()
	if length, expected := len(c.Keys()), 0; length != expected {
		t.Errorf("Cache had %d keys, expected %d", length, expected)
	}
}

func TestDelete(t *testing.T) {
	c := New()
	for i := 0; i < 1000; i++ {
		item := rand.Int()
		key := strconv.Itoa(item)
		c.Add(key, item)
	}

	for _, key := range c.Keys() {
		c.Delete(key)
	}

	if length, expected := len(c.Keys()), 0; length != expected {
		t.Errorf("Cache had %d keys, expected %d", length, expected)
	}
}

func TestDeleteBasic(t *testing.T) {
	c := New()
	c.Add("1", 1)
	c.Delete("1")

	if result := c.Get("1"); result != nil {
		t.Errorf("Result was %#v, expected nil", result)
	}
}

func TestClearEvery(t *testing.T) {
	c := New()
	ticker := c.ClearEvery(time.Nanosecond)

	for i := 0; i < 10; i++ {
		c.Add(strconv.Itoa(i), i)
	}

	// wait for first tick - delayed as the goroutine starts up
	<-ticker.C

	if length, expected := len(c.Keys()), 0; length != expected {
		t.Errorf("Cache had %d keys, expected %d", length, expected)
	}
}

func TestGet(t *testing.T) {
	c := New()
	c.Add("1", 1)

	if result, expected := c.Get("1"), 1; !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}

	if result := c.Get("2"); result != nil {
		t.Errorf("Result was %#v, expected nil", result)
	}
}

func TestGetf(t *testing.T) {
        c := New()
        c.Add("1", 1)

	result, exists := c.Getf("1")
	if !exists {
		t.Error("Exists was false, expected true")
	}

	if expected := 1; !reflect.DeepEqual(result, expected) {
                t.Errorf("Result was %#v, expected %#v", result, expected)
        }
}


func TestItems(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Add(strconv.Itoa(i), i)
	}

	expected := map[string]interface{}{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	}

	if result := c.Items(); !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}
}

func TestKeys(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Add(strconv.Itoa(i), i)
	}

	expected := []string{"0", "1", "2", "3", "4"}
	if result := c.Keys(); !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}
}
