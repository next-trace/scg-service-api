package di_test

import (
	"testing"

	di "github.com/next-trace/scg-service-api/infrastructure/di"
)

type A struct{ Name string }

type B struct{ A A }

func newA() A    { return A{Name: "a"} }
func newB(a A) B { return B{A: a} }

func TestContainer_ProvideResolveInvoke(t *testing.T) {
	c := di.NewContainer()
	if err := c.Provide(newA); err != nil {
		t.Fatalf("provide A: %v", err)
	}
	if err := c.Provide(newB); err != nil {
		t.Fatalf("provide B: %v", err)
	}

	var gotB B
	if err := c.Resolve(&gotB); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if gotB.A.Name != "a" {
		t.Fatalf("unexpected resolved value: %+v", gotB)
	}

	// Invoke a function that needs B
	called := false
	fn := func(b B) {
		if b.A.Name == "a" {
			called = true
		}
	}
	if err := c.Invoke(fn); err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if !called {
		t.Fatalf("expected function to be called with resolved deps")
	}

	c.Reset()
}
