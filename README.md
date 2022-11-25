# linefix
A small CLI tool checking files in a folder for new line

## Installation

```bash
go build -o linefix
```

## Usage

### Scan all files in a path

```bash
linefix scan <path>
```

```
Scanning 100% |█████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (7/7, 28 it/s)        

2 files affected by newline issues in . and subdirectories: 
README.md 
lel.txt
```

### Fix all files in a path

```bash
linefix fix <path>
```

```
Scanning for files with no newlines
Scanning 100% |█████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (7/7, 28 it/s)        

2 files fixed newline issues in . and subdirectories: 
README.md 
lel.txt
```
