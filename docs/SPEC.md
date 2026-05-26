# Carlang Language Specification

Carlang is a tiny stack language modeled around Invoker-style spell preparation.

## Program

A program is a sequence of statements.

```txt
program = statement* EOF
```

## Statements

```txt
statement = invoke
          | cast
          | combo_call
          | combo_definition
          | block_push
```

## Chants

A chant is any non-empty sequence made only from `Q`, `W`, and `E`.

```txt
chant = ("Q" | "W" | "E")+
```

A chant followed by `R` invokes a raw spell recipe:

```carl
QQQ R
```

A chant by itself calls a memorized combo:

```carl
QQQ
```

## Blocks

Blocks defer execution.

```carl
[ QQQ R D QQQ R D WWE R D ]
```

A standalone block is pushed to the Grimoire Stack.

A block followed by a chant and `R` defines a named combo:

```carl
[ WWE R D ] QQ R
```

## Spell invocation

Only three-orb chants are valid raw spell recipes.

Recipe order is normalized by counting Q, W, and E.

```txt
QWE = QEW = WQE = WEQ = EQW = EWQ
```

All normalize to `QWE`.

## Spellbook

There are two spell slots: `D` and `F`.

Every invocation shifts the spellbook:

```txt
new spell -> D
old D     -> F
old F     -> discarded
```

## Conditional execution

`QWE R D` casts Deafening Blast.

It pops a condition from the Mana Stack and two blocks from the Grimoire Stack.

The top block is the false branch. The block beneath it is the true branch.

```carl
[ true branch ] [ false branch ] condition QWE R D
```

If `condition != 0`, the true branch executes. Otherwise, the false branch executes.

## Core spells

| Recipe | Spell | Operation |
|---|---|---|
| `QQQ` | Cold Snap | Push `1` |
| `QQW` | Ghost Walk | Drop top integer |
| `QQE` | Ice Wall | Duplicate top integer |
| `QWW` | Tornado | Swap top two integers |
| `WWW` | EMP | Subtract |
| `WWE` | Alacrity | Add |
| `WEE` | Chaos Meteor | Multiply |
| `QEE` | Forge Spirit | Duplicate top block |
| `EEE` | Sun Strike | Print top integer as an ASCII character |
| `QWE` | Deafening Blast | Conditional execution |
