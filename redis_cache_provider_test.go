package cache

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type Person struct {
	Name string
	Age  int
}

type Data struct {
	Bool bool

	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uintptr uintptr

	Float32 float32
	Float64 float64

	Array [2]int
	Map   map[string]int

	Slice  []int
	String string

	Person Person

	Time     time.Time
	UnixTime UnixTime
}

var data = Data{
	Bool:    true,
	Int:     1,
	Int8:    2,
	Int16:   3,
	Int32:   4,
	Int64:   5,
	Uint:    6,
	Uint8:   7,
	Uint16:  8,
	Uint32:  9,
	Uint64:  10,
	Uintptr: 11,
	Float32: 12.12,
	Float64: 13.13,

	Array:    [2]int{16, 16},
	Map:      map[string]int{"A": 17, "B": 18},
	Slice:    []int{19, 20},
	String:   "21",
	Person:   Person{"Jerry", 22},
	Time:     time.Now(),
	UnixTime: UnixTime(time.Now()),
}

func TestRedisCache(t *testing.T) {
	cli := NewRedisCacheProvider(redisC)
	dv, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(dv))

	// 确保key不存在
	key := "go.redis.test.cli"
	cli.MustRemove(key)

	// 获取不存在的缓存
	var v Data = Data{Bool: true} //设置一个非 0值
	if cli.MustGet(key, &v); !v.Bool {
		t.Error("在没有缓存key的情况下，接收者值被改变")
	}

	// 尝试获取不存在的缓存
	if exist, err := cli.TryGet(key, &v); exist || err != nil || !v.Bool {
		t.Error("预期之外的行为", exist, v.Bool, err)
	}

	// 设置缓存
	cli.MustSet(key, data, 0)

	// 获取缓存值，对比是否和之前的一样。
	cli.MustGet(key, &v)
	tv, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	if string(dv) != string(tv) {
		t.Error("存进去、取出来 值发生了改变")
	}

	// SETNX
	result, err := cli.Create(key, Data{Int: 99}, 0)
	if err != nil {
		t.Error(err)
	}

	if result {
		t.Error("预期 key 存在，添加缓存失败")
	}

	cli.MustGet(key, &v)
	if v.Int == 99 {
		t.Error("预期 key 存在，添加缓存失败，缓存值还是之前的值")
	}

	// 删除缓存key， 存在和不存在时的不同返回值。
	r1 := cli.MustRemove(key)
	r2 := cli.MustRemove(key)
	if !r1 || r2 {
		t.Error("预期 r1 == true , r2 == false")
	}
}

func TestNewRedisCacheProvider(t *testing.T) {
	type args struct {
		cli redis.Cmdable
	}
	tests := []struct {
		name    string
		args    args
		want    *RedisCacheProvider
		wantErr bool
	}{
		{"Client", args{&redis.Client{}}, &RedisCacheProvider{&redis.Client{}}, false},
		{"Ring client", args{&redis.Ring{}}, &RedisCacheProvider{&redis.Ring{}}, false},
		{"ClusterClient client", args{&redis.ClusterClient{}}, &RedisCacheProvider{&redis.ClusterClient{}}, false},
		{"Tx client", args{&redis.Tx{}}, &RedisCacheProvider{&redis.Tx{}}, false},
		{"nil value", args{nil}, nil, true},
		// 哨兵客户端，不支持。
		//{"client", args{&redis.SentinelClient{}}, &RedisCacheProvider{&redis.Client{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := func() (p *RedisCacheProvider, e error) {
				defer func() {
					err := recover()
					if err != nil {
						e = err.(error)
						p = nil
					}
				}()
				return NewRedisCacheProvider(tt.args.cli), nil
			}()

			if !reflect.DeepEqual(got, tt.want) || (err != nil) != tt.wantErr {
				t.Errorf("NewRedisCacheProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_Get(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.Get(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisCacheProvider_MustGet(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.MustGet(tt.args.key, tt.args.value)
		})
	}
}

func TestRedisCacheProvider_TryGet(t *testing.T) {
	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.TryGet(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.TryGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCacheProvider.TryGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_Create(t *testing.T) {
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.Create(tt.args.key, tt.args.value, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCacheProvider.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_MustCreate(t *testing.T) {
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cli.MustCreate(tt.args.key, tt.args.value, tt.args.t); got != tt.want {
				t.Errorf("RedisCacheProvider.MustCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_Set(t *testing.T) {
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cli.Set(tt.args.key, tt.args.value, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisCacheProvider_MustSet(t *testing.T) {
	type args struct {
		key   string
		value any
		t     time.Duration
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cli.MustSet(tt.args.key, tt.args.value, tt.args.t)
		})
	}
}

func TestRedisCacheProvider_Remove(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		cli     *RedisCacheProvider
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.Remove(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisCacheProvider.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisCacheProvider.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedisCacheProvider_MustRemove(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		cli  *RedisCacheProvider
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cli.MustRemove(tt.args.key); got != tt.want {
				t.Errorf("RedisCacheProvider.MustRemove() = %v, want %v", got, tt.want)
			}
		})
	}
}