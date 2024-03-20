package msg

import "fmt"

func LogMessage(m string, keyvals ...interface{}) string {
	out := m

	if len(keyvals) > 0 && len(keyvals)%2 == 1 {
		keyvals = append(keyvals, nil)
	}

	for i := 0; i <= len(keyvals)-2; i += 2 {
		out = fmt.Sprintf("%s: %v=%v", out, keyvals[i], keyvals[i+1])
		//out = fmt.Sprintf("%s %v", out, keyvals[i])
	}

	return out
}
