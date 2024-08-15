package types

import "strings"

type MultiError struct {
	Errors []error
}

func (m *MultiError) Add(err error) {
	m.Errors = append(m.Errors, err)
}

func (m *MultiError) HasError() bool {
	return len(m.Errors) > 0
}

func (m *MultiError) Error() string {
	var errMsg strings.Builder
	for _, err := range m.Errors {
		errMsg.WriteString(err.Error() + "; ")
	}
	return errMsg.String()
}
