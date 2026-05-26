# Carlang Language Guide

Carlang is intentionally tiny and stack-based. To print text, programs build integer character codes from spell combos and fire them with Sun Strike.

## Core Alphabet

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

Whitespace is optional between Q/W/E sequences and control tokens:

```carl
QQQ R D
QQQRD
```

Line comments start with `//` or `#`.

## Runtime Model

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

## Spell Table

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

## Spellbook Behavior

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

## Named Combos

A block followed by a Q/W/E name and `R` memorizes that block as a combo.

```carl
[ WWE R D ] QQ R
```

This defines `QQ` as an add operation. After that:

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

That distinction lets Carlang have readable helper combos and raw Invoker-style spellcasting.

## Blocks And Conditionals

Blocks defer execution:

```carl
[ QQQ R D QQQ R D WWE R D ]
```

A standalone block is pushed to the Grimoire Stack.

`QWE R D` casts Deafening Blast. It pops a condition from the Mana Stack and two blocks from the Grimoire Stack. The top block is the false branch; the block beneath it is the true branch.

```carl
[ true branch ] [ false branch ] condition QWE R D
```

If `condition != 0`, the true branch executes. Otherwise, the false branch executes.

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

The full sample is in [`../examples/hello_world.carl`](../examples/hello_world.carl).

## Execution Paths

This repo includes two execution paths:

```txt
source -> lexer -> parser -> AST -> evaluator -> output
source -> lexer -> parser -> AST -> compiler -> bytecode -> VM -> output
```

The interpreter preserves high-level Carlang semantics directly. The compiler lowers supported invoke/cast sequences into bytecode instructions. See [compiler notes](COMPILER.md) for details.
