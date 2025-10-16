# go-file-word-checker
A simple, concurrent command-line tool written in Go to find a specific word within files.

## Usage
Run the program from your terminal with the word you want to find, followed by one or more file paths or patterns.

```Bash
./go-file-word-checker <word-to-find> <file_or_path_pattern1> [file_or_path_pattern2] ...
```

## Example

Search for "brother" in all txt files:
```Bash
./go-file-word-checker 'brother' '*.txt'
```