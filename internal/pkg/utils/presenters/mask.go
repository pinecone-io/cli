package presenters

func MaskHeadTail(s string, head, tail int) string {
	if head < 0 {
		head = 0
	}
	if tail < 0 {
		tail = 0
	}

	runes := []rune(s)
	length := len(runes)
	if length == 0 || head+tail >= length {
		return s
	}

	start := string(runes[:min(head, length)])
	end := ""
	if tail > 0 {
		end = string(runes[length-min(tail, length):])
	}

	return start + "***" + end
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
