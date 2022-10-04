package validators

import (
	"gno.land/p/avl"
)

//----------------------------------------
// Realm (package) state

var (
	// address -> *Validator
	gValidators = avl.NewMutTree()
)
