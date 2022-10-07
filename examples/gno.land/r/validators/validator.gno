package validators

import (
	"encoding/base64"
	"std"

	"gno.land/p/ufmt"
)

type Validator struct {
	name    string
	address std.Address
	pubkey  string
	vp      int64
	joined  int64
}

func newValidator(address std.Address, name, pubkey string, vp, joined int64) *Validator {
	if name == "" || len(name) > 70 {
		panic("invalid validator name")
	}

	pubk, err := base64.StdEncoding.DecodeString(pubkey)
	if err != nil || len(pubk) != 32 {
		panic("invalid pubkey")
	}

	if vp < 1 {
		panic("amount must be positive")
	}

	return &Validator{
		name:    name,
		address: address,
		pubkey:  pubkey,
		vp:      vp,
		joined:  std.GetHeight(),
	}
}

func (v *Validator) AddVotingPower(vp int64) {
	if vp < 1 {
		panic("stake must be positive")
	}

	v.vp += vp
}

func (v *Validator) Json() string {
	format := `{"name": "%s", "address": "%s", "pubkey": "%s", "vp": %d}`

	return ufmt.Sprintf(format, v.name, v.address.String(), v.pubkey, v.vp)
}

func (v *Validator) Render() string {
	format := `|%s|%s|%d|%d|`

	return ufmt.Sprintf(format, v.name, v.address.String(), v.vp, v.joined)
}

func validatorIDKey(address std.Address) string {
	return address.String()
}

func getValidator(val std.Address) *Validator {
	valkey := validatorIDKey(val)
	valoper, exists := gValidators.Get(valkey)
	if !exists {
		return nil
	}

	return valoper.(*Validator)
}
