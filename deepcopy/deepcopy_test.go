package deepcopy

import (
	"reflect"
	"testing"
)

type st struct {
	Int          int
	PInt         *int
	String       string
	PString      *string
	ArrayInt     [2]int
	SliceInt     []int
	MapInt       map[int]int
	MapPInt      map[int]*int
	Interface    interface{}
	PInterface   *interface{}
	Struct       stt
	PStruct      *stt
	ArrayStruct  [2]st1
	SliceStruct  []st1
	SlicePStruct []*st1
}

type stt struct {
	Int        int
	PInt       *int
	String     string
	PString    *string
	SliceInt   []int
	MapInt     map[int]int
	Interface  interface{}
	PInterface *interface{}
}

type st1 struct {
	Value int
}

func TestDeepCopy(t *testing.T) {
	{
		s1 := st{
			Int:        10,
			PInt:       new(int),
			String:     "string",
			PString:    new(string),
			ArrayInt:   [2]int{1, 2},
			SliceInt:   []int{1, 2},
			MapInt:     map[int]int{1: 1, 2: 2},
			MapPInt:    map[int]*int{1: new(int), 2: new(int)},
			Interface:  "interface",
			PInterface: new(interface{}),
			Struct: stt{
				Int:        20,
				PInt:       nil,
				String:     "string2",
				PString:    nil,
				SliceInt:   []int{3, 4},
				MapInt:     map[int]int{3: 3, 4: 4},
				Interface:  "interface2",
				PInterface: nil,
			},
			PStruct:      new(stt),
			ArrayStruct:  [2]st1{{1}, {2}},
			SliceStruct:  []st1{{1}, {2}},
			SlicePStruct: []*st1{{1}, {2}},
		}
		*s1.PInt = 30
		*s1.PString = "string"
		*s1.PInterface = "interface"
		*s1.PStruct = s1.Struct
		*s1.MapPInt[1] = 1
		*s1.MapPInt[2] = 2
		s2 := Copy(s1).(st)
		t.Logf("%+v %+v", s1, s2)
		if !reflect.DeepEqual(s1, s2) {
			t.Fatal("fail")
		}
	}
}

func TestDeepCopyGeneric(t *testing.T) {
	s1 := st{
		Int:        10,
		PInt:       new(int),
		String:     "string",
		PString:    new(string),
		ArrayInt:   [2]int{1, 2},
		SliceInt:   []int{1, 2},
		MapInt:     map[int]int{1: 1, 2: 2},
		MapPInt:    map[int]*int{1: new(int), 2: new(int)},
		Interface:  "interface",
		PInterface: new(interface{}),
		Struct: stt{
			Int:        20,
			PInt:       nil,
			String:     "string2",
			PString:    nil,
			SliceInt:   []int{3, 4},
			MapInt:     map[int]int{3: 3, 4: 4},
			Interface:  "interface2",
			PInterface: nil,
		},
		PStruct:      new(stt),
		ArrayStruct:  [2]st1{{1}, {2}},
		SliceStruct:  []st1{{1}, {2}},
		SlicePStruct: []*st1{{1}, {2}},
	}
	*s1.PInt = 30
	*s1.PString = "string"
	*s1.PInterface = "interface"
	*s1.PStruct = s1.Struct
	*s1.MapPInt[1] = 1
	*s1.MapPInt[2] = 2
	s2 := CopyGeneric(s1)
	t.Logf("%+v %+v", s1, s2)
	if !reflect.DeepEqual(s1, s2) {
		t.Fatal("fail")
	}
}
