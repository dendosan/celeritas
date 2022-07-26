package celeritas

import (
	"fmt"
	"regexp"
	"runtime"
	"time"
)

func (c *Celeritas) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")

	c.InfoLog.Printf(fmt.Sprintf("Load Time: %s took %s\n", name, elapsed))
}
