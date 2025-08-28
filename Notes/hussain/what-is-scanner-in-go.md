
# What `bufio.NewScanner` is

`bufio.NewScanner(r io.Reader)` returns a `*bufio.Scanner` that reads text **token by token** from any `io.Reader` (files, network connections, stdin…).  
By default, each **token = one line**.

# Core API (90% of usage)

```go
s := bufio.NewScanner(r)   // r: any io.Reader (e.g., *os.File)
for s.Scan() {             // advances to next token
    line := s.Text()       // token as string (copied)
    // use line
}
if err := s.Err(); err != nil {
    // handle read/scan error (not EOF)
}
```

- `Scan()` → moves to next token; returns `false` at EOF or on error.
    
- `Text()` → current token as `string` (safe to keep).
    
- `Bytes()` → current token as `[]byte` (valid only until next `Scan()`).
    
- `Err()` → returns non-EOF error, if any.
    

# Tokenization (how it splits input)

- **Default**: `bufio.ScanLines` (line-by-line, trims trailing `\n`).
    
- Built-ins:
    
    - `bufio.ScanWords` (word-by-word; whitespace-delimited)
        
    - `bufio.ScanRunes` (rune-by-rune; Unicode-aware)
        
- Custom:
    

```go
s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    // Return how many bytes to advance, the token slice, and error if any.
})
```

Example: switch to words

```go
s := bufio.NewScanner(r)
s.Split(bufio.ScanWords)
for s.Scan() {
    word := s.Text()
}
```

# The 64KB gotcha (max token size)

`Scanner` limits a single token to **64 KiB by default**. If a line (or token) is longer, `Scan()` will return an error like “token too long”.

Fix it with a bigger buffer **and** a larger max token size:

```go
s := bufio.NewScanner(r)
buf := make([]byte, 0, 1024*1024)          // initial buffer (1 MiB)
s.Buffer(buf, 10*1024*1024)                // max token size (10 MiB)
```

> If you might have large lines (logs, JSONL, CSV), do this.

# When to use `Scanner`

Use `Scanner` when you want **simple, clean iteration** over:

- Lines (config files, logs, text files)
    
- Words (tokenization for quick parsing)
    
- Runes (character-level processing)
    
- Custom tokens (e.g., records separated by `"---"`)
    

It’s perfect for **human-readable text** where tokens aren’t huge and you like a tidy loop.

# When **not** to use `Scanner` (choose `bufio.Reader` instead)

- **Very long tokens/lines** (multi-MB CSV rows, minified JSON): `Scanner` can be made to work, but `Reader` is often better.
    
- **Binary protocols** or mixed-size records where you need exact byte control.
    
- You need methods like `ReadBytes`, `ReadString`, `Peek`, `ReadSlice`, or buffered `Read` performance tricks.
    

`bufio.Reader` snippet for long lines:

```go
r := bufio.NewReader(f)
for {
    line, err := r.ReadString('\n')  // can read arbitrarily long lines
    if err == io.EOF {
        if len(line) > 0 { /* handle last line */ }
        break
    }
    if err != nil { return err }
    // use line
}
```

# Common examples

**1) Read lines from a file**

```go
f, err := os.Open("input.txt")
if err != nil { log.Fatal(err) }
defer f.Close()

s := bufio.NewScanner(f)
for s.Scan() {
    fmt.Println(s.Text())
}
if err := s.Err(); err != nil {
    log.Fatal(err)
}
```

**2) Read from stdin**

```go
s := bufio.NewScanner(os.Stdin)
for s.Scan() {
    fmt.Println("You typed:", s.Text())
}
if err := s.Err(); err != nil {
    log.Fatal(err)
}
```

**3) Words instead of lines**

```go
s := bufio.NewScanner(strings.NewReader("hello world from Go"))
s.Split(bufio.ScanWords)
for s.Scan() {
    fmt.Println(s.Text())
}
```

**4) Custom delimiter (e.g., `---` between records)**

```go
delim := []byte("---")
splitFunc := func(data []byte, atEOF bool) (int, []byte, error) {
    if i := bytes.Index(data, delim); i >= 0 {
        return i + len(delim), data[:i], nil
    }
    if atEOF && len(data) > 0 {
        return len(data), data, nil
    }
    return 0, nil, nil
}

s := bufio.NewScanner(r)
s.Buffer(make([]byte, 0, 64*1024), 10*1024*1024) // raise limit if needed
s.Split(splitFunc)
for s.Scan() {
    record := s.Text()
    // process record
}
if err := s.Err(); err != nil {
    log.Fatal(err)
}
```

**5) Rune-by-rune (Unicode safe)**

```go
s := bufio.NewScanner(r)
s.Split(bufio.ScanRunes)
for s.Scan() {
    fmt.Printf("rune: %q\n", s.Text())
}
```

# Performance notes

- `Scanner` is **simple and safe**, not the absolute fastest for heavy I/O.
    
- `Text()` returns a **string copy**; `Bytes()` avoids copying but is only valid until next `Scan()`.
    
- For very high throughput or complex parsing, `bufio.Reader` (or `Read` into your own buffer) can be more efficient.
    

# Error & EOF behavior

- Loop ends when `Scan()` is `false`. Then:
    
    - `s.Err()` is `nil` → that’s **EOF** (normal).
        
    - `s.Err()` is non-nil → a real **error** occurred (I/O error, token too long, your split func error).
        
- Final line without trailing newline is still returned as a token. Good for reading files that don’t end with `\n`.
    

# Concurrency

- A `Scanner` isn’t goroutine-safe. Use it from a single goroutine per scanner/reader.
    

# TL;DR

- Use `Scanner` for **clean token iteration** (lines/words/runes/custom), **small-to-moderate tokens**, and **simple code**.
    
- Remember the **64 KiB** token limit; bump it with `Scanner.Buffer` when needed.
    
- For **huge lines**, **binary**, or **fine-grained control**, go with `bufio.Reader`.
    

