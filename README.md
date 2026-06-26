# 🐋 headroom-eval-cli

Interactive pure-ASCII TUI for the [Headroom Eval Space](https://huggingface.co/spaces/PeetPedro/headroom-eval).
Live status, animated ASCII dodecahedron, steganography, genesis seal.

## Install

```bash
go install github.com/peterlodri-sec/headroom-eval-cli@latest
headroom-eval
```

Or using the one-line installer (downloads a pre-built binary):

```bash
curl -fsSL https://raw.githubusercontent.com/peterlodri-sec/headroom-eval-cli/main/install.sh | bash
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
headroom-eval --help      # usage and key reference
headroom-eval --seal      # genesis hash
headroom-eval --version   # version + build time
```

## headroom-logs

`headroom-logs` is a companion CLI for streaming Hugging Face Space logs with color output.

```bash
go install github.com/peterlodri-sec/headroom-eval-cli/cmd/logs@latest
```

```
headroom-logs [flags]

Flags:
  --token    HF API token (default: $HF_TOKEN)
  --space    HF Space owner/name (default: PeetPedro/headroom-eval)
  --mode     run or build (default: run)
  -f         follow / stream logs
  --filter   show only lines containing STR
  --color    force color output
```

Example:

```bash
export HF_TOKEN=hf_...
headroom-logs -f                      # stream run logs
headroom-logs --mode build            # last build log line
headroom-logs -f --filter ERROR       # stream, errors only
```

## Build

```bash
make build-upx   # 2.6MB binary
```

## Links

- [Headroom Eval Space](https://huggingface.co/spaces/PeetPedro/headroom-eval)
- [kompress paper (ICLR 2027)](https://kompress.vaked.dev)
- [LoopKit](https://github.com/peterlodri-sec/loopkit)
- [pocoo.vaked.dev](https://pocoo.vaked.dev)
- [Sponsor](https://github.com/sponsors/peterlodri-sec)
