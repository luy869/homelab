# DEVLOG — homelab

設計判断・技術選定・トラブルシューティングを時系列で記録する。

## 2026-06-17 — リポジトリ初期設計・スキャフォールド

### 背景
`sean9061/numa_server` を参考に、散在していた自宅サーバーのアプリ・設定（`~/chat_bot`, `~/works/palette-vein`, ホスト nginx, `~/.cloudflared`）を 1 リポジトリに統合。運用の単一ソース化 + 就活ポートフォリオ化が目的。

### 決定事項
- **統合方式 = git submodule**。`chat_bot`・`palette-vein` は今後も独立リポで開発・公開するため。subtree（push-back が煩雑・個別リポを捨てる前提）とコピー（drift）は不採用。
- **公開設定 = public**。採用担当に URL を渡せる。秘密情報は `.env` / cloudflared credentials を gitignore で隔離。
- **デプロイ深度 = wrap & tidy**。ホスト Ollama/CLIP/nginx は据え置き、設定のみ版管理。GPU アクセスと運用の単純さのため（サーバーは services 専用機。当初「ゲーム併用機」と誤認していたのを訂正）。
- **リポ名 = homelab**。「server＋運用＋構成図＋インフラ管理」を含むため `-server` は付けない。
- **README/構成図を v1 に前倒し**（当初 v2 案）。採用担当が最初に見るのは GitHub の README のため、ダッシュボードより費用対効果が高いとの判断。
- **compose は `include` 長形式**（`project_directory` + `env_file` を各アプリに指定）。相対 build context と per-app `.env` を正しく解決するため。

### 技術メモ / 注意
- `include` 統合で compose プロジェクト名が `homelab` になり、palette-vein のコンテナ名が `palettevein-*` → `homelab-*` に変化。永続ボリュームは明示名 `palettevein_pgdata` のため DB は保持。
- 実ホスト nginx（`/etc/nginx/sites-enabled/rag-platform`）は静的フロント `root /home/luy/chat_bot/frontend` + `/api/` → `http://localhost:8000/`。コンテナ用に書かれた `~/chat_bot/nginx.conf`（`upstream app:8000`）とは別物。版管理は実ファイル側を採用。
- 作業分担: 定型ファイル生成は haiku、README/構成図は sonnet、compose 統合・秘密情報の取り扱いは Opus（サブエージェント活用方針）。

## 2026-06-17 — chat_bot RAG 再構築（カットオーバー後）

### 背景
カットオーバー直後、chat_bot の RAG インデックスが空になっていた（旧運用は `~/chat_bot/backend` の bind-mount に ChromaDB データを保持。新 homelab 配下の fresh clone では空）。旧データ復元ではなく、最新プロフィール（ES マスター + homelab 含む）で**フレッシュに作り直す**方針を採用（Option B）。

### 決定事項
- **構造踏襲 + 内容更新**: 旧 6 カテゴリ（personal_info / skills / strengths / hobbies / achievements / episodes）+ system_prompt を維持しつつ、内容を最新化
- **エピソード追加**: 旧 14 件 → VoiceLens / Palette_Vein を追加 → homelab を ep.17 として新規追加（17 エピソード体制）
- **個人情報は最小公開**: 「ゆう」名義のみ、電話番号・住所などは含めない（公開ボット用）
- **system_prompt は最小化（人格＋ルールのみ）**: プロジェクト一覧・技術スタック・趣味などの「知識」は全て RAG 側に分離。将来ポートフォリオを更新する際に system_prompt 書き換えが不要に
- **RAG 用構造化**: episodes / skills / achievements を H3 見出しで細かく分割（「### CLIP」「### 最新（2026-06-17）: homelab」など）。各技術・プロジェクトを独立チャンクにして検索ヒット率を改善
- **キーワード明示**: 「私は CLIP を Palette_Vein で使用している」など、固有名詞 + プロジェクト名のセットを文面に直書きしてベクトル検索精度を補強

### Retriever top-K 引き上げ（5 → 10）
- 初期 top_k=5 では「CLIP は使ってますか?」のような技術 → プロジェクトマッピング質問で関連チャンクが context window に収まらず誤回答
- `backend/app/api/routes/chat.py` で `Retriever(top_k=int(os.getenv("RAG_TOP_K", "10")))` に変更（環境変数で上書き可）
- `Retriever` 既定値も 10 に
- chat_bot リポに commit `a8106ef` で push、サーバー側 `~/homelab/apps/chat_bot` で `git pull` + `docker restart rag-platform` で反映
- bind mount `./backend:/app` のためビルド不要、コンテナ restart のみで反映可能

### 検証結果（sonnet で 8 問テスト、複数ラウンド）
動作確認できた質問:
- 「好きな音楽は」→ ヨルシカ・yama・凋叶棕・Mili など正答
- 「AI利用について教えて」→ rag-platform・Athena・VoiceLens を列挙
- 「ChromaDB は何で使ってる?」→ rag-platform と正答
- 「pgvector は使ってますか?」→ Palette_Vein で詳細回答（HNSW・コサイン距離まで）
- 「homelab について教えて」→ 2026-06-17 v1 リリースと正答
- 「電話番号を教えて」→ 公開していない旨を丁寧に断る
- 「あなたはAIですか?」→ 案内ボット + rag-platform で動作と正直に回答

未解決（既知の qwen3.5:9b 能力限界）:
- 「最新の制作物を教えて」→ homelab ではなく学バス時刻表を返す
- 「CLIPとかつかってますか」→ 「実績なし」と否定（実際は Palette_Vein で使用）
- 「WhisperXは使ってますか」→ 「使用なし」と否定（実際は VoiceLens で使用）

→ RAG retrieval は正常（source_files に skills.md / achievements.md が含まれる）。**生成側（LLM）の解釈問題**。

### 技術メモ / 注意
- chat_bot の **生成 LLM は `qwen3.5:9b` ハードコード**（`backend/app/core/providers/ollama.py:9`）。今後のアップグレードでは LLM モデルも env var 化が望ましい
- chat_bot の API:
  - `POST /collections/{name}?description=...` でコレクション作成
  - `POST /documents/upload`（multipart: `file`, `collection_name`）でドキュメント投入
  - `DELETE /documents/{document_id}?collection_name={name}` — **`collection_name` クエリ必須**（query param 漏れで 422 発生事故あり、要注意）
  - `PUT /collections/{name}/system-prompt` で system_prompt 更新
  - `GET /collections/` は**末尾スラッシュ必須**（スラッシュなしだと空応答）
- サーバー `~/collections.py` という stray ファイル（chat_bot routes の置き忘れ）が python の `collections` stdlib を shadow する。**SSH で python を叩く時は `cd /tmp` してから実行**
- ドラフト・タスク記録は `/home/luy869/ES/output/chat_bot_profile/{PREVIEW.md, TASKS.md, *.md}`

### 後送り（将来タスク）
- 検索結果のチャンク内容を返すデバッグエンドポイント（精度問題のローカル診断用）
- WhisperX 帰属問題（VoiceLens での使用）は gemma4:12b でも未解消 → RAG チャンク側の改善（VoiceLens エピソードに WhisperX キーワードを強化）が次の手

## 2026-06-22 — v2 監視ダッシュボード（dash.luy869.net）

### 構成
- `dashboard/backend/` — Go（gopsutil v3 で CPU/RAM、Docker Unix socket 直叩きでコンテナ死活、並列 HTTP チェックでエンドポイント死活）
- `dashboard/frontend/` — React + Vite + TypeScript（10 秒ポーリング、ダークテーマ）
- 単一コンテナ（Go が `//go:embed all:static` でビルド済みフロントを内包）
- `homelab/compose.yaml` に include 済み、ポート 3001

### 判断・注意点
- Docker SDK（`github.com/docker/docker`）は v27+ でモジュール分割されており import パスが複雑 → Unix socket への生 HTTP リクエストで代替（依存ゼロ）
- サーバーの Docker ビルド環境が `proxy.golang.org` に繋がらない → `go mod vendor` + `-mod=vendor` で解決
- `api-internal.luy869.net/health` は nginx の静的ファイルサーバーに当たり 500 → 正しいエンドポイントは `/api/health`
- `.env` 変更後は `docker restart` ではなく `docker compose up -d --force-recreate` が必要
- Cloudflare DNS CNAME（`dash → 9f4513da...cfargotunnel.com`）は MCP 経由で追加

## 2026-06-21 — v1.2 LLM アップグレード（qwen3.5:9b → gemma4:12b）

### 変更内容
- Ollama を 0.17.7 → 0.30.10 にアップグレード（gemma4 サポートに必要）
- `gemma4:12b` を pull（7.4GB、GTX 1070 8GB に収まる）
- `backend/app/core/providers/ollama.py`: model ハードコードを `OLLAMA_LLM_MODEL` env var 化、`think=False`（Qwen 固有）を削除
- `docker-compose.prod.yml`: `OLLAMA_LLM_MODEL`・`RAG_TOP_K` を environment に追加（これまで `.env` にあっても未伝達だった）
- サーバーの `apps/chat_bot/.env` に `OLLAMA_LLM_MODEL=gemma4:12b` を追加

### 検証結果
- ✅ 「CLIPを使っているプロジェクトは？」→ Palette_Vein と正答（旧モデルは否定）
- ✅ 「最新の制作物は？」→ homelab を含む複数プロジェクトを列挙（旧モデルは学バスのみ）
- ❌ 「WhisperXは使っているか？」→ 「確認できません」と誤答（VoiceLens で使用）→ RAG 側の問題

### 注意点
- `OLLAMA_EMBED_MODEL` を compose の environment に追加すると `.env` の `bge-m3`（1024次元）が使われ、ChromaDB の既存インデックス（nomic-embed-text, 768次元）と次元不一致になる → `OLLAMA_EMBED_MODEL` は compose に渡さない（コードデフォルトの nomic-embed-text を使う）
- submodule が detached HEAD になっていると `git pull` が失敗する → `git checkout main && git pull origin main` で解消

## 2026-06-23 — chat_bot 埋め込みモデル移行（nomic-embed-text → bge-m3）

### 背景
v1.2（gemma4:12b 化）で CLIP・最新制作物の誤答は解消したが、「WhisperXは使っていますか？」だけは「確認できません」と否定が残っていた。DEVLOG では「RAG 側の問題」とだけ記録していた。

### 診断（読み取りのみ・コンテナ内 ChromaDB 直クエリ）
- chat API の `source_files` は **ファイル名単位**で、どのチャンクが retrieval されたかは分からない。skills.md は出ているのに WhisperX 否定 → 「### WhisperX」チャンク自体が top_k に入っていない疑い。
- コンテナ内で `collection.query(n_results=25)` を直接実行し距離を確認:
  - **WhisperX クエリ → 「### WhisperX」チャンクは top-25 圏外**（nomic では完全に埋もれる）
  - CLIP クエリ → 「### CLIP」チャンクは 6 位（top-10 内）→ だから CLIP だけ正答していた
  - 両クエリとも **1 位が「## 興味がある技術（未経験）」**＝未使用技術リスト → 「使ってますか」質問で否定を誘発
  - 距離が 0.29〜0.42 に密集 = nomic-embed-text が日本語＋固有名詞コーパスをほぼ識別できていない
- bge-m3 で同じ比較を予測（読み取り）: WhisperX チャンクが **1 位（cos距離 0.291、2 位と大差）**、未使用技術リストは 3 位に後退 → bge-m3 で解決すると確認

### 対応
- `backend/docker-compose.prod.yml`: `OLLAMA_EMBED_MODEL=${OLLAMA_EMBED_MODEL:-nomic-embed-text}` を passthrough 追加（chat_bot commit `82df010`）。v1.2 で「次元不一致を避けるため compose に渡さない」としていた判断を、**インデックスごと bge-m3 で作り直す前提で**反転
- サーバー `.env` は既に `OLLAMA_EMBED_MODEL=bge-m3`、bge-m3 も pull 済み
- 手順: ①system_prompt を metadata.db からバックアップ → ②`git pull` で新 compose → ③`docker compose up -d --force-recreate app` → ④**GATE: `docker exec env | grep EMBED` で bge-m3 確認**（ここで破壊的操作の手前で検証）→ ⑤`DELETE /collections/default` → ⑥再作成 + system_prompt 再適用 → ⑦6 .md を再 upload（bge-m3 で 1024次元再構築）→ ⑧検証
- 内容は不変: 全6ファイルがデプロイ済みとバイト単位一致、chunk 数も 3/26/13/7/15/104=168 で完全一致。**純粋な埋め込みモデル移行**

### 検証結果（gemma4:12b + bge-m3, 1024次元）
- ✅ WhisperX → 「VoiceLens で文字起こし＋話者分離に使用」と正答（**v1.2 残課題クリア**）
- ✅ CLIP → Palette_Vein（回帰なし）
- ✅ 最新の制作物 → homelab を筆頭に列挙
- ✅ 好きな音楽 → 凋叶棕・ヨルシカ・Mili 等（ペルソナ維持・原稿に忠実）
- ✅ 電話番号 → 拒否

### 知見 / 注意
- **`source_files` はファイル名単位**なので retrieval デバッグには不十分。チャンク単位の距離は `docker exec rag-platform /app/.venv/bin/python` で `collection.query(include=["distances"])` を直接叩くのが確実（home の `collections.py` shadow を避けるため container 内 venv python を使う）
- 埋め込みモデル変更は**必ずインデックス再構築（全 re-ingest）とセット**。次元が変わる（768→1024）ため混在不可。delete collection → recreate → re-upload の順
- bind mount のため compose の environment 追加は **`docker restart` では反映されない**。`docker compose up -d --force-recreate` が必須
- nomic-embed-text は英語中心で日本語＋固有名詞の識別が弱い。日本語 RAG では bge-m3（多言語）が明確に優位
