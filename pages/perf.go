package pages

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type ServerTiming struct {
	Name        string
	Duration    time.Duration
	Description string
}

func (st ServerTiming) String() string {
	builder := strings.Builder{}
	builder.WriteString(st.Name)

	if len(st.Description) != 0 {
		builder.WriteString(fmt.Sprintf(";desc=\"%s\"", st.Description))
	}

	if st.Duration != 0 {
		builder.WriteString(fmt.Sprintf(";dur=%d", st.Duration.Milliseconds()))
	}
	return builder.String()
}

type Timing struct {
	name  string
	desc  string
	start time.Time
}

func (t Timing) Done(ctx context.Context) {
	st := &ServerTiming{Name: t.name, Description: t.desc, Duration: time.Since(t.start)}
	pctx, ok := ctx.(*PageContext)
	if ok {
		pctx.timings = append(pctx.timings, st)
	} else {
		lctx, ok := ctx.(*LayoutContext)
		if ok {
			lctx.timings = append(lctx.timings, st)
		}
	}
}

func NewTiming(name, desc string) *Timing {
	return &Timing{start: time.Now(), name: name, desc: desc}
}
