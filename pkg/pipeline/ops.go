package pipeline

func ForceChange(ctx *Context) error {
	return ForceChangeIf(true)(ctx)
}

func ForceChangeIf(condition bool) FuncT {
	return func(ctx *Context) error {
		if condition {
			ctx.setResultVirtual(true)
		}
		return nil
	}
}

func ForceUpdateIfAny() FuncT {
	return func(ctx *Context) error {
		if ctx.PipelineHasChanges() {
			ctx.setResultVirtual(true)
		}
		return nil
	}
}

func DummyOperation(change bool, err error) FuncT {
	return func(ctx *Context) error {
		ctx.SetResult(change)
		return err
	}
}
