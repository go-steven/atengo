package show_source_line

import (
	"bytes"
	"fmt"
	"strings"
)

func GetSourceWithLineNum(script string) string {
	buff := bytes.NewBufferString("")
	lines := strings.Split(script, "\n")
	for idx, v := range lines {
		buff.WriteString(fmt.Sprintf("%d:		%s\n", idx+1, v))
	}
	return buff.String()
}
