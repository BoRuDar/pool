package pool

import (
	"encoding/json"
	"math/rand/v2"
	"reflect"
	"testing"
	"time"
)

type userTestStr struct {
	Active bool
	Name   string
	Age    int
	Items  []string
	Spouse *userTestStr
}

var (
	even = userTestStr{
		Active: true,
		Name:   "even",
		Age:    22,
		Items:  []string{"even"},
	}
	odd = userTestStr{
		Active: false,
		Name:   "odd",
		Age:    33,
		Items:  []string{"odd"},
	}
)

func BenchmarkNoPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf []byte

		for n := range rand.IntN(10) {
			if n/2 == 0 {
				buf = toJSON(&userTestStr{
					Active: even.Active,
					Name:   even.Name,
					Age:    n + 2,
					Items:  even.Items,
				})
			} else {
				buf = toJSON(&userTestStr{
					Active: odd.Active,
					Name:   odd.Name,
					Age:    n + 4,
					Items:  odd.Items,
				})
			}
		}

		if len(buf) > 0 {
			buf = nil
		}
	}
}

func BenchmarkWithPool(b *testing.B) {
	pool := New[userTestStr]()

	for i := 0; i < b.N; i++ {
		var buf []byte

		for n := range rand.IntN(10) {
			if n/2 == 0 {
				u1 := pool.Get()
				u1.Age = n + 2
				u1.Name = even.Name
				u1.Items = even.Items
				u1.Active = even.Active

				buf = toJSON(u1)
				pool.Return(u1)
			} else {
				u2 := pool.Get()
				u2.Age = n + 4
				u2.Name = odd.Name
				u2.Items = odd.Items
				u2.Active = odd.Active

				buf = toJSON(u2)
				pool.Return(u2)
			}
		}

		if len(buf) > 0 {
			buf = nil
		}
	}
}

func toJSON(ptr *userTestStr) []byte {
	b, _ := json.Marshal(ptr)
	return b
}

func TestConcurrentAccess(t *testing.T) {
	t.Parallel()

	pool := New[userTestStr]()

	tcs := []struct {
		name     string
		expected userTestStr
		setFn    func(str *userTestStr)
	}{
		{
			name:     "empty struct",
			expected: userTestStr{},
			setFn: func(str *userTestStr) {
				// do nothing here
			},
		},
		{
			name: "name",
			expected: userTestStr{
				Name: "one",
				Age:  42,
			},
			setFn: func(u *userTestStr) {
				u.Name = "one"
				u.Age = 42
			},
		},
		{
			name: "name and age",
			expected: userTestStr{
				Name: "one",
			},
			setFn: func(u *userTestStr) {
				u.Name = "one"
			},
		},
		{
			name: "embedded struct",
			expected: userTestStr{
				Name: "one",
				Spouse: &userTestStr{
					Name: "two",
				},
			},
			setFn: func(u *userTestStr) {
				u.Name = "one"
				u.Spouse = &userTestStr{
					Name: "two",
				}
			},
		},
		{
			name: "embedded struct and slice",
			expected: userTestStr{
				Name: "one",
				Spouse: &userTestStr{
					Name:  "two",
					Items: []string{"keys"},
				},
				Items: []string{"phone"},
			},
			setFn: func(u *userTestStr) {
				u.Name = "one"
				u.Spouse = &userTestStr{
					Name:  "two",
					Items: []string{"keys"},
				}
				u.Items = []string{"phone"}
			},
		},
		{
			name: "all fields",
			expected: userTestStr{
				Name:   "one",
				Age:    42,
				Active: true,
				Spouse: &userTestStr{
					Name:  "two",
					Items: []string{"keys"},
				},
				Items: []string{"phone"},
			},
			setFn: func(u *userTestStr) {
				u.Name = "one"
				u.Age = 42
				u.Active = true
				u.Spouse = &userTestStr{
					Name:  "two",
					Items: []string{"keys"},
				}
				u.Items = []string{"phone"}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			time.Sleep(time.Millisecond * time.Duration(rand.IntN(10)))

			u := pool.Get()
			defer pool.Return(u)

			tc.setFn(u)

			if !reflect.DeepEqual(*u, tc.expected) {
				t.Fatalf("structs didn't match: %#v != %#v", u, tc.expected)
			}
		})
	}
}
