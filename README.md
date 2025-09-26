# Go Mdict Dictionary

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

This is a modern, high-performance MDict dictionary web application based on Go and React.

## Main Features

- **High-Performance Backend**: Developed in Go for fast dictionary query services.
- **Modern Frontend**: Built with React, Vite, TypeScript, and shadcn/ui for a beautiful and responsive interface.
- **MDict Format Support**: Compatible with common `.mdx` and `.mdd` dictionary file formats.
- **Single-File Deployment**: The backend and frontend can be packaged into a single binary for easy deployment and distribution.

## Tech Stack

- **Backend**: Go
- **Frontend**: React, TypeScript, Vite, Tailwind CSS, shadcn/ui
- **Package Manager**: pnpm

## Getting Started

### 1. Download binary

Download binary file from release.

### 2. Prepare Dictionary Files

Place your MDict dictionary files (e.g., `.mdx`, `.mdd` files) into the `dict/` folder in the project root. Each dictionary should have its own subdirectory.

For example:
```
dictionary-server
dict/
    └── MyDictionary/
        ├── MyDictionary.mdx
        ├── MyDictionary.mdd
        └── MyDictionary.css
```

### 3. Start service

## Build and Deployment

This project can be packaged into a single executable binary that includes all frontend static assets.

### 1. Build the Frontend

```bash
cd web
pnpm install
pnpm build
```

### 2. Build the Backend

The build script will automatically embed the frontend assets from the `web/dist` directory into the final Go executable.

In the project root directory, execute:

```bash
go build -o dictionary-server .
```

After the build is complete, you will get an executable file named `dictionary-server`. You can run it directly to start the entire application without needing to start the frontend service separately.

## LICENSE

[go-mdict](https://github.com/terasum/medict/blob/develop/internal/libs/go-mdict/mdict.go) GPLv3

This project is licensed under the [GPLv3](LICENSE) license. See the `LICENSE` file for details.
