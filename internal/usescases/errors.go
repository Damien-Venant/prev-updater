package usescases

import "errors"

var (
	ErrBranchNameNotExist error = errors.New("The branch name doesn't exist in git repository")
)
