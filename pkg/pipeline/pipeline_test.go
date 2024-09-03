package pipeline

import (
	"io"
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/internal/test"
)

type pipelineTestCase struct {
	Name                    string
	Funcs                   []FuncT
	ExpectedStepCtr         uint32
	ExpectedAllStepCtr      uint32
	ExpectedUpdateCtr       uint32
	ExpectedAllUpdateCtr    uint32
	ExpectedPipelineChanged bool
	ExpectedPreviousChanged bool
	ExpectError             bool
}

func TestPipelines(t *testing.T) {
	cases := []pipelineTestCase{
		{
			Name:                    "No Operations",
			Funcs:                   []FuncT{},
			ExpectedStepCtr:         0,
			ExpectedAllStepCtr:      0,
			ExpectedUpdateCtr:       0,
			ExpectedAllUpdateCtr:    0,
			ExpectedPipelineChanged: false,
			ExpectedPreviousChanged: false,
			ExpectError:             false,
		},
		{
			Name: "Single Operation",
			Funcs: []FuncT{
				DummyOperation(true, nil),
			},
			ExpectedStepCtr:         1,
			ExpectedAllStepCtr:      1,
			ExpectedUpdateCtr:       1,
			ExpectedAllUpdateCtr:    1,
			ExpectedPipelineChanged: true,
			ExpectedPreviousChanged: true,
			ExpectError:             false,
		},
		{
			Name: "Operation error",
			Funcs: []FuncT{
				DummyOperation(true, nil),
				DummyOperation(true, nil),
				DummyOperation(false, nil),
				DummyOperation(false, io.EOF),
				DummyOperation(true, nil),
			},
			ExpectedStepCtr:         4,
			ExpectedAllStepCtr:      4,
			ExpectedUpdateCtr:       2,
			ExpectedAllUpdateCtr:    2,
			ExpectedPipelineChanged: true,
			ExpectedPreviousChanged: false,
			ExpectError:             true,
		},
		{
			Name: "Nested pipeline",
			Funcs: []FuncT{
				DummyOperation(false, nil),
				DummyOperation(false, nil),
				NewPipelineP([]FuncT{
					DummyOperation(false, nil),
					DummyOperation(true, nil),
				}),
				DummyOperation(false, nil),
			},
			ExpectedStepCtr:         4,
			ExpectedAllStepCtr:      6,
			ExpectedUpdateCtr:       1,
			ExpectedAllUpdateCtr:    2,
			ExpectedPipelineChanged: true,
			ExpectedPreviousChanged: false,
			ExpectError:             false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx, r := test.DefaultPreamble(t, time.Second*10)
			p := NewPipeline(ctx, tc.Funcs)
			err := p.Run()
			if tc.ExpectError {
				r.Error(err)
			} else {
				r.NoError(err)
			}
			r.Equal(tc.ExpectedStepCtr, p.StepCtr(false), "step ctr mistmatch")
			r.Equal(tc.ExpectedAllStepCtr, p.StepCtr(true), "all step ctr mistmatch")
			r.Equal(tc.ExpectedUpdateCtr, p.UpdateCtr(false), "update ctr mistmatch")
			r.Equal(tc.ExpectedAllUpdateCtr, p.UpdateCtr(true), "all update ctr mistmatch")
			r.Equal(tc.ExpectedPipelineChanged, p.PipelineHasChanges(), "pipeline changed mismatch")
			r.Equal(tc.ExpectedPreviousChanged, p.PreviousHasChanges(), "previous changed mismatch")
		})
	}
}
