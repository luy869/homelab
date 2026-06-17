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
- chat_bot の生成 LLM を qwen3.5:9b → より高性能なモデル（gemma2:27b / qwen2.5:32b など）に切り替えて、技術帰属問題の解消を試す
- chat_bot に環境変数で LLM モデル指定（`OLLAMA_LLM_MODEL` 等）を追加
- 検索結果のチャンク内容を返すデバッグエンドポイント（精度問題のローカル診断用）
