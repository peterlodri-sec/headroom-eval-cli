# 🐋 headroom-eval-cli

Interactive ASCII TUI for the [Headroom Eval Space](https://huggingface.co/spaces/PeetPedro/headroom-eval).
Live status, animated dodecahedron, steganography.

## Install

```bash
go install github.com/peterlodri-sec/headroom-eval-cli@latest
headroom-eval
```

## Features

- **Live Space status** — polls HF API every 3s
- **Animated ASCII dodecahedron** — two-frame animation
- **Steganography** — press `s` to reveal hidden message
- **Genesis seal** — cryptographic hash of the binary
- **Quick links** — open Space, paper, LoopKit from TUI

## Keys

| Key | Action |
|---|---|
| `r` | Refresh status |
| `o` | Open Space in browser |
| `p` | Open paper (kompress.vaked.dev) |
| `g` | Open LoopKit repo |
| `s` | Toggle steganography |
| `q` | Quit |

## CLI Flags

```bash
headroom-eval --seal      # print genesis hash
headroom-eval --version   # print version + build time
```

## Build from source

```bash
git clone https://github.com/peterlodri-sec/headroom-eval-cli
cd headroom-eval-cli
make build-upx   # → 2.6MB binary
```

## Part of

- [Headroom Eval Space](https://huggingface.co/spaces/PeetPedro/headroom-eval)
- [LoopKit](https://github.com/peterlodri-sec/loopkit)
- [kompress paper (ICLR 2027)](https://kompress.vaked.dev)
