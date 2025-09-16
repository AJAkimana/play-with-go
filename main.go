package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const defaultCode = `package main

import "fmt"

func main() {
	fmt.Println("Hello, Go Playground!")
}`

type PlaygroundResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/run", handleRun)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting Go Playground on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Go Playground</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f8f9fa;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            color: #007d9c;
            text-align: center;
            margin-bottom: 30px;
        }
        .playground {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 20px;
        }
        .editor-container {
            margin-bottom: 20px;
        }
        textarea {
            width: 100%;
            height: 300px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 10px;
            resize: vertical;
            outline: none;
            box-sizing: border-box;
        }
        textarea:focus {
            border-color: #007d9c;
            box-shadow: 0 0 0 2px rgba(0,125,156,0.2);
        }
        .controls {
            margin-bottom: 20px;
        }
        .run-button {
            background-color: #007d9c;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 500;
        }
        .run-button:hover {
            background-color: #006080;
        }
        .run-button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        .output {
            background-color: #f8f9fa;
            border: 1px solid #ddd;
            border-radius: 4px;
            padding: 15px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            white-space: pre-wrap;
            min-height: 150px;
            max-height: 400px;
            overflow-y: auto;
        }
        .error {
            color: #dc3545;
        }
        .loading {
            color: #6c757d;
            font-style: italic;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Go Playground</h1>
        <div class="playground">
            <div class="editor-container">
                <textarea id="code" placeholder="Write your Go code here...">{{.Code}}</textarea>
            </div>
            <div class="controls">
                <button class="run-button" onclick="runCode()">Run</button>
            </div>
            <div class="output" id="output">Click "Run" to execute your code.</div>
        </div>
    </div>

    <script>
        async function runCode() {
            const code = document.getElementById('code').value;
            const outputEl = document.getElementById('output');
            const runButton = document.querySelector('.run-button');
            
            runButton.disabled = true;
            runButton.textContent = 'Running...';
            outputEl.textContent = 'Running your code...';
            outputEl.className = 'output loading';
            
            try {
                const response = await fetch('/run', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: 'code=' + encodeURIComponent(code)
                });
                
                const result = await response.json();
                
                if (result.error) {
                    outputEl.textContent = result.error;
                    outputEl.className = 'output error';
                } else {
                    outputEl.textContent = result.output || 'Program finished with no output.';
                    outputEl.className = 'output';
                }
            } catch (error) {
                outputEl.textContent = 'Error: ' + error.message;
                outputEl.className = 'output error';
            }
            
            runButton.disabled = false;
            runButton.textContent = 'Run';
        }

        // Allow Ctrl+Enter to run code
        document.getElementById('code').addEventListener('keydown', function(e) {
            if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
                e.preventDefault();
                runCode();
            }
        });
    </script>
</body>
</html>`

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Code string
	}{
		Code: defaultCode,
	}

	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, data)
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "No code provided", http.StatusBadRequest)
		return
	}

	result, err := runGoCode(code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"error": %q}`, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"output": %q}`, result)
}

func runGoCode(code string) (string, error) {
	// Create a temporary directory for the Go code
	tempDir, err := os.MkdirTemp("", "go-playground-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write the Go code to a file
	codePath := filepath.Join(tempDir, "main.go")
	err = os.WriteFile(codePath, []byte(code), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write code file: %v", err)
	}

	// Initialize go module in temp directory
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "mod", "init", "playground")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to initialize module: %v", err)
	}

	// Run the Go code
	cmd = exec.CommandContext(ctx, "go", "run", "main.go")
	cmd.Dir = tempDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		// Check if it's a compilation error
		if stderr.Len() > 0 {
			return "", fmt.Errorf("compilation error:\n%s", stderr.String())
		}
		return "", fmt.Errorf("runtime error: %v", err)
	}

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n--- stderr ---\n"
		}
		output += stderr.String()
	}

	// Trim excessive whitespace but preserve intentional formatting
	output = strings.TrimSpace(output)
	
	return output, nil
}