# xdir - Directory to XML converter for AI context

A memory-efficient directory to XML converter written in Go. This tool recursively processes directories and converts text files into a structured XML format, with smart file filtering and automatic exclusion of problematic directories.

## Features

- **Memory Efficient**: Uses buffered reading and reasonable memory limits
- **Smart Filtering**: Automatically excludes problematic directories (node_modules, .git, etc.)
- **File Type Support**: Processes common text files while skipping binary files
- **Progress Tracking**: Real-time progress updates with detailed statistics
- **Error Handling**: Graceful error handling with clear reporting
- **Compression Support**: Optional gzip compression for output files

## Installation

### Using Pre-built Binaries

1. Download the latest release for your platform from the releases page
2. Extract the archive:
   ```bash
   # For Linux/macOS:
   tar -xzf xdir-<os>-<arch>.tar.gz
   
   # For Windows:
   unzip xdir-windows-<arch>.zip
   ```
3. Install the binary:
   ```bash
   # Linux/macOS:
   sudo ./install.sh
   
   # Windows:
   # Move xdir.exe to a directory in your PATH
   ```

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/xdir-context.git
   cd xdir-context
   ```

2. Build the binary:
   ```bash
   go build -o xdir
   ```

3. Install to system (optional):
   ```bash
   sudo mv xdir /usr/local/bin/
   ```

## Usage

### Basic Usage

```bash
xdir [flags] sourcedir [output.xml]
```

If no output file is specified, it defaults to `output.xml`.

### Command Line Flags

- `-patterns`: File extensions to include (comma-separated)
  ```bash
  xdir -patterns=js,ts,md sourcedir output.xml
  ```

- `-glob`: Glob patterns to match files (comma-separated)
  ```bash
  xdir -glob="*.min.js,*.config.js" sourcedir output.xml
  ```

- `-compress`: Enable gzip compression for output
  ```bash
  xdir -compress sourcedir output.xml.gz
  ```

- `-max-size`: Maximum file size in bytes (default: 10MB)
  ```bash
  xdir -max-size=5242880 sourcedir output.xml
  ```

- `-unsafe`: Allow processing of normally excluded paths
  ```bash
  xdir -unsafe sourcedir output.xml
  ```

### Output Format

The tool generates XML in the following format:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<files>
  <file name="path/to/file.js" size="1234">
    <![CDATA[
    // File contents here
    ]]>
  </file>
  <!-- Additional files -->
</files>
```

## Default Behavior

### Excluded Paths
By default, the following paths are excluded:
- node_modules
- .git
- .env files
- .DS_Store
- Hidden files/directories (starting with .)

Use the `-unsafe` flag to process these paths.

### Default File Types
The tool supports these file types by default:
- Web: `.html`, `.css`, `.js`, `.jsx`, `.ts`, `.tsx`, `.json`
- Documents: `.txt`, `.md`, `.mdx`, `.xml`, `.yaml`, `.yml`
- Programming: `.go`, `.py`, `.java`, `.c`, `.cpp`, `.h`, `.hpp`, `.rs`, `.rb`, `.php`

## Performance Considerations

- Memory usage remains constant regardless of directory size
- Files larger than max-size are skipped (default 10MB)
- Binary files (images, videos, etc.) are automatically skipped
- Efficient buffered reading for large files
- Progress reporting for long-running operations

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
