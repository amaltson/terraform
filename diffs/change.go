package diffs

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

// Change represents a change to a complex object, such as a resource.
//
// Rather than instantiating this struct directly, prefer instead to use the
// NewCreate, NewRead, NewUpdate, NewDelete, and NewReplace methods, which
// ensure that expected invariants are met.
type Change struct {
	// Action identifies the kind of action represented by this change.
	Action Action

	// Type defines the type of value this change applies to. The old
	// and new values must have a type that conforms to this type, as
	// defined by cty.
	Type cty.Type

	// Old and New are the expected values before and after the action
	// respectively. The usage of these varies by Action:
	//
	// For Create, Old must be a null value of the change type and New
	// is the value that will be created.
	//
	// For Read, Old must be a null value of the change type and New
	// describes the object that will be read. A New value for Read will
	// generally consist primarily of unknown values, thus indicating
	// the expected structure of the result even though the values are not
	// yet known.
	//
	// For Update and Replace, Old describes the expected existing value to
	// update and New is its replacement. For Replace, ForceNew is also
	// populated.
	//
	// For Delete, Old is the expected value to destroy and New must be a null
	// value of the change type.
	Old, New cty.Value

	// ForcedReplace is populated for changes with action Replace to indicate
	// which paths within the value prompted the change to be a Replace rather
	// than an Update.
	//
	// The set may be populated with paths from the old or new value, or both,
	// depending on the nature of the change being described. A diff renderer
	// should annotate its rendering of a particular path with an
	// indication that it prompted replacement if that path is present in
	// this set.
	//
	// ForceNew is nil for all other actions and will panic if accessed.
	ForcedReplace PathSet
}

func NewCreate(ty cty.Type, v cty.Value) *Change {
	if errs := ty.TestConformance(v.Type()); len(errs) != 0 {
		panic(fmt.Errorf("value does not conform to type: %s", errs[0]))
	}

	return &Change{
		Action: Create,
		Type:   ty,
		Old:    cty.NullVal(ty),
		New:    v,
	}
}

func NewRead(ty cty.Type, v cty.Value) *Change {
	if errs := ty.TestConformance(v.Type()); len(errs) != 0 {
		panic(fmt.Errorf("value does not conform to type: %s", errs[0]))
	}

	return &Change{
		Action: Read,
		Type:   ty,
		Old:    cty.NullVal(ty),
		New:    v,
	}
}

func NewUpdate(ty cty.Type, old, new cty.Value) *Change {
	if errs := ty.TestConformance(old.Type()); len(errs) != 0 {
		panic(fmt.Errorf("old value does not conform to type: %s", errs[0]))
	}
	if errs := ty.TestConformance(new.Type()); len(errs) != 0 {
		panic(fmt.Errorf("new value does not conform to type: %s", errs[0]))
	}

	return &Change{
		Action: Update,
		Type:   ty,
		Old:    old,
		New:    new,
	}
}

func NewReplace(ty cty.Type, old, new cty.Value, forcedReplace PathSet) *Change {
	if errs := ty.TestConformance(old.Type()); len(errs) != 0 {
		panic(fmt.Errorf("old value does not conform to type: %s", errs[0]))
	}
	if errs := ty.TestConformance(new.Type()); len(errs) != 0 {
		panic(fmt.Errorf("new value does not conform to type: %s", errs[0]))
	}

	return &Change{
		Action: Replace,
		Type:   ty,
		Old:    old,
		New:    new,
	}
}

func NewDelete(ty cty.Type, v cty.Value, forcedReplace PathSet) *Change {
	if errs := ty.TestConformance(v.Type()); len(errs) != 0 {
		panic(fmt.Errorf("value does not conform to type: %s", errs[0]))
	}

	return &Change{
		Action: Replace,
		Type:   ty,
		Old:    v,
		New:    cty.NullVal(ty),
	}
}
