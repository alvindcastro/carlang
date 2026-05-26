# Carlang

Carlang is an unofficial Invoker-inspired esoteric programming language.

Programs are written by combining **Quas**, **Wex**, and **Exort** into three-orb recipes, invoking those recipes with `R`, and casting from two spell slots with `D` and `F`.

The language is intentionally tiny, stack-based, and cursed in the correct way: to print text, you build ASCII values out of spell combos and fire them with Sun Strike.

> This is a fan project and is not affiliated with Valve, Dota 2, or any official Dota property.

## Quick start

```bash
go test ./...
go run ./cmd/carl run examples/hello_world.carl
```

Expected output:

```txt
Hello, World!
```

Build the CLI:

```bash
go build -o bin/carl ./cmd/carl
bin/carl run examples/hello_world.carl
```

Compile to bytecode and run on the VM:

```bash
bin/carl compile examples/hello_world.carl -o dist/hello_world.cbc
bin/carl vm dist/hello_world.cbc
```

Disassemble a source file:

```bash
bin/carl dis examples/hello_world.carl
```

Start the REPL:

```bash
bin/carl repl
```

## Core syntax

The core alphabet is:

```txt
Q W E R D F [ ]
```

| Token | Meaning |
|---|---|
| `Q` | Quas |
| `W` | Wex |
| `E` | Exort |
| `R` | Invoke a spell, or memorize a combo after a block |
| `D` | Cast newest invoked spell |
| `F` | Cast older invoked spell |
| `[` `]` | Define a deferred code block |

Whitespace is optional between Q/W/E sequences and the control tokens, so these are equivalent:

```carl
QQQ R D
QQQRD
```

Line comments start with `//`.

## Runtime model

Carlang has four main runtime pieces:

```txt
Mana Stack      integer stack
Grimoire Stack  stack of deferred code blocks
Spellbook       two invoked spell slots: D and F
Tome            dictionary of named combos
```

Example:

```carl
QQQ R D
```

Meaning:

```txt
QQQ   prepare Cold Snap recipe
R     invoke it into the D slot
D     cast the D slot
```

Result:

```txt
Mana Stack: [1]
```

## Spell table

Recipes are normalized by counting `Q`, `W`, and `E`. Order does not matter for spell invocation:

```carl
QWE R D
WEQ R D
EQW R D
```

All three invoke Deafening Blast.

| Recipe | Spell | Carlang operation | Stack effect |
|---|---|---|---|
| `QQQ` | Cold Snap | Push one | `[] -> [1]` |
| `QQW` | Ghost Walk | Drop | `[a] -> []` |
| `QQE` | Ice Wall | Duplicate | `[a] -> [a, a]` |
| `QWW` | Tornado | Swap | `[a, b] -> [b, a]` |
| `WWW` | EMP | Subtract | `[a, b] -> [a - b]` |
| `WWE` | Alacrity | Add | `[a, b] -> [a + b]` |
| `WEE` | Chaos Meteor | Multiply | `[a, b] -> [a * b]` |
| `QEE` | Forge Spirit | Duplicate top block | `Grimoire: [B] -> [B, B]` |
| `EEE` | Sun Strike | Print ASCII character | `[a] -> []` and writes `char(a)` |
| `QWE` | Deafening Blast | Conditional block execution | pops condition and two blocks |

## Spellbook behavior

Carlang has two spell slots.

```carl
QQQ R   // D = Cold Snap
EEE R   // D = Sun Strike, F = Cold Snap
F       // cast Cold Snap
D       // cast Sun Strike
```

When a new spell is invoked:

```txt
new spell -> D
old D     -> F
old F     -> forgotten
```

## Named combos

A block followed by a Q/W/E name and `R` memorizes that block as a combo.

```carl
[ WWE R D ] QQ R
```

This defines `QQ` as an add operation.

After that:

```carl
QQQ R D   // push 1
QQQ R D   // push 1
QQ        // call combo QQ, which adds
```

The stack now contains `2`.

A name by itself calls a combo:

```carl
QQ
```

A name followed by `R` invokes a raw spell recipe:

```carl
QQQ R
```

That distinction is what lets Carlang have both readable helper combos and raw Invoker-style spellcasting.

## Hello World

Carlang has no string literals in the core language. `Hello, World!` is produced by building ASCII values from stack operations and printing each character with Sun Strike.

```carl
// Core helper combos
[ WWE R D ] QQ R
[ EEE R D ] EQE R
[ QQE R D QQ ] EEQ R

// Powers of two
[ QQQ R D QQQ R D QQ ] QQQ R
[ QQQ EEQ ] QQE R
[ QQE EEQ ] QQW R
[ QQW EEQ ] QEQ R
[ QEQ EEQ ] QEE R
[ QEE EEQ ] QEW R

// H = 72 = 64 + 8
QEW QQW QQ EQE
```

The full sample is in [`examples/hello_world.carl`](examples/hello_world.carl).

## Interpreter and compiler

This repo includes two execution paths:

```txt
source -> lexer -> parser -> AST -> evaluator -> output
source -> lexer -> parser -> AST -> compiler -> bytecode -> VM -> output
```

The interpreter preserves the high-level Carlang semantics directly.

The compiler lowers raw invoke/cast sequences into bytecode instructions such as:

```txt
OpPushOne
OpDup
OpAdd
OpPrintChar
```

For example:

```carl
QQQ R D
QQE R D
WWE R D
EEE R D
```

compiles to:

```txt
0000 OpPushOne
0001 OpDup
0002 OpAdd
0003 OpPrintChar
```

## Repository layout

```txt
cmd/carl/      CLI entrypoint
token/         token definitions
lexer/         lexer
ast/           AST nodes
parser/        parser
spell/         spell names and recipe normalization
evaluator/     tree-walking evaluator and interpreter VM
code/          bytecode instruction format
compiler/      source-to-bytecode compiler
vm/            bytecode virtual machine
repl/          interactive mode
examples/      sample programs
docs/          language notes
```

## Current status

Implemented:

- Lexing and parsing
- Named combos
- Interpreter execution
- Spellbook slots `D` and `F`
- Core arithmetic/stack spells
- ASCII output through Sun Strike
- Bytecode compiler
- Bytecode VM
- CLI commands: `run`, `compile`, `vm`, `dis`, `repl`

Partially implemented:

- Block operations and conditionals are supported by the interpreter, but the compiler currently focuses on stack/ASCII programs.

## License

MIT
