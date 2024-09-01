package pipeline

import (
	"context"
	"sync/atomic"
)

// Context is a wrapper around context.Context
type Context struct {
	context.Context

	pipelineChanged *atomic.Bool
	previousChanged *atomic.Bool

	allStepCtr   *atomic.Uint32
	allUpdateCtr *atomic.Uint32

	stepCtr   atomic.Uint32
	updateCtr atomic.Uint32
}

// UpgradeContext upgrades a stdlib context to a Context
func UpgradeContext(ctx context.Context) *Context {
	res := &Context{
		Context:         ctx,
		pipelineChanged: &atomic.Bool{},
		previousChanged: &atomic.Bool{},
		allStepCtr:      &atomic.Uint32{},
		allUpdateCtr:    &atomic.Uint32{},
		stepCtr:         atomic.Uint32{},
		updateCtr:       atomic.Uint32{},
	}
	if pContext, ok := ctx.(*Context); ok {
		res.pipelineChanged = pContext.pipelineChanged
		res.previousChanged = pContext.previousChanged
		res.allStepCtr = pContext.allStepCtr
		res.allUpdateCtr = pContext.allUpdateCtr
	}
	return res
}

// BackgroundContext returns a new Context from a background context
func BackgroundContext() *Context {
	return UpgradeContext(context.Background())
}

// PipelineHasChanges returns true if any operation in the pipeline
// had changes
func (c *Context) PipelineHasChanges() bool {
	return c.pipelineChanged.Load()
}

// PreviousHasChanges returns true if the previous operation had changed
func (c *Context) PreviousHasChanges() bool {
	return c.previousChanged.Load()
}

// SetResult sets the result of an operation in the pipeline
func (c *Context) SetResult(changed bool) {
	c.stepCtr.Add(1)
	c.allStepCtr.Add(1)
	if changed {
		c.updateCtr.Add(1)
		c.allUpdateCtr.Add(1)
		c.pipelineChanged.Store(true)
	}
	c.previousChanged.Store(changed)
}

// setResultVirtual sets the result of virtual control flow operation
// without incrementing the step counter
func (c *Context) setResultVirtual(changed bool) {
	if changed {
		c.pipelineChanged.Store(true)
	}
	c.previousChanged.Store(changed)
}

// StepCtr returns the number of steps in the pipeline
func (c *Context) StepCtr(all bool) uint32 {
	if all {
		return c.allStepCtr.Load()
	}
	return c.stepCtr.Load()
}

// UpdateCtr returns the number of steps that had changes
func (c *Context) UpdateCtr(all bool) uint32 {
	if all {
		return c.allUpdateCtr.Load()
	}
	return c.updateCtr.Load()
}
