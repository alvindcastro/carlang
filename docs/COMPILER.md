# Carlang Compiler Notes

The compiler converts parsed Carlang programs into bytecode for a small stack VM.

## Pipeline

```txt
source
  -> lexer
  -> parser
  -> AST
  -> compiler
  -> bytecode
  -> VM
  -> output
```

## Compile-time expansion

Named combos are expanded at compile time.

This source:

```carl
[ WWE R D ] QQ R
QQ
```

compiles as if it were:

```carl
WWE R D
```

## Spellbook tracking

The compiler tracks `D` and `F` slots while compiling.

```carl
QQQ R
EEE R
F
D
```

The compiler sees:

```txt
D = Cold Snap
D = Sun Strike, F = Cold Snap
F -> OpPushOne
D -> OpPrintChar
```

## Initial opcodes

```txt
OpPushOne
OpDrop
OpDup
OpSwap
OpSub
OpAdd
OpMul
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

The compiler currently targets stack and ASCII-output programs. Interpreter-only features, such as general block conditionals, can be added to bytecode later with jump instructions.
