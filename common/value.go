package common

// Value represents a value in raw and quoted form.
type Value struct {
	Quoted string
	Raw    string
}

// String always returns the quoted value representation.
func (v Value) String() string {
	return v.Quoted
}
