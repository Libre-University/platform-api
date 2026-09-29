package core

import (
	"context"
	"io/fs"
	"testing"

	"github.com/go-chi/chi/v5"
)

type stub struct{ name string }

func (s stub) Name() string                          { return s.name }
func (s stub) Version() string                       { return "0.0.0" }
func (s stub) Migrations() fs.FS                     { return nil }
func (s stub) Permissions() []Permission             { return nil }
func (s stub) Init(context.Context, *Services) error { return nil }
func (s stub) Routes(chi.Router)                     {}

func TestValidName(t *testing.T) {
	ok := []string{"obs", "lms", "campus_card", "a1"}
	bad := []string{"", "O", "Obs", "1obs", "obs-lms", "a", "çok", "this_name_is_way_too_long_for_a_schema"}
	for _, n := range ok {
		if !ValidName(n) {
			t.Errorf("%q geçerli olmalıydı", n)
		}
	}
	for _, n := range bad {
		if ValidName(n) {
			t.Errorf("%q geçersiz olmalıydı", n)
		}
	}
}

func TestRegisterSortsAndRejectsDuplicates(t *testing.T) {
	resetRegistry()
	t.Cleanup(resetRegistry)

	Register(stub{"lms"})
	Register(stub{"obs"})
	got := Registered()
	if len(got) != 2 || got[0].Name() != "lms" || got[1].Name() != "obs" {
		t.Fatalf("beklenmeyen sıra: %v", got)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("tekrar eden ad panic üretmeliydi")
		}
	}()
	Register(stub{"obs"})
}

func TestRegisterRejectsInvalidName(t *testing.T) {
	resetRegistry()
	t.Cleanup(resetRegistry)
	defer func() {
		if recover() == nil {
			t.Fatal("geçersiz ad panic üretmeliydi")
		}
	}()
	Register(stub{"Not-Valid"})
}
