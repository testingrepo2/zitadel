// Code generated by "stringer -type Type"; DO NOT EDIT.

package milestone

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[unknown-0]
	_ = x[InstanceCreated-1]
	_ = x[AuthenticationSucceededOnInstance-2]
	_ = x[ProjectCreated-3]
	_ = x[ApplicationCreated-4]
	_ = x[AuthenticationSucceededOnApplication-5]
	_ = x[InstanceDeleted-6]
	_ = x[typesCount-7]
}

const _Type_name = "unknownInstanceCreatedAuthenticationSucceededOnInstanceProjectCreatedApplicationCreatedAuthenticationSucceededOnApplicationInstanceDeletedtypesCount"

var _Type_index = [...]uint8{0, 7, 22, 55, 69, 87, 123, 138, 148}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
