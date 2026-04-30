# Relinted
Relinted reformats C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, C#, PHP, and Swift source code to visually resemble Python by aligning braces (`{`, `}`) and semicolons (`;`) to the far right. This creates a clean, Python-like visual structure while preserving the original syntax. The output should still compile, but the word "should" may be doing some heavy lifting. Not for production. Output may contain syntax known to the state of California to cause code cancer, cert defects, or other deconstructive harm.

See Releases for Linux, macOS, and Windows builds (which are safe to play with and legal everywhere questionable taste has not been outlawed).

## Features
- **Relocates braces**: Moves leading `{` and `}` to the end of the previous line
- **Right alignment**: Aligns trailing `;`, `{`, and `}` one space past the longest line in the file
- **Preserves indentation**: Maintains original indentation when relocating leading braces
- **Removes new empty lines**: Removes lines that only become empty after brace relocation, but keeps originally empty lines
- **Normalizes tabs**: Expands tabs to 4 spaces and strips trailing whitespace for consistent visual alignment
- **Outputs wherever**: Outputs to `stdout` by default, or writes to a specified output file
- **Language autodetection**: Supports C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, C#, PHP, and Swift with auto-detection from file extension

### Arguments
| Argument   | Description                                                                                                                          |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| `-l`       | *(Optional)* Language to use: `c`, `perl`, `rust`, `js`, `ts`, `go`, `java`, `cs`, `php`, or `swift`. Overrides extension detection. |
| `<input>`  | Path to the source file to reformat                                                                                                  |
| `[output]` | *(Optional)* Path to write the reformatted output; if omitted, output goes to `stdout`                                               |

## Usage Examples
```bash
# Auto-detect language from extension
./relinted source.c

# Force a specific language
./relinted -l perl perlsourcefile

# Write to output file
./relinted input.rs output.rs
```

## Example Transformation
**Before**
```c
int main() {
    printf("Hello World\n");
    return 0;
}
```

**Relinted**
```c
int main()                  {
    printf("Hello World\n") ;
    return 0                ;}
```

### Language Detection
Relinted auto-detects the language via file extension:

| Extension                 | Language   |
| ------------------------- | ---------- |
| `.c`, `.h`, `.cpp`, `.cc` | C          |
| `.pl`, `.pm`              | Perl       |
| `.rs`                     | Rust       |
| `.js`                     | JavaScript |
| `.ts`                     | TypeScript |
| `.go`                     | Go         |
| `.java`                   | Java       |
| `.cs`                     | C#         |
| `.php`                    | PHP        |
| `.swift`                  | Swift      |

The `-l` flag can be used to override extension detection.

## Manual Build
Building manually requires Go 1.22 or higher.

```bash
just build                           # Build with Just
go build -o relinted ./cmd/relinted/ # Or build directly
```

## Tests
The repository contains example read-only source code test files and their read-only ground truth counterparts; running Relinted on each linted-example file should produce output identical to the matching ground truth relinted-example file:

| Linted source                    | Relinted output                    |
| -------------------------------- | ---------------------------------- |
| examples/linted-example-1.c      | examples/relinted-example-1.c      |
| examples/linted-example-2.c      | examples/relinted-example-2.c      |
| examples/linted-example-3.c      | examples/relinted-example-3.c      |
| examples/linted-example-4.pl     | examples/relinted-example-4.pl     |
| examples/linted-example-5.rs     | examples/relinted-example-5.rs     |
| examples/linted-example-6.js     | examples/relinted-example-6.js     |
| examples/linted-example-7.go     | examples/relinted-example-7.go     |
| examples/linted-example-8.java   | examples/relinted-example-8.java   |
| examples/linted-example-9.ts     | examples/relinted-example-9.ts     |
| examples/linted-example-10.cs    | examples/relinted-example-10.cs    |
| examples/linted-example-11.php   | examples/relinted-example-11.php   |
| examples/linted-example-12.swift | examples/relinted-example-12.swift |

Tests can be executed with `just test`:
just test           # runs all unit tests
just test-c         # runs C/C++ tests only
just test-perl      # runs Perl tests only
just test-rust      # runs Rust tests only
just test-js        # runs JavaScript tests only
just test-go        # runs Go tests only
just test-java      # runs Java tests only
just test-php       # runs PHP tests only
just test-swift     # runs Swift tests only
just run <input>    # runs the formatter

## Notes
- Does not modify semantic meaning, valid syntax, or compiler behavior
- Does not guarantee resulting code will compile but aims to not make breaking changes
- Assumes well-formed input with standard brace/semicolon placement
- Tab expansion attempts to use 4 spaces for consistent terminal rendering
- Supports C/C++, Perl, Rust, JavaScript, TypeScript, Go, Java, C#, PHP, and Swift; other languages can be added via different parsers
- Relinted may exhibit clank jank and makes no guarantee about the sanity of the code, the supposed coder, or the user
