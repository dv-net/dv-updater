package retry

import "fmt"

type Policy int

const (
	PolicyLinear Policy = iota
	PolicyBackoff
	PolicyInfinite
)

func (r Policy) Validate() error {
	switch r {
	case PolicyLinear, PolicyBackoff, PolicyInfinite:
		return nil
	default:
		return fmt.Errorf("invalid retry policy")
	}
}
