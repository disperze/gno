package validators

import (
	"std"

	"testing"
)

func TestCreateValidator(t *testing.T) {
	pkgAddr := std.TestDerivePkgAddr("gno.land/r/validators")
	caller := std.Address("g1u7y667z64x2h7vc6fmpcprgey4ck233jaww9zq")
	std.TestSetOrigPkgAddr(pkgAddr)
	std.TestSetOrigSend(std.Coins{{"ugnot", 10}}, nil)
	std.TestSetOrigCaller(caller)

	pubkey := "cTJh+fxKUSHGVLDKSSn5zShXGWle5QYyIkRpg7FAbI8="
	id := CreateValidator("test", pubkey)
	if id != 1 {
		t.Fatalf("id expected %q, got %q.", 1, id)
	}
}

func TestUpdateMaxValidators(t *testing.T) {
	pkgAddr := std.TestDerivePkgAddr("gno.land/r/validators")
	anon := std.Address("g1us8428u2a5satrlxzagqqa5m6vmuze025anjlj")
	admin := std.Address("g19yltltjed5688sullftzurhdshmjwtm8vlqdr2")
	std.TestSetOrigPkgAddr(pkgAddr)

	std.TestSetOrigCaller(anon)
	assertPanic(t, func() {
		UpdateMaxValidators(15)
	})

	std.TestSetOrigCaller(admin)
	assertPanic(t, func() {
		UpdateMaxValidators(5)
	})

	std.TestSetOrigCaller(admin)
	UpdateMaxValidators(100)
}

func assertPanic(t *testing.T, f func()) {
	defer func() { _ = recover() }()
    f()
	t.Errorf("The code did not panic")
}
