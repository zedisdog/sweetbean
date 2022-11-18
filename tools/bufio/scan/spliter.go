package scan

import (
	"bufio"
	"bytes"
	"context"
)

// SplitByHeadAndFoot  scanner的分割函数，ctx可以控制退出
func SplitByHeadAndFoot(ctx context.Context, head, foot []byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		select {
		case <-ctx.Done():
			return 0, nil, bufio.ErrFinalToken
		default:
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			i := bytes.Index(data, head)
			if i == -1 {
				return len(data), nil, nil
			} else if i > 0 {
				return i - 1, nil, nil
			}

			if i := bytes.Index(data, foot); i >= 0 {
				// We have a full newline-terminated line.
				return i + len(foot), data[:i+len(foot)], nil
			}

			// If we're at EOF, we have a final, non-terminated line. Return it.
			if atEOF {
				return 0, nil, bufio.ErrFinalToken
			}

			// Request more data.
			return 0, nil, nil
		}
	}
}
