package c

import "relinted/internal/formatter"

// Format takes C/C++ source code and reformats it by relocating braces and
// semicolons to the far right, creating a Python-like visual structure.
func Format(input string) string {
	return formatter.Format(input)
}
