# Go Playground

A simple, web-based Go code playground that allows you to write, compile, and run Go programs directly in your browser.

![Go Playground Screenshot](https://github.com/user-attachments/assets/54d9fdd5-63f7-4be2-a1f2-7b71f3951d15)

## Features

- 🚀 **Real-time Go code execution** - Write and run Go code instantly
- 🎨 **Clean, modern interface** - Syntax-highlighted code editor with a responsive design
- ⚡ **Fast compilation** - Compile and execute Go programs with minimal latency
- 🔍 **Error handling** - Clear display of compilation and runtime errors
- ⌨️ **Keyboard shortcuts** - Use `Ctrl+Enter` (or `Cmd+Enter` on Mac) to run code
- 🔒 **Safe execution** - Code runs in isolated temporary directories

## Getting Started

### Prerequisites

- Go 1.19 or later

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/AJAkimana/play-with-go.git
   cd play-with-go
   ```

2. Build the application:
   ```bash
   go build -o playground main.go
   ```

3. Run the playground:
   ```bash
   ./playground
   ```

4. Open your web browser and navigate to `http://localhost:8080`

### Using a Custom Port

You can specify a custom port by setting the `PORT` environment variable:

```bash
PORT=3000 ./playground
```

## Usage

1. **Write Go code** in the editor area
2. **Click "Run"** or use `Ctrl+Enter` to execute your code
3. **View results** in the output area below
4. **See errors** highlighted in red if your code has compilation or runtime issues

### Example Programs

Try these example programs in the playground:

**Hello World:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go Playground!")
}
```

**Working with loops and time:**
```go
package main

import (
    "fmt"
    "time"
)

func main() {
    for i := 1; i <= 5; i++ {
        fmt.Printf("Count: %d\n", i)
    }
    
    fmt.Println("Current time:", time.Now().Format("2006-01-02 15:04:05"))
}
```

## Development

### Running from Source

```bash
go run main.go
```

### Building

```bash
go build -o playground main.go
```

## Security

The playground runs user code in isolated temporary directories and uses Go's built-in compilation and execution. Each code execution:

- Runs with a 10-second timeout
- Uses temporary directories that are cleaned up after execution
- Has no access to the host filesystem beyond the temp directory

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.