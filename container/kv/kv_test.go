package kv

import (
	"testing"
)

// 测试用键, 需复用同一实例才能命中同一条记录.
var (
	testKeyInt    = NewK[int]("int")
	testKeyString = NewK[string]("string")
	testKeyPtr    = NewK[*int]("ptr")
	testKeySlice  = NewK[[]int]("slice")
)

// mustPanic 断言 fn 会 panic.
func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("期待 panic, 但未发生")
		}
	}()
	fn()
}

// TestM_PlainKey 覆盖非 K 的普通键.
func TestM_PlainKey(t *testing.T) {
	m := NewM()
	m.Set("a", 1)
	m.Set(2, "b")

	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v", v, ok)
	}
	if v, ok := m.Get(2); !ok || v != "b" {
		t.Fatalf("Get(2) = %v, %v", v, ok)
	}
	// 键不存在时返回 nil, false
	if v, ok := m.Get("not-exist"); ok || v != nil {
		t.Fatalf("Get(not-exist) = %v, %v", v, ok)
	}
}

// TestNewM 覆盖 NewM 的构造与参数校验.
func TestNewM(t *testing.T) {
	m := NewM(testKeyInt, 1, "b", 2)
	if m.Len() != 2 {
		t.Fatalf("Len = %d, 期待 2", m.Len())
	}
	if v, ok := testKeyInt.Get(m); !ok || v != 1 {
		t.Fatalf("testKeyInt.Get = %v, %v", v, ok)
	}
	if v, ok := m.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b) = %v, %v", v, ok)
	}

	// 参数个数为奇数时 panic
	mustPanic(t, func() { NewM(testKeyInt) })
	// 参数个数为偶数但为空时正常
	if m := NewM(); m.Len() != 0 {
		t.Fatalf("Len = %d, 期待 0", m.Len())
	}
}

// TestK_KeyIdentity 键的判等基于实例身份而非名称.
func TestK_KeyIdentity(t *testing.T) {
	k1 := NewK[int]("same")
	k2 := NewK[int]("same")

	// 名称相同但实例不同, 是两个不同的键
	if k1 == k2 {
		t.Fatal("同名的两个键不应相等")
	}

	m := NewM(k1, 1, k2, 2)
	if m.Len() != 2 {
		t.Fatalf("Len = %d, 期待 2", m.Len())
	}
	if v, ok := k1.Get(m); !ok || v != 1 {
		t.Fatalf("k1.Get = %v, %v", v, ok)
	}
	if v, ok := k2.Get(m); !ok || v != 2 {
		t.Fatalf("k2.Get = %v, %v", v, ok)
	}

	// 值类型不同、名称相同的键也是不同的键
	ki, ks := NewK[int]("x"), NewK[string]("x")
	if m := NewM(ki, 1, ks, "a"); m.Len() != 2 {
		t.Fatalf("Len = %d, 期待 2", m.Len())
	}
}

// TestK_GetNotExist 键不存在时返回零值与 false.
func TestK_GetNotExist(t *testing.T) {
	m := NewM()
	if v, ok := testKeyInt.Get(m); ok || v != 0 {
		t.Fatalf("Get = %v, %v", v, ok)
	}
	if v, ok := testKeyString.Get(m); ok || v != "" {
		t.Fatalf("Get = %v, %v", v, ok)
	}
}

// TestK_SetGet 覆盖 K 的正常存取与更新.
func TestK_SetGet(t *testing.T) {
	m := NewM()
	m.Set(testKeyInt, 1)
	m.Set(testKeyString, "a")

	if v, ok := testKeyInt.Get(m); !ok || v != 1 {
		t.Fatalf("testKeyInt.Get = %v, %v", v, ok)
	}
	if v, ok := testKeyString.Get(m); !ok || v != "a" {
		t.Fatalf("testKeyString.Get = %v, %v", v, ok)
	}

	// 更新已有键
	m.Set(testKeyInt, 2)
	if v, ok := testKeyInt.Get(m); !ok || v != 2 {
		t.Fatalf("更新后 testKeyInt.Get = %v, %v", v, ok)
	}
}

// TestK_NilValue 覆盖 nil 值的处理.
func TestK_NilValue(t *testing.T) {
	m := NewM()

	// V 可赋 nil 时允许存 nil
	m.Set(testKeyPtr, nil)
	if v, ok := testKeyPtr.Get(m); !ok || v != nil {
		t.Fatalf("Get = %v, %v", v, ok)
	}

	// 带类型的 nil 指针同样视为合法值
	m.Set(testKeyPtr, (*int)(nil))
	if v, ok := testKeyPtr.Get(m); !ok || v != nil {
		t.Fatalf("Get = %v, %v", v, ok)
	}
	// 直接走 M.Get 时, 接口内持有带类型的 nil 指针
	if v, ok := m.Get(testKeyPtr); !ok || v != (*int)(nil) {
		t.Fatalf("M.Get = %v, %v", v, ok)
	}

	// V 不可赋 nil 时存 nil 会 panic
	mustPanic(t, func() { m.Set(testKeyInt, nil) })
}

// TestK_AnyValueType 值类型为接口时接受任意值.
func TestK_AnyValueType(t *testing.T) {
	k := NewK[any]("any")
	m := NewM(k, "x")
	if v, ok := k.Get(m); !ok || v != "x" {
		t.Fatalf("Get = %v, %v", v, ok)
	}
}

// TestK_TypeCheck 值类型不匹配时 panic.
func TestK_TypeCheck(t *testing.T) {
	m := NewM()
	m.Set(testKeyInt, 1)

	mustPanic(t, func() { m.Set(testKeyInt, "1") })
	mustPanic(t, func() { m.Set(testKeyString, 1) })
}

// TestM_GetTypeMismatchPanic 读取时也会做类型校验.
func TestM_GetTypeMismatchPanic(t *testing.T) {
	// 复用同一内部指针但改变值类型, 构造出存储值类型与键不匹配的情况
	ki := NewK[int]("forged")
	ks := K[string]{string: ki.string}

	m := NewM(ki, 1)
	mustPanic(t, func() { m.Get(ks) })
}

// TestK_ZeroValuePanic 未初始化的零值键不可用, 避免内部指针为 nil 时互相冲突.
func TestK_ZeroValuePanic(t *testing.T) {
	var k K[int]
	var m M

	mustPanic(t, func() { m.Set(k, 1) })
	mustPanic(t, func() { m.Get(k) })
	mustPanic(t, func() { k.Get(&m) })
}

// TestM_ZeroValue M 的零值可直接使用.
func TestM_ZeroValue(t *testing.T) {
	var m M
	if m.Len() != 0 {
		t.Fatalf("Len = %d, 期待 0", m.Len())
	}
	if _, ok := m.Get("a"); ok {
		t.Fatal("零值 M 不应有数据")
	}

	m.Set("a", 1)
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v", v, ok)
	}
	if m.Len() != 1 {
		t.Fatalf("Len = %d, 期待 1", m.Len())
	}
}

// TestM_Len 覆盖 Len.
func TestM_Len(t *testing.T) {
	m := NewM()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("a", 3) // 覆盖已有键, 数量不变
	if m.Len() != 2 {
		t.Fatalf("Len = %d, 期待 2", m.Len())
	}
}

// TestM_Clone 覆盖浅拷贝.
func TestM_Clone(t *testing.T) {
	p := new(int)
	m := NewM(testKeyInt, 1, "b", 2, testKeyPtr, p)
	c := m.Clone()

	if c.Len() != m.Len() {
		t.Fatalf("Len = %d, 期待 %d", c.Len(), m.Len())
	}
	if v, ok := testKeyInt.Get(c); !ok || v != 1 {
		t.Fatalf("副本 Get = %v, %v", v, ok)
	}
	// 浅拷贝: 引用类型的值与原映射共享
	if v, ok := c.Get(testKeyPtr); !ok || v != p {
		t.Fatalf("副本 Get = %v, %v", v, ok)
	}

	// 两个映射相互独立
	c.Set(testKeyInt, 100)
	if v, _ := testKeyInt.Get(m); v != 1 {
		t.Fatalf("修改副本影响了原映射: %v", v)
	}
	m.Set("b", 200)
	if v, _ := c.Get("b"); v != 2 {
		t.Fatalf("修改原映射影响了副本: %v", v)
	}

	// 零值映射的副本可直接使用
	var zero M
	z := zero.Clone()
	if z.Len() != 0 {
		t.Fatalf("Len = %d, 期待 0", z.Len())
	}
	z.Set("a", 1)
	if v, ok := z.Get("a"); !ok || v != 1 {
		t.Fatalf("Get = %v, %v", v, ok)
	}
}

// TestM_DeepClone 覆盖深拷贝: 值被递归拷贝, 键与原映射共享.
func TestM_DeepClone(t *testing.T) {
	s := []int{1, 2, 3}
	p := new(int)
	*p = 10
	m := NewM(testKeySlice, s, testKeyPtr, p)

	c := m.DeepClone()
	if c.Len() != m.Len() {
		t.Fatalf("Len = %d, 期待 %d", c.Len(), m.Len())
	}

	// 键与原映射共享: 同一个 K 实例仍能命中副本
	cs, ok := testKeySlice.Get(c)
	if !ok || len(cs) != 3 || cs[0] != 1 {
		t.Fatalf("副本 Get = %v, %v", cs, ok)
	}

	// 深拷贝: slice 值不共享底层数组
	cs[0] = 100
	if v, _ := testKeySlice.Get(m); v[0] != 1 {
		t.Fatalf("修改副本的 slice 影响了原映射: %v", v)
	}

	// 深拷贝: 指针值指向不同的对象
	cp, ok := testKeyPtr.Get(c)
	if !ok || cp == nil || cp == p || *cp != 10 {
		t.Fatalf("副本 Get = %v, %v", cp, ok)
	}
	*cp = 20
	if *p != 10 {
		t.Fatalf("修改副本的指针影响了原映射: %v", *p)
	}

	// 零值映射的深拷贝可直接使用
	var zero M
	z := zero.DeepClone()
	if z.Len() != 0 {
		t.Fatalf("Len = %d, 期待 0", z.Len())
	}
	z.Set("a", 1)
	if v, ok := z.Get("a"); !ok || v != 1 {
		t.Fatalf("Get = %v, %v", v, ok)
	}
}

// deepCopyCounter 实现 deepcopy.CopyGenericI, 用于验证 DeepClone 优先调用 DeepCopy.
type deepCopyCounter struct {
	calls int
}

func (d *deepCopyCounter) DeepCopy() interface{} {
	return &deepCopyCounter{calls: d.calls + 1}
}

// TestM_DeepClone_CopyGenericI 值的类型实现 deepcopy.CopyGenericI 时优先调用其 DeepCopy.
func TestM_DeepClone_CopyGenericI(t *testing.T) {
	k := NewK[*deepCopyCounter]("counter")
	m := NewM(k, &deepCopyCounter{})

	v, ok := k.Get(m.DeepClone())
	if !ok || v == nil {
		t.Fatalf("Get = %v, %v", v, ok)
	}
	if v.calls != 1 {
		t.Fatalf("DeepCopy 调用次数 = %d, 期待 1", v.calls)
	}
}

// TestM_Clear 覆盖清空.
func TestM_Clear(t *testing.T) {
	m := NewM(testKeyInt, 1, "b", 2)
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("Len = %d, 期待 0", m.Len())
	}
	// 清空后仍可继续使用
	m.Set("a", 1)
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v", v, ok)
	}

	// 零值映射清空不 panic
	var zero M
	zero.Clear()
}

// TestM_Del 覆盖删除.
func TestM_Del(t *testing.T) {
	m := NewM(testKeyInt, 1, "b", 2)

	// 删除普通键
	if v, ok := m.Del("b"); !ok || v != 2 {
		t.Fatalf("Del(b) = %v, %v", v, ok)
	}
	if _, ok := m.Get("b"); ok {
		t.Fatal("删除后仍存在")
	}
	// 删除不存在的普通键
	if v, ok := m.Del("b"); ok || v != nil {
		t.Fatalf("Del(b) = %v, %v", v, ok)
	}

	// 删除 K 键
	if v, ok := m.Del(testKeyInt); !ok || v != 1 {
		t.Fatalf("Del(testKeyInt) = %v, %v", v, ok)
	}
	if _, ok := testKeyInt.Get(m); ok {
		t.Fatal("删除后仍存在")
	}
	if m.Len() != 0 {
		t.Fatalf("Len = %d, 期待 0", m.Len())
	}
	// 删除不存在的 K 键
	if v, ok := m.Del(testKeyInt); ok || v != nil {
		t.Fatalf("Del(testKeyInt) = %v, %v", v, ok)
	}

	// 零值映射删除普通键
	var zero M
	if v, ok := zero.Del("a"); ok || v != nil {
		t.Fatalf("Del(a) = %v, %v", v, ok)
	}
}

// TestM_DelTypeMismatchPanic 用伪造的键删除时会做类型校验.
func TestM_DelTypeMismatchPanic(t *testing.T) {
	ki := NewK[int]("forged")
	ks := K[string]{string: ki.string}
	m := NewM(ki, 1)
	mustPanic(t, func() { m.Del(ks) })
}

// TestK_String 覆盖键的字符串描述与零值守卫.
func TestK_String(t *testing.T) {
	if s := NewK[int]("n").String(); s != "n(int)" {
		t.Fatalf("String = %q, 期待 %q", s, "n(int)")
	}
	// 零值键未初始化, String 会 panic
	var zk K[int]
	mustPanic(t, func() { _ = zk.String() })
}

// TestK_Set 通过 K.Set 写入.
func TestK_Set(t *testing.T) {
	m := NewM()
	testKeyInt.Set(m, 1)
	if v, ok := testKeyInt.Get(m); !ok || v != 1 {
		t.Fatalf("Get = %v, %v", v, ok)
	}
}

// TestK_Must 覆盖 Must: 存在时返回值, 不存在时 panic, nil 值时返回零值.
func TestK_Must(t *testing.T) {
	m := NewM(testKeyInt, 1)
	if v := testKeyInt.Must(m); v != 1 {
		t.Fatalf("Must = %v, 期待 1", v)
	}
	mustPanic(t, func() { testKeyString.Must(m) })

	// 存储的值为 nil 属于合法值, Must 应返回零值而非 panic
	m.Set(testKeyPtr, nil)
	if v := testKeyPtr.Must(m); v != nil {
		t.Fatalf("Must = %v, 期待 nil", v)
	}
}

// TestGetMustFunc 覆盖包级 Get/Must.
func TestGetMustFunc(t *testing.T) {
	m := NewM(testKeyInt, 1)
	if v, ok := Get(m, testKeyInt); !ok || v != 1 {
		t.Fatalf("Get = %v, %v", v, ok)
	}
	if v := Must(m, testKeyInt); v != 1 {
		t.Fatalf("Must = %v, 期待 1", v)
	}
	mustPanic(t, func() { Must(m, testKeyString) })
}

// TestM_Must 覆盖 M.Must.
func TestM_Must(t *testing.T) {
	m := NewM(testKeyInt, 1, "b", 2)
	if v := m.Must(testKeyInt); v != 1 {
		t.Fatalf("Must = %v, 期待 1", v)
	}
	if v := m.Must("b"); v != 2 {
		t.Fatalf("Must = %v, 期待 2", v)
	}
	// 键不存在时 panic
	mustPanic(t, func() { m.Must("not-exist") })
	// 存储 nil 值时正常返回
	m.Set(testKeyPtr, nil)
	if v := m.Must(testKeyPtr); v != nil {
		t.Fatalf("Must = %v, 期待 nil", v)
	}
}

// TestMI_KIImpl 确认 *M 实现 MI, K 实现 KI.
func TestMI_KIImpl(t *testing.T) {
	var _ MI = &M{}
	var _ KI = K[int]{}
}

// TestM_SetPairs 覆盖批量设置.
func TestM_SetPairs(t *testing.T) {
	m := NewM()
	m.SetPairs(testKeyInt, 1, "b", 2)
	if m.Len() != 2 {
		t.Fatalf("Len = %d, 期待 2", m.Len())
	}
	if v, ok := testKeyInt.Get(m); !ok || v != 1 {
		t.Fatalf("Get = %v, %v", v, ok)
	}
	// 参数个数为奇数时 panic
	mustPanic(t, func() { m.SetPairs(testKeyInt) })
}

// TestCanBeNil 覆盖类型可赋 nil 的判断.
func TestCanBeNil(t *testing.T) {
	if !canBeNil[any]() || !canBeNil[*int]() || !canBeNil[[]int]() ||
		!canBeNil[map[int]int]() || !canBeNil[chan int]() || !canBeNil[func()]() {
		t.Error("可赋 nil 的类型判断错误")
	}
	if canBeNil[int]() || canBeNil[string]() || canBeNil[struct{}]() || canBeNil[[2]int]() {
		t.Error("不可赋 nil 的类型判断错误")
	}
}

// benchK、benchM 为基准测试的输出 sink, 用于阻止编译器优化掉被测调用.
var (
	benchK K[int]
	benchM *M
)

func BenchmarkNewK(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchK = NewK[int]("bench")
	}
}

func BenchmarkM_SetK(b *testing.B) {
	k := NewK[int]("set")
	m := NewM(k, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(k, i)
	}
}

func BenchmarkM_GetK(b *testing.B) {
	k := NewK[int]("get")
	m := NewM(k, 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := m.Get(k); !ok {
			b.Fatal("键不存在")
		}
	}
}

func BenchmarkK_Get(b *testing.B) {
	k := NewK[int]("get")
	m := NewM(k, 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := k.Get(m); !ok {
			b.Fatal("键不存在")
		}
	}
}

// BenchmarkM_SetPlain 测量普通键的覆盖写, 与 SetK 对齐(同一键反复写入, 避免 map 无界增长).
func BenchmarkM_SetPlain(b *testing.B) {
	m := NewM(1, 0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(1, i)
	}
}

func BenchmarkM_GetPlain(b *testing.B) {
	m := NewM(1, 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := m.Get(1); !ok {
			b.Fatal("键不存在")
		}
	}
}

// BenchmarkK_GetNil 覆盖值为 nil 的读取路径, 该路径会判断 V 是否可赋 nil.
func BenchmarkK_GetNil(b *testing.B) {
	k := NewK[*int]("nil")
	m := NewM(k, nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := k.Get(m); !ok {
			b.Fatal("键不存在")
		}
	}
}

// BenchmarkM_Clone 测量浅拷贝.
func BenchmarkM_Clone(b *testing.B) {
	m := NewM("a", 1, "b", 2, "c", 3, "d", 4)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchM = m.Clone()
	}
}

// BenchmarkM_DeepClone 测量深拷贝.
func BenchmarkM_DeepClone(b *testing.B) {
	m := NewM("a", 1, "b", 2, "c", 3, "d", 4)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchM = m.DeepClone()
	}
}
