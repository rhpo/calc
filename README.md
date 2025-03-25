
# Calc - Command Line Calculator

## Description

`Calc` is a command-line calculator written in Go that supports basic arithmetic and advanced mathematical operations. It works on both Linux and Windows systems and provides an interactive way to perform calculations directly from the terminal.

## Features

- Basic arithmetic operations (addition, subtraction, multiplication, division).
- Advanced functions (trigonometric, logarithmic, and more).
- Works on Linux and Windows platforms.
- Interactive terminal-based interface.

## Installation

### Linux / Windows

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/calc.git
   ```

2. Navigate to the project folder:

   ```bash
   cd calc
   ```

3. Build the Go application:

   ```bash
   go build
   ```

4. Run the calculator:

   ```bash
   ./calc   # For Linux
   calc.exe # For Windows
   ```

## Usage

Once the program is running, you can start entering mathematical expressions.

Example usage:

```bash
$ ./calc
Calc v1.0, type 'exit' to leave.
>> 1+1
2
>> 1 + 3.032
4.032
>> cos(1 + pi / e)
-0.5521418740547764
>> exit
```

Type `exit` to quit the program.

## License

MIT License. See [LICENSE](LICENSE) for more details.
