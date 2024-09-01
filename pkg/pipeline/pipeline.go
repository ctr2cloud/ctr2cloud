package pipeline

import "context"

type FuncT func(ctx *Context) error

type Pipeline struct {
	*Context
	operations []FuncT
}

// NewPipeline creates a new pipeline
func NewPipeline(ctx context.Context, operations []FuncT) *Pipeline {
	return &Pipeline{
		Context:    UpgradeContext(ctx),
		operations: operations,
	}
}

// NewPipelineP creates a nested pipeline
func NewPipelineP(operations []FuncT) FuncT {
	return func(ctx *Context) error {
		p := NewPipeline(ctx, operations)
		err := p.Run()
		ctx.SetResult(p.UpdateCtr(false) > 0)
		return err
	}
}

func (p *Pipeline) Run() error {
	for _, operation := range p.operations {
		err := operation(p.Context)
		if err != nil {
			return err
		}
	}
	return nil
}
