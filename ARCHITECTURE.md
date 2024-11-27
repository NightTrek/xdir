# Simplified High-Performance Architecture

## Core Principles
1. Stream Processing - Never load entire files into memory
2. Single Pass - Process each file exactly once
3. Minimal State - Keep only essential information in memory
4. Simple Interfaces - Clear, focused responsibilities

## Proposed Architecture

### 1. Streaming File Processor
```go
// Simple interface for processing files
type FileProcessor interface {
    ProcessFile(path string, w io.Writer) error
}

// Main processor implementation
type StreamProcessor struct {
    bufferSize int
    scanner    *bufio.Scanner
}

// Process one file at a time, streaming directly to XML
func (p *StreamProcessor) ProcessFile(path string, w io.Writer) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write file start tag
    fmt.Fprintf(w, `  <file path="%s">`, path)
    
    // Stream process the file
    scanner := bufio.NewScanner(file)
    scanner.Buffer(make([]byte, p.bufferSize), p.bufferSize)
    
    // Process in chunks, never loading entire file
    for scanner.Scan() {
        // Process and write immediately
        processChunk(scanner.Bytes(), w)
    }
    
    // Write file end tag
    fmt.Fprintf(w, "  </file>\n")
    return scanner.Err()
}
```

### 2. Context Analyzers
```go
// Simple interface for analyzing file context
type ContextAnalyzer interface {
    AnalyzeContext(r io.Reader) (Context, error)
}

// Minimal context structure
type Context struct {
    Dependencies []string
    FileType    string
    LineCount   int
}

// Language-specific analyzers
type GoAnalyzer struct{}
type JSAnalyzer struct{}

// Example of streaming analysis
func (a *GoAnalyzer) AnalyzeContext(r io.Reader) (Context, error) {
    var ctx Context
    scanner := bufio.NewScanner(r)
    
    for scanner.Scan() {
        line := scanner.Text()
        // Quick pattern matching without storing file
        if strings.HasPrefix(line, "import") {
            // Extract import immediately
            ctx.Dependencies = append(ctx.Dependencies, extractImport(line))
        }
    }
    return ctx, scanner.Err()
}
```

### 3. Main Processing Flow
```go
func ProcessDirectory(dir string, output io.Writer) error {
    // Write XML header once
    fmt.Fprintf(output, "<?xml version=\"1.0\"?>\n<files>\n")
    
    // Walk directory streaming each file
    return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
        
        // Get appropriate analyzer based on file extension
        analyzer := getAnalyzer(path)
        if analyzer == nil {
            return nil // Skip unsupported files
        }
        
        // Process file with streaming
        return processFileWithContext(path, analyzer, output)
    })
    
    // Write XML footer
    fmt.Fprintf(output, "</files>")
}
```

## Key Benefits

1. **Memory Efficiency**
   - Files are never fully loaded into memory
   - Context is processed and written immediately
   - No need for intermediate data structures

2. **Processing Speed**
   - Single pass through each file
   - No redundant file reads
   - Direct streaming to output

3. **Simplicity**
   - Clear, focused interfaces
   - Minimal state management
   - Easy to extend for new file types

4. **Reliability**
   - Predictable memory usage
   - Graceful handling of large files
   - No complex concurrency needed

## Implementation Notes

1. **Buffer Sizes**
   ```go
   const (
       defaultBufferSize = 32 * 1024  // 32KB chunks
       maxBufferSize    = 1024 * 1024 // 1MB max
   )
   ```

2. **File Processing**
   - Use `bufio.Scanner` with custom buffer sizes
   - Process line-by-line for text files
   - Stream directly to XML output

3. **Context Analysis**
   - Quick pattern matching without storing content
   - Immediate processing of dependencies
   - Skip binary files automatically

4. **Error Handling**
   - Continue on non-critical errors
   - Log issues without stopping processing
   - Clear error context for debugging

## Usage Example

```go
func main() {
    output, _ := os.Create("output.xml")
    defer output.Close()
    
    processor := NewStreamProcessor(defaultBufferSize)
    err := processor.ProcessDirectory("./project", output)
    if err != nil {
        log.Fatal(err)
    }
}
```

This architecture maintains the core functionality while significantly improving:
- Memory efficiency through streaming
- Processing speed through single-pass analysis
- Code simplicity through clear interfaces
- Reliability through predictable resource usage

The system remains focused on its primary goal: providing useful context in XML format while processing directories efficiently.
