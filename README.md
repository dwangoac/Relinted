# Relinted

Relinted reformats C/C++ and Perl source code to visually resemble Python by aligning braces (`{`, `}`) and semicolons (`;`) to the far right of the codebase. This creates a clean, Python-like visual structure while preserving the original syntax and semantics.

## Build

**Requirements:** Go 1.22+

```bash
just build          # compiles the relinted binary
just test           # runs all unit tests
just test-c         # runs C/C++ tests only
just test-perl      # runs Perl tests only
just test-rust      # runs Rust tests only
just run <input>    # runs the formatter
```

Alternatively:

```bash
go build -o relinted ./cmd/relinted/
./relinted [-l|--lang lang] <input> [output]
```

## ✨ Features

- 🔀 **Brace Relocation**: Moves leading `{` and `}` to the end of the previous line
- 📏 **Right Alignment**: Aligns trailing `;`, `{`, and `}` exactly one space past the longest line in the file
- 📐 **Indentation Preservation**: Maintains original indentation when relocating leading braces
- 🗑️ **Smart Empty Line Handling**: Removes lines that only become empty after brace relocation, but keeps originally empty lines
- 🐇 **Tab & Whitespace Normalization**: Expands tabs to 4 spaces and strips trailing whitespace for consistent visual alignment
- 📤 **Flexible I/O**: Outputs to `stdout` by default, or writes to a specified output file
- 🌐 **Multi-Language**: Supports C/C++, Perl, and Rust with auto-detection from file extension

### Arguments
| Argument        | Description                                                                            |
| --------------- | -------------------------------------------------------------------------------------- |
| `-l`, `--lang`  | *(Optional)* Language to use: `c`, `perl`, or `rust`. Overrides extension detection.   |
| `<input>`       | Path to the source file to reformat                                                    |
| `[output]`      | *(Optional)* Path to write the reformatted output; if omitted, output goes to `stdout` |

### Language Detection

Relinted auto-detects the language from the file extension:

| Extension | Language |
|-----------|----------|
| `.c`, `.h`, `.cpp`, `.cc` | C |
| `.pl`, `.pm` | Perl |
| `.rs` | Rust |

The `-l`/`--lang` flag overrides extension detection.

## Usage Examples

```bash
# Auto-detect language from extension
./relinted source.c
./relinted script.pl
./relinted game.rs

# Force a specific language
./relinted -l perl source.c
./relinted --lang c script.pl
./relinted -l rust source.c

# Write to output file
./relinted input.pl output.pl
```

## 🔍 How It Works

1. **Reads & Normalizes**: Opens the input file, expands tabs to 4 spaces if present, and strips trailing whitespace/newlines
2. **Calculates Alignment Column**: Finds the maximum line length (`max_len`) in the original file
3. **Processes Lines**:
    - Extracts leading `{` or `}` and queues them for the previous line
    - Extracts trailing `;`, `{`, `}`, or `,` (Rust match arms) and queues them for right-alignment
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

**Before (Perl)**
```perl
if ($x =~ /pattern/) {
    print "match";
}
```

**After (Python-like Visual)**
```perl
if ($x =~ /pattern/)           {
    print "match"              ;
}                              ;}
```

**Before (Rust)**
```rust
match x {
    Ok(n) => n,
    Err(_) => continue,
}
```

**After (Python-like Visual)**
```rust
match x                             {
    Ok(n) => n                      ,
    Err(_) => continue              ,}
```

## Tests

The repository contains read-only source code test example files and their read-only ground truth counterparts; running Relinted on each linted-example file should produce output identical to the matching ground truth relinted-example file:

| Linted source      | Relinted output      |
| ------------------ | -------------------- |
| linted-example-1.c | relinted-example-1.c |
| linted-example-2.c | relinted-example-2.c |
| linted-example-3.c | relinted-example-3.c |
| linted-example-4.pl| relinted-example-4.pl|
| linted-example-5.rs| relinted-example-5.rs|

## 📜 Notes

- Does not modify semantic meaning, valid syntax, or compiler behavior
- Assumes well-formed input with standard brace/semicolon placement
- Tab expansion uses 4 spaces for consistent terminal rendering
- Supports C/C++, Perl, and Rust; other languages can be added via different parsers
