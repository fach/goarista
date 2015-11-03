// Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package key_test

import (
	"fmt"
	"testing"

	. "github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/test"
)

func TestKeyEqual(t *testing.T) {
	tests := []struct {
		a      Key
		b      Key
		result bool
	}{{
		a:      New("foo"),
		b:      New("foo"),
		result: true,
	}, {
		a:      New("foo"),
		b:      New("bar"),
		result: false,
	}, {
		a:      New(map[string]interface{}{}),
		b:      New("bar"),
		result: false,
	}, {
		a:      New(map[string]interface{}{}),
		b:      New(map[string]interface{}{}),
		result: true,
	}, {
		a:      New(map[string]interface{}{"a": 3}),
		b:      New(map[string]interface{}{}),
		result: false,
	}, {
		a:      New(map[string]interface{}{"a": 3}),
		b:      New(map[string]interface{}{"b": 4}),
		result: false,
	}, {
		a:      New(map[string]interface{}{"a": 3}),
		b:      New(map[string]interface{}{"a": 4}),
		result: false,
	}, {
		a:      New(map[string]interface{}{"a": 3}),
		b:      New(map[string]interface{}{"a": 3}),
		result: true,
	}}

	for _, tcase := range tests {
		if tcase.a.Equal(tcase.b) != tcase.result {
			t.Errorf("Wrong result for case:\na: %#v\nb: %#v\nresult: %#v",
				tcase.a,
				tcase.b,
				tcase.result)
		}
	}

	if New("a").Equal(32) == true {
		t.Error("Wrong result for different types case")
	}
}

func TestIsHashable(t *testing.T) {
	tests := []struct {
		k interface{}
		h bool
	}{{
		true,
		true,
	}, {
		false,
		true,
	}, {
		uint8(3),
		true,
	}, {
		uint16(3),
		true,
	}, {
		uint32(3),
		true,
	}, {
		uint64(3),
		true,
	}, {
		int8(3),
		true,
	}, {
		int16(3),
		true,
	}, {
		int32(3),
		true,
	}, {
		int64(3),
		true,
	}, {
		float32(3.2),
		true,
	}, {
		float64(3.3),
		true,
	}, {
		"foobar",
		true,
	}, {
		map[string]interface{}{"foo": "bar"},
		false,
	}}

	for _, tcase := range tests {
		if New(tcase.k).IsHashable() != tcase.h {
			t.Errorf("Wrong result for case:\nk: %#v",
				tcase.k)

		}
	}
}

func TestGetFromMap(t *testing.T) {
	tests := []struct {
		k     Key
		m     map[Key]interface{}
		v     interface{}
		found bool
	}{{
		k:     New("a"),
		m:     map[Key]interface{}{New("a"): "b"},
		v:     "b",
		found: true,
	}, {
		k:     New(uint32(35)),
		m:     map[Key]interface{}{New(uint32(35)): "c"},
		v:     "c",
		found: true,
	}, {
		k:     New(uint32(37)),
		m:     map[Key]interface{}{New(uint32(36)): "c"},
		found: false,
	}, {
		k:     New(uint32(37)),
		m:     map[Key]interface{}{},
		found: false,
	}, {
		k: New(map[string]interface{}{"a": "b", "c": uint64(4)}),
		m: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
		},
		v:     "foo",
		found: true,
	}}

	for _, tcase := range tests {
		v, ok := tcase.k.GetFromMap(tcase.m)
		if tcase.found != ok {
			t.Errorf("Wrong retrieval result for case:\nk: %#v\nm: %#v\nv: %#v",
				tcase.k,
				tcase.m,
				tcase.v)
		} else if tcase.found && !ok {
			t.Errorf("Unable to retrieve value for case:\nk: %#v\nm: %#v\nv: %#v",
				tcase.k,
				tcase.m,
				tcase.v)
		} else if tcase.found && !test.DeepEqual(tcase.v, v) {
			t.Errorf("Wrong result for case:\nk: %#v\nm: %#v\nv: %#v",
				tcase.k,
				tcase.m,
				tcase.v)
		}
	}
}

func TestDeleteFromMap(t *testing.T) {
	tests := []struct {
		k Key
		m map[Key]interface{}
		r map[Key]interface{}
	}{{
		k: New("a"),
		m: map[Key]interface{}{New("a"): "b"},
		r: map[Key]interface{}{},
	}, {
		k: New("b"),
		m: map[Key]interface{}{New("a"): "b"},
		r: map[Key]interface{}{New("a"): "b"},
	}, {
		k: New("a"),
		m: map[Key]interface{}{},
		r: map[Key]interface{}{},
	}, {
		k: New(uint32(35)),
		m: map[Key]interface{}{New(uint32(35)): "c"},
		r: map[Key]interface{}{},
	}, {
		k: New(uint32(36)),
		m: map[Key]interface{}{New(uint32(35)): "c"},
		r: map[Key]interface{}{New(uint32(35)): "c"},
	}, {
		k: New(uint32(37)),
		m: map[Key]interface{}{New(uint32(36)): "c"},
		r: map[Key]interface{}{New(uint32(36)): "c"},
	}, {
		k: New(uint32(37)),
		m: map[Key]interface{}{},
		r: map[Key]interface{}{},
	}, {
		k: New(map[string]interface{}{"a": "b", "c": uint64(4)}),
		m: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
		},
		r: map[Key]interface{}{},
	}}

	for _, tcase := range tests {
		tcase.k.DeleteFromMap(tcase.m)
		if !test.DeepEqual(tcase.m, tcase.r) {
			t.Errorf("Wrong result for case:\nk: %#v\nm: %#v\nr: %#v",
				tcase.k,
				tcase.m,
				tcase.r)
		}
	}
}

func TestSetToMap(t *testing.T) {
	tests := []struct {
		k Key
		v interface{}
		m map[Key]interface{}
		r map[Key]interface{}
	}{{
		k: New("a"),
		v: "c",
		m: map[Key]interface{}{New("a"): "b"},
		r: map[Key]interface{}{New("a"): "c"},
	}, {
		k: New("b"),
		v: uint64(56),
		m: map[Key]interface{}{New("a"): "b"},
		r: map[Key]interface{}{
			New("a"): "b",
			New("b"): uint64(56),
		},
	}, {
		k: New("a"),
		v: "foo",
		m: map[Key]interface{}{},
		r: map[Key]interface{}{New("a"): "foo"},
	}, {
		k: New(uint32(35)),
		v: "d",
		m: map[Key]interface{}{New(uint32(35)): "c"},
		r: map[Key]interface{}{New(uint32(35)): "d"},
	}, {
		k: New(uint32(36)),
		v: true,
		m: map[Key]interface{}{New(uint32(35)): "c"},
		r: map[Key]interface{}{
			New(uint32(35)): "c",
			New(uint32(36)): true,
		},
	}, {
		k: New(uint32(37)),
		v: false,
		m: map[Key]interface{}{New(uint32(36)): "c"},
		r: map[Key]interface{}{
			New(uint32(36)): "c",
			New(uint32(37)): false,
		},
	}, {
		k: New(uint32(37)),
		v: "foobar",
		m: map[Key]interface{}{},
		r: map[Key]interface{}{New(uint32(37)): "foobar"},
	}, {
		k: New(map[string]interface{}{"a": "b", "c": uint64(4)}),
		v: "foobar",
		m: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
		},
		r: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foobar",
		},
	}, {
		k: New(map[string]interface{}{"a": "b", "c": uint64(7)}),
		v: "foobar",
		m: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
		},
		r: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
			New(map[string]interface{}{"a": "b", "c": uint64(7)}): "foobar",
		},
	}, {
		k: New(map[string]interface{}{"a": "b", "d": uint64(6)}),
		v: "barfoo",
		m: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
		},
		r: map[Key]interface{}{
			New(map[string]interface{}{"a": "b", "c": uint64(4)}): "foo",
			New(map[string]interface{}{"a": "b", "d": uint64(6)}): "barfoo",
		},
	}}

	for i, tcase := range tests {
		tcase.k.SetToMap(tcase.m, tcase.v)
		if !test.DeepEqual(tcase.m, tcase.r) {
			t.Errorf("Wrong result for case %d:\nk: %#v\nm: %#v\nr: %#v",
				i,
				tcase.k,
				tcase.m,
				tcase.r)
		}
	}
}

func BenchmarkSetToMapWithStringKey(b *testing.B) {
	m := map[Key]interface{}{
		New("a"):   true,
		New("a1"):  true,
		New("a2"):  true,
		New("a3"):  true,
		New("a4"):  true,
		New("a5"):  true,
		New("a6"):  true,
		New("a7"):  true,
		New("a8"):  true,
		New("a9"):  true,
		New("a10"): true,
		New("a11"): true,
		New("a12"): true,
		New("a13"): true,
		New("a14"): true,
		New("a15"): true,
		New("a16"): true,
		New("a17"): true,
		New("a18"): true,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(fmt.Sprintf("b%d", i)).SetToMap(m, true)
	}
}

func BenchmarkSetToMapWithUint64Key(b *testing.B) {
	m := map[Key]interface{}{
		New(uint64(1)):  true,
		New(uint64(2)):  true,
		New(uint64(3)):  true,
		New(uint64(4)):  true,
		New(uint64(5)):  true,
		New(uint64(6)):  true,
		New(uint64(7)):  true,
		New(uint64(8)):  true,
		New(uint64(9)):  true,
		New(uint64(10)): true,
		New(uint64(11)): true,
		New(uint64(12)): true,
		New(uint64(13)): true,
		New(uint64(14)): true,
		New(uint64(15)): true,
		New(uint64(16)): true,
		New(uint64(17)): true,
		New(uint64(18)): true,
		New(uint64(19)): true,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(uint64(i)).SetToMap(m, true)
	}
}

func BenchmarkGetFromMapWithMapKey(b *testing.B) {
	m := map[Key]interface{}{
		New(map[string]interface{}{"a0": true}):  true,
		New(map[string]interface{}{"a1": true}):  true,
		New(map[string]interface{}{"a2": true}):  true,
		New(map[string]interface{}{"a3": true}):  true,
		New(map[string]interface{}{"a4": true}):  true,
		New(map[string]interface{}{"a5": true}):  true,
		New(map[string]interface{}{"a6": true}):  true,
		New(map[string]interface{}{"a7": true}):  true,
		New(map[string]interface{}{"a8": true}):  true,
		New(map[string]interface{}{"a9": true}):  true,
		New(map[string]interface{}{"a10": true}): true,
		New(map[string]interface{}{"a11": true}): true,
		New(map[string]interface{}{"a12": true}): true,
		New(map[string]interface{}{"a13": true}): true,
		New(map[string]interface{}{"a14": true}): true,
		New(map[string]interface{}{"a15": true}): true,
		New(map[string]interface{}{"a16": true}): true,
		New(map[string]interface{}{"a17": true}): true,
		New(map[string]interface{}{"a18": true}): true,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := New(map[string]interface{}{fmt.Sprintf("a%d", i%19): true})
		_, found := key.GetFromMap(m)
		if !found {
			b.Fatalf("WTF: %#v", key)
		}
	}
}