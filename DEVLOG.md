# DEVLOG — homelab

設計判断・技術選定・トラブルシューティングを時系列で記録する。

## 2026-06-17 — リポジトリ初期設計・スキャフォールド

### 背景
`sean9061/numa_server` を参考に、散在していた自宅サーバーのアプリ・設定（`~/chat_bot`, `~/works/palette-vein`, ホスト nginx, `~/.cloudflared`）を 1 リポジトリに統合。運用の単一ソース化 + 就活ポートフォリオ化が目的。

### 決定事項
- **統合方式 = git submodule**。`chat_bot`・`palette-vein` は今後も独立リポで開発・公開するため。subtree（push-back が煩雑・個別リポを捨てる前提）とコピー（drift）は不採用。
- **公開設定 = public**。採用担当に URL を渡せる。秘密情報は `.env` / cloudflared credentials を gitignore で隔離。
- **デプロイ深度 = wrap & tidy**。ホスト Ollama/CLIP/nginx は据え置き、設定のみ版管理。GPU アクセスと日常/ゲーム機運用を壊さないため。
- **リポ名 = homelab**。「server＋運用＋構成図＋インフラ管理」を含むため `-server` は付けない。
- **README/構成図を v1 に前倒し**（当初 v2 案）。採用担当が最初に見るのは GitHub の README のため、ダッシュボードより費用対効果が高いとの判断。
- **compose は `include` 長形式**（`project_directory` + `env_file` を各アプリに指定）。相対 build context と per-app `.env` を正しく解決するため。

### 技術メモ / 注意
- `include` 統合で compose プロジェクト名が `homelab` になり、palette-vein のコンテナ名が `palettevein-*` → `homelab-*` に変化。永続ボリュームは明示名 `palettevein_pgdata` のため DB は保持。
- 実ホスト nginx（`/etc/nginx/sites-enabled/rag-platform`）は静的フロント `root /home/luy/chat_bot/frontend` + `/api/` → `http://localhost:8000/`。コンテナ用に書かれた `~/chat_bot/nginx.conf`（`upstream app:8000`）とは別物。版管理は実ファイル側を採用。
- 作業分担: 定型ファイル生成は haiku、README/構成図は sonnet、compose 統合・秘密情報の取り扱いは Opus（サブエージェント活用方針）。
