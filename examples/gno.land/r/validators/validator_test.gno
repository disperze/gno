package validators

import (
	"std"

	"testing"
)

func Test(t *testing.T) {
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
