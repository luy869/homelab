# homelab — 自宅サーバー統合モノレポ

自宅サーバー（`luy-XA7C-R38`）で公開運用する個人開発アプリ群を 1 リポジトリで束ね、起動・設定・構成を一元管理する。設計背景は [docs/architecture.md](docs/architecture.md)、決定の経緯は [DEVLOG.md](DEVLOG.md) を参照。

## 技術スタック / 構成
- オーケストレーション: Docker Compose（ルート `compose.yaml` が各アプリ prod compose を `include`）
- リバースプロキシ: ホスト nginx（:80）→ `proxy/nginx/`
- 公開: Cloudflare Tunnel → `proxy/cloudflared/`（IP 秘匿・自動 TLS）
- ローカル LLM: Ollama（ホスト :11434）→ `services/ollama.md`
- 画像推薦: CLIP gRPC（ホスト :50051）→ `services/clip.md`
- アプリ（submodule）:
  - `apps/chat_bot` — RAG ポートフォリオチャット（FastAPI, container `rag-platform` :8000）
  - `apps/palette-vein` — CLIP 画像推薦（Go + Python CLIP + pgvector, frontend :8090）

## 起動手順
1. `git clone --recurse-submodules <repo>`（既存 clone は `git submodule update --init`）
2. 各アプリに `.env` を配置: `apps/chat_bot/.env` / `apps/palette-vein/.env`（`.env.example` 参照）
3. ホストサービス（Ollama / CLIP / nginx / cloudflared）が稼働中か確認 → `./scripts/status.sh`
4. `./scripts/up.sh`（= `docker compose up -d --build`）
5. `./scripts/status.sh` で死活確認

## 設計方針（wrap & tidy）
- ホスト常駐サービス（Ollama/CLIP/nginx/Tunnel）はコンテナ化せず据え置き（GPU アクセス・日常/ゲーム併用機のため）。設定のみ版管理。
- 各アプリは独立 GitHub リポのまま submodule 参照。アプリのコード変更は各リポで行い、homelab は submodule ポインタを更新する。
- 秘密情報は `apps/<app>/.env` と `proxy/cloudflared/*.json` に隔離し gitignore。

## 注意点
- `include` 統合で compose プロジェクト名が `homelab` になり、palette-vein のコンテナ名が `palettevein-*` → `homelab-*` に変化する。永続ボリュームは明示名 `palettevein_pgdata` のため DB データは保持される。
- ホスト nginx は `root /home/luy/chat_bot/frontend` を配信。chat_bot を `apps/chat_bot` に集約した場合は nginx の root パス更新が必要（カットオーバー時）。

## マイルストーン
- [x] v1 設計（plan 承認済み）
- [~] v1 スキャフォールド（リポ雛形・submodule・compose・README/構成図・scripts）← 進行中
- [ ] v1 カットオーバー（サーバーを homelab 起点デプロイへ切替）※ユーザー確認後に実施
- [ ] v2 監視ダッシュボード（`dashboard/`, dash.luy869.net）
- [ ] v3 Open WebUI（chat.luy869.net）
