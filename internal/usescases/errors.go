package usescases

import "errors"

var (
	ErrBranchNameNotExist error = errors.New("the branch name doesn't exist in git repository")
)
