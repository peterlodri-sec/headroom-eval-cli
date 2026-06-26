# 🐋 headroom-eval-cli

Interactive pure-ASCII TUI for the [Headroom Eval Space](https://huggingface.co/spaces/PeetPedro/headroom-eval).
Live status, animated ASCII dodecahedron, steganography, genesis seal.

## Install

```bash
go install github.com/peterlodri-sec/headroom-eval-cli@latest
headroom-eval
```

## Keys

| Key | Action |
|---|---|
| `r` | Refresh Space status |
| `o` | Open Space in browser |
| `p` | Open paper (kompress.vaked.dev) |
| `g` | Open LoopKit |
| `d` | Sponsor/donate |
| `s` | Toggle steganography |
| `q` | Quit |

## Flags

```bash
headroom-eval --seal     # genesis hash
headroom-eval --version  # version + build time
```

## Links

- [Headroom Eval Space](https://huggingface.co/spaces/PeetPedro/headroom-eval)
- [kompress paper (ICLR 2027)](https://kompress.vaked.dev)
- [LoopKit](https://github.com/peterlodri-sec/loopkit)
- [pocoo.vaked.dev](https://pocoo.vaked.dev)
- [Sponsor](https://github.com/sponsors/peterlodri-sec)

## Build

```bash
make build-upx   # 2.6MB binary
```
