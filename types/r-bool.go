package types

// RBool is a boolean which supports R's NA
type RBool byte

const (
	// NA represents R's NA
	NA RBool = 128
	// TRUE == true
	TRUE RBool = 1
	// FALSE == false
	FALSE RBool = 0
)

// IsNA returns true if the RBool is an NA
func (rb RBool) IsNA() bool {
	return rb == NA
}

// ToBoolean converts the RBool to a native go boolean.
// If the second boolean returned is false, the RBool is NA.
func (rb RBool) ToBoolean() (boolean bool, ok bool) {
	if rb == TRUE {
		boolean = true
		ok = true
	} else if rb == FALSE {
		boolean = false
		ok = true
	} else {
		boolean = false
		ok = false
	}
	return
}
