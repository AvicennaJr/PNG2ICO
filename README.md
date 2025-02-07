# PNG to ICO Converter

A simple command-line tool to convert PNG images to ICO format.

## Installation

Download the latest release for your platform from the [releases page](https://github.com/AvicennaJr/png2ico/releases).

## Usage

Convert a single file:
```bash
png2ico image.png
```

### Flags

- `-o`, `--output`  : Specify the output directory for the .ico files.
- `-f`, `--force`   : Overwrite existing files.
- `-v`, `--verbose` : Enable verbose output.

### Examples

Convert a single file and specify the output directory:
```bash
png2ico image.png -o /path/to/output
```

Convert a single file and overwrite existing files:
```bash
png2ico image.png -f
```

Convert a single file with verbose output:
```bash
png2ico image.png -v
```

Combine flags:
```bash
png2ico image.png -o /path/to/output -f -v
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.