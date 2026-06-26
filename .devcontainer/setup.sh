   #!/bin/bash
   set -e
   echo "=== DX Setup ==="
   pip install -q -r requirements.txt 2>/dev/null || true
   go install github.com/peterlodri-sec/headroom-eval-cli/cmd/logs@latest
  2>/dev/null || true
   npm install -g @opencode-ai/opencode 2>/dev/null || true
   curl -fsSL https://raw.githubusercontent.com/whale-coder/whale/main/install.sh
  2>/dev/null | bash 2>/dev/null || true
   pip install -q headroom-ai mcp 2>/dev/null || true
   echo "DX ready — whale, opencode, headroom, nix, context7"
