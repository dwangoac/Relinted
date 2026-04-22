# Relinted

Relinted reformats C/C++ source code to visually resemble Python by aligning braces (`{`, `}`) and semicolons (`;`) to the far right of the codebase. This creates a clean, Python-like visual structure while preserving the original C/C++ syntax and semantics.

## ✨ Features

- 🔀 **Brace Relocation**: Moves leading `{` and `}` to the end of the previous line
- 📏 **Right Alignment**: Aligns trailing `;`, `{`, and `}` exactly one space past the longest line in the file
- 📐 **Indentation Preservation**: Maintains original indentation when relocating leading braces
- 🗑️ **Smart Empty Line Handling**: Removes lines that only become empty after brace relocation, but keeps originally empty lines
- 🐇 **Tab & Whitespace Normalization**: Expands tabs to 4 spaces and strips trailing whitespace for consistent visual alignment
- 📤 **Flexible I/O**: Outputs to `stdout` by default, or writes to a specified output file

### Arguments
| Argument     | Description                                                                            |
| ------------ | -------------------------------------------------------------------------------------- |
| `<input.c>`  | Path to the C/C++ source file to reformat                                              |
| `[output.c]` | *(Optional)* Path to write the reformatted output; if omitted, output goes to `stdout` |

## 🔍 How It Works

1. **Reads & Normalizes**: Opens the input file, expands tabs to 4 spaces if present, and strips trailing whitespace/newlines
2. **Calculates Alignment Column**: Finds the maximum line length (`max_len`) in the original file
3. **Processes Lines**:
   - Extracts leading `{` or `}` and queues them for the previous line
   - Extracts trailing `;`, `{`, or `}` and queues them for right-alignment
   - Preserves original indentation on lines where leading braces are moved
4. **Filters & Pads**: Removes lines that only became empty due to brace relocation, then pads each line to `max_len` and appends queued punctuation exactly one space past the maximum width
5. **Outputs**: Writes the reformatted code to `stdout` or the specified output file

## 📝 Example Transformation

**Before (Standard C)**
```c
int main() {
    printf("Hello World\n");
    return 0;
}
```

**After (Python-like Visual)**
```c
int main()                  {
    printf("Hello World\n") ;
    return 0                ;}
```

## Tests

The repository contains read-only source code test example files and their read-only ground truth counterparts; running Relinted on each linted-example file should produce output identical to the matching ground truth relinted-example file:

| Linted source      | Relinted output      |
| ------------------ | -------------------- |
| linted-example-1.c | relinted-example-1.c |
| linted-example-2.c | relinted-example-2.c |
| linted-example-3.c | relinted-example-3.c |

## 📜 Notes

- Does not modify semantic meaning, valid syntax, or compiler behavior
- Assumes well-formed C/C++ input with standard brace/semicolon placement
- Tab expansion uses 4 spaces for consistent terminal rendering
- Currently supports C; other languages can be added in the future via different parsers
