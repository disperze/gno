package validators

import (
	"std"

	"gno.land/p/avl"
	"gno.land/p/ufmt"
)

func CreateValidator(name, pubkey string) int {
	std.AssertOriginCall()

	if gValidators.Size()+1 > maxValidators {
		panic("validator set is full")
	}

	caller := std.GetOrigCaller()
	valKey := validatorIDKey(caller)
	exists := gValidators.Has(valKey)
	if exists {
		panic("validator already exists")
	}

	vp := std.GetOrigSend().AmountOf(nativeDenom)
	validator := newValidator(caller, name, pubkey, vp, std.GetHeight())

	gValidators.Set(valKey, validator)
	return gValidators.Size()
}

func Stake(validator std.Address) {
	std.AssertOriginCall()

	valoper := getValidator(validator)
	if valoper == nil {
		panic("validator does not exist")
	}

	vp := std.GetOrigSend().AmountOf(nativeDenom)
	valoper.AddVotingPower(vp)

	// TODO: add delegation
}

func ValidatorSet() string {
	colon := ""
	output := ""
	gValidators.Iterate("", "", func(n *avl.Tree) bool {
		validator := n.Value().(*Validator)
		output = output + colon + validator.Json()
		colon = ","
		return false
	})

	return "[" + output + "]"
}

func UpdateMaxValidators(max int) {
	std.AssertOriginCall()
	caller := std.GetOrigCaller()
	if caller != admin {
		panic("unauthorized")
	}

	if max <= maxValidators {
		panic("new max must be greater than current value")
	}

	maxValidators = max
}

func Render(path string) string {
	str := ""
	str += ufmt.Sprintf("## Total Validators: %d | Max: %d\n\n", gValidators.Size(), maxValidators)
	str += "|Name|Address|Voting Power|Joined|\n"
	str += "|---|---|---|---|\n"
	total := int64(0)
	gValidators.Iterate("", "", func(n *avl.Tree) bool {
		val := n.Value().(*Validator)
		str += val.Render() + "\n"
		total += val.vp
		return false
	})
	str += ufmt.Sprintf("\n\nTotal Staked: %d\n", total)

	return str
}
