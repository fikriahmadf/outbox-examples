package failure

import "github.com/fikriahmadf/outbox-examples/shared/caller"

// AddStack will wrap error with caller function name.
func AddFuncName(err error) error {
	return Wrap(err, caller.FuncName(caller.WithSkip(1)))
}
