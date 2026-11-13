package deepcopy

import (
	"reflect"
	"testing"
	"time"
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

// unexportedSt 含未导出字段, 用于验证未导出字段也能被深拷贝.
type unexportedSt struct {
	Exported   int
	unexported int
	ptr        *int
	slice      []int
	m          map[int]int
	inner      stt
}

func TestDeepCopyUnexported(t *testing.T) {
	p := new(int)
	*p = 5
	s1 := unexportedSt{
		Exported:   1,
		unexported: 2,
		ptr:        p,
		slice:      []int{1, 2},
		m:          map[int]int{1: 1},
		inner:      stt{Int: 9, SliceInt: []int{3}},
	}
	s2 := Copy(s1).(unexportedSt)

	// 基本与嵌套字段被拷贝
	if s2.Exported != 1 || s2.unexported != 2 {
		t.Fatalf("基本字段拷贝失败: %+v", s2)
	}
	if s2.inner.Int != 9 || s2.inner.SliceInt[0] != 3 {
		t.Fatalf("嵌套结构体拷贝失败: %+v", s2.inner)
	}

	// 引用类型为深拷贝, 不共享底层数据
	if s2.ptr == s1.ptr || *s2.ptr != 5 {
		t.Fatalf("指针未深拷贝: %v", s2.ptr)
	}
	s2.slice[0] = 100
	if s1.slice[0] != 1 {
		t.Fatalf("切片未深拷贝: %v", s1.slice)
	}
	s2.m[1] = 100
	if s1.m[1] != 1 {
		t.Fatalf("map 未深拷贝: %v", s1.m)
	}
	s2.inner.SliceInt[0] = 100
	if s1.inner.SliceInt[0] != 3 {
		t.Fatalf("嵌套切片未深拷贝: %v", s1.inner.SliceInt)
	}
}

func TestDeepCopyTime(t *testing.T) {
	type wrapper struct{ T time.Time }
	now := time.Now()
	s1 := wrapper{T: now}
	s2 := Copy(s1).(wrapper)
	if !s2.T.Equal(now) {
		t.Fatalf("time.Time 未正确拷贝: %v", s2.T)
	}
}
