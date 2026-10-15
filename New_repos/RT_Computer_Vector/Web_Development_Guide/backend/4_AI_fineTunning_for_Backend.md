# AI fine Tunning guide (for dummies like me)

## Step 1: Picking a base model 

Picking the right base model is one of the highest-leverage decisions in a fine-tuning project. Get it wrong, and no amount of hyperparameter tuning or extra training data will fully compensate. Get it right, and you'll spend less on compute, hit your accuracy targets faster, and avoid legal headaches down the line. This guide walks through the full decision process, step by step.

### Why the Base Model Choice Matters

Fine-tuning adapts a model's existing knowledge to a narrower task — it doesn't teach the model from scratch. That means the base model's pretraining data, architecture, and scale set a ceiling on what fine-tuning can realistically achieve. A model with no exposure to code will struggle to become a strong coding assistant no matter how much you fine-tune it, and a model that's too small may lack the capacity to learn nuanced behavior even with a perfect dataset. Choosing well upfront saves you from discovering these limits after you've already spent your compute budget.

### Step 1: Define Your Task Requirements

Before comparing models, write down what you actually need:

- **Task type** — classification, generation, summarization, code completion, conversational agent, etc.
- **Input/output length** — short-form responses vs. long documents or multi-turn context.
- **Latency and throughput needs** — real-time inference vs. batch processing.
- **Deployment environment** — cloud API, on-prem server, or edge/local device.
- **Accuracy bar** — what's "good enough" for this to ship.

These constraints will immediately rule out some models (e.g., a 70B model is a poor fit for on-device deployment) and narrow your search considerably.

### Step 2: Match Pretraining Domain to Your Use Case

Models specialize based on what they were trained on. A model pretrained heavily on source code will adapt faster to coding tasks than a general-purpose text model; a model trained on multilingual corpora will transfer better to non-English tasks. Look for:

- Published details on the pretraining data mix (when available).
- Existing benchmark results in your domain (code, math, medical, legal, multilingual, etc.).
- Evidence of instruction-tuning or chat-tuning if your task is conversational, since this affects how much fine-tuning is needed to get usable behavior.

If your domain is niche (e.g., legal contracts, chemistry notation), prioritize models with any prior exposure to that kind of text over models that are simply larger.

### Step 3: Weigh Model Size Against Compute Budget

Bigger models generally generalize better and follow instructions more reliably, but the cost curve is steep:

- **Training cost** scales with parameter count — memory, GPU-hours, and storage all increase.
- **Inference cost** is often the bigger long-term expense, especially at production scale.
- **Fine-tuning method matters here too** — full fine-tuning of a large model may be infeasible, whereas parameter-efficient methods (LoRA, QLoRA, adapters) can make even large models tractable on modest hardware.

A practical rule: pick the smallest model that can plausibly hit your accuracy bar, and validate that assumption with a small pilot run before committing to a full training cycle. Going bigger "just in case" is usually a waste of budget.

### Step 4: Evaluate Benchmark Performance

Public benchmarks are a useful filter, not a guarantee:

- Look for benchmarks that resemble your actual task (don't over-index on general leaderboards if your task is narrow).
- Treat benchmark scores as a starting shortlist, not a final ranking — your own dataset quality and task framing will affect real-world results more than a half-point difference on a public eval.
- Where possible, run your own small-scale comparison: fine-tune 2–3 candidate models on a subset of your data and compare results directly.

### Step 5: Check the License

This step gets skipped more often than it should, and it can be the most expensive mistake on this list:

- Confirm the license permits your intended use — some open-weight models restrict **commercial use**, **redistribution**, or use of outputs to train other models.
- Check for usage caps (e.g., licenses that change terms above a certain number of users or revenue threshold).
- If you plan to redistribute a fine-tuned checkpoint, verify that's explicitly allowed, not just implied.

Confirm licensing *before* investing compute — reworking a project after training because of a license conflict is a painful and avoidable cost.

### Step 6: Consider Context Length and Architecture

If your task involves long documents, multi-turn conversations, or large retrieved context, check the base model's native context window and how well it actually performs at long range — advertised max length and effective usable length aren't always the same. Also consider architectural details relevant to your deployment: quantization support, compatibility with your fine-tuning framework, and whether the model has known quirks (e.g., poor tool-calling support, weak multi-turn coherence) that would require extra work to fix.

### Step 7: Assess Ecosystem and Tooling Support

A technically strong model with poor tooling support can slow you down significantly. Check for:

- Compatibility with common fine-tuning frameworks (Hugging Face, Axolotl, etc.).
- Active community support, documentation, and known-issue tracking.
- Availability of quantized or distilled variants if you'll need to compress the model post-training.
- Precedent — has anyone published successful fine-tunes of this model for a similar task? That's a strong practical signal.

### Common Mistakes to Avoid

- **Defaulting to the largest available model** without validating that the task needs that capacity.
- **Ignoring the license** until after training is complete.
- **Over-trusting leaderboard rankings** for tasks that don't resemble the benchmark.
- **Skipping a small pilot run** before committing to full-scale fine-tuning.
- **Underestimating inference cost**, which usually dwarfs training cost over a model's production lifetime.

### Quick Reference Checklist

- [ ] Task requirements defined (type, length, latency, deployment target)
- [ ] Pretraining domain matches or is adjacent to your use case
- [ ] Model size validated against compute budget and a pilot run
- [ ] Benchmarks checked against tasks similar to yours
- [ ] License confirmed for your specific commercial/redistribution needs
- [ ] Context length and architecture fit deployment constraints
- [ ] Tooling and community support confirmed

### Example: Popular Base Models at a Glance

The table below is a snapshot of widely used open-weight base model families as of late 2026. Specs, licenses, and rankings shift quickly in this space, so treat this as a starting shortlist and verify current details before committing.

| Model Family | Typical Size Range | Best For | License Notes | Key Strength |
|---|---|---|---|---|
| Llama 4 (Meta) | 1B – 400B+ | General-purpose tasks, broad tooling support | Community license, free below a large monthly-user threshold | Largest ecosystem and community fine-tuning resources |
| Qwen 2.5 / 3.5 (Alibaba) | 0.5B – 72B | Multilingual tasks, coding (via Coder variants) | Mostly Apache 2.0 | Widest size range, strong coding-specific checkpoints |
| Mistral Small / Large (Mistral AI) | 7B – 123B | Efficient reasoning, function calling, non-English scripts | Apache 2.0 (Small), commercial license (Large) | Strong performance per parameter, efficient tokenizer |
| Gemma 3 / 4 (Google DeepMind) | 1B – 27B | High-quality small models, on-device deployment | Gemma license (commercial use with some restrictions) | High intelligence density per parameter |
| Phi-4 (Microsoft) | 3.8B – 14B | Math, reasoning, and low-resource/edge deployment | MIT | Trained on curated synthetic data, runs on minimal hardware |
| DeepSeek-V4 / R1-Distill (DeepSeek) | 14B – 670B (MoE) | Advanced reasoning, distilled reasoning in smaller models | MIT | Strong reasoning ability distilled into compact dense models |
| GLM-5.2 (Zhipu AI) | Large MoE | Coding-focused tasks and agentic workflows | Open-weight | Leads open-weight coding benchmarks |
| Nemotron 3 (NVIDIA) | ~30B – ~550B | Tool-use agents, long-horizon planning, throughput-sensitive pipelines | OpenMDW-1.1 (largest tier fully open source) | Purpose-built for multi-agent, high-throughput systems |

For most fine-tuning projects on a single consumer or prosumer GPU, models in the 7B–14B range with QLoRA are the practical sweet spot — enough capability for strong task performance without needing multi-GPU infrastructure.

## Step 2: Finding datasets 

### Finding Datasets for Fine-Tuning

Once you've picked a base model, the next bottleneck is almost always data. A great model fine-tuned on a small, messy, or mismatched dataset will underperform a modest model fine-tuned on clean, well-targeted data. This guide covers how to think about sourcing datasets, what to check before you commit to one, and where to actually go looking.

### Why Dataset Quality Matters More Than Quantity

Fine-tuning works by showing the model examples of the behavior you want. If those examples are inconsistent, mislabeled, or only loosely related to your actual task, the model learns the noise along with the signal. A few thousand clean, well-labeled, on-task examples routinely outperform a much larger dataset scraped without curation. Before searching for a dataset, it helps to be clear-eyed about this trade-off: scale is useful, but only once quality is established.

### Step 1: Define What You're Actually Looking For

Be specific before you start searching:

- **Task format** — classification, instruction-following, conversational turns, summarization pairs, code completion, etc.
- **Domain specificity** — general-purpose vs. a narrow field like legal, medical, or financial text.
- **Language(s) and scripts** required.
- **Volume needed** — this varies hugely by task and fine-tuning method (parameter-efficient methods often need far less data than full fine-tuning).
- **Label quality requirements** — do you need human-verified labels, or is weak/noisy labeling acceptable for your use case?

A vague search ("chatbot dataset") returns thousands of irrelevant results. A specific one ("multi-turn customer support dialogues in English with resolution labels") narrows things down fast.

### Step 2: Check Licensing Before You Download Anything

This is the step most likely to bite you later:

- Confirm the dataset's license permits your intended use — many datasets are restricted to **non-commercial** or **research-only** use.
- Check for **attribution requirements** you'll need to satisfy in a shipped product.
- Watch for datasets built by scraping platforms with their own terms of service (e.g., social media, forums) — the data being publicly visible doesn't always mean it's freely licensed for training.
- If a dataset's license is unclear or unstated, treat that as a red flag rather than an oversight to ignore.

### Step 3: Evaluate Data Quality Before Committing

Before building a full pipeline around a dataset, spot-check it:

- Pull a random sample of rows and manually review them for label accuracy and consistency.
- Check for **duplicate or near-duplicate entries**, which can silently bias training.
- Look for **class imbalance** if it's a classification task — a dataset that's 95% one label will teach the model to just guess that label.
- Check for **contamination** — overlap between the dataset and any benchmark you plan to evaluate on, which would make your evaluation results meaningless.

### Step 4: Match Dataset Structure to Your Fine-Tuning Format

Different fine-tuning approaches expect different data shapes:

- **Instruction fine-tuning** needs prompt/response pairs, often with a system or instruction field.
- **Preference tuning (DPO/RLHF-style)** needs paired responses with a preference label (chosen vs. rejected).
- **Classification** needs text paired with a clean label taxonomy that matches your production categories.
- **Domain adaptation / continued pretraining** can often use raw unlabeled text, which loosens quality requirements somewhat but raises the importance of topical relevance.

Reformatting a dataset after the fact is doable but time-consuming — checking the schema upfront saves rework.

### Step 5: Consider Building or Augmenting Your Own

Public datasets won't always match your task closely enough. Common ways to close the gap:

- **Synthetic data generation** using a capable model to produce examples in your target format, then filtering for quality.
- **Combining multiple datasets** — mixing a general instruction dataset with a smaller domain-specific set often improves both task performance and generalization.
- **Human-in-the-loop labeling** for a smaller, high-precision seed set, especially when public data doesn't exist for a niche domain.

A hybrid approach — some public data, some synthetic, a smaller high-quality human-reviewed core — is common in practice rather than relying on any single source.

### Step 6: Hold Out a Real Evaluation Set

Set aside a portion of your data — or find a separate dataset entirely — purely for evaluation, and don't let it leak into training. It's easy to accidentally contaminate your eval set when combining multiple sources, so check overlap explicitly (e.g., via hashing or near-duplicate detection) rather than assuming a clean split.

### Common Mistakes to Avoid

- **Downloading first, checking the license later** — sometimes after it's too late to swap the dataset out.
- **Trusting dataset size over dataset quality.**
- **Skipping a manual spot-check** of sample rows before building a full pipeline.
- **Ignoring benchmark contamination**, which inflates evaluation results and hides real-world weaknesses.
- **Using a single source** when a blended dataset would generalize better.

### Quick Reference Checklist

- [ ] Task format, domain, language, and volume needs clearly defined
- [ ] License confirmed for your specific use (commercial, redistribution, attribution)
- [ ] Sample rows manually reviewed for label quality and consistency
- [ ] Duplicates, class imbalance, and benchmark contamination checked
- [ ] Dataset schema matches your fine-tuning method's expected format
- [ ] Gaps filled with synthetic data, blended sources, or targeted human labeling if needed
- [ ] A clean, non-overlapping evaluation set held out

### Example: Well-Known Sites for Finding Datasets

| Site | Best For | Notes |
|---|---|---|
| Hugging Face Datasets | LLM fine-tuning datasets of nearly every kind | Huge library with built-in filtering by task, size, and language; integrates directly with common training libraries |
| Kaggle Datasets | General ML and data science projects | Large, searchable catalog with community notebooks showing how each dataset has been used |
| UCI Machine Learning Repository | Classic structured/tabular ML datasets | Long-standing academic repository, well-documented but skews toward smaller, traditional ML tasks |
| Google Dataset Search | Broad discovery across the web | A search engine specifically for datasets hosted anywhere online, useful for tracking down niche or domain-specific sources |
| Papers With Code (Datasets) | Datasets tied to published research and benchmarks | Good for finding the exact dataset behind a paper's reported results, with linked leaderboards |
| AWS Open Data Registry | Large-scale public datasets hosted on cloud infrastructure | Useful when you need datasets already staged for cloud-based training pipelines |
| data.gov | US government and public sector data | Strong for civic, economic, health, and demographic data |
| OpenML | Structured ML datasets with built-in benchmarking | Datasets come with metadata and prior benchmark results, useful for quick comparison |
| Common Crawl | Raw web text at massive scale | A common source for pretraining or domain-adaptation corpora, requires heavy filtering before use |
| LLMDataHub (GitHub) | Instruction fine-tuning datasets specifically | Community-curated index pointing to popular instruction-tuning and alignment datasets |

As with base models, treat this table as a starting point rather than an exhaustive list — new domain-specific and community-curated dataset hubs appear regularly, so it's worth searching for options specific to your field before settling on a general-purpose source.

## step 3: Dataset formatting 

### Formatting Datasets for Fine-Tuning

Having the right dataset isn't enough if it's shaped wrong. Every model family expects its training examples in a specific structure — the wrong chat template, a missing special token, or an inconsistent schema can silently degrade results even when the underlying data is excellent. This guide covers how to think about formatting, what to check before training, and which formats pair with which models and use cases.

### Why Format Choice Matters

Fine-tuning teaches a model to predict text in the exact shape it's shown. If your training data doesn't match the format the model was originally instruction-tuned with — its chat template, its special tokens, its role structure — the model has to relearn that structure on top of learning your task, which wastes capacity and often produces inconsistent behavior at inference time. Getting the format right up front means every training example is doing useful work.

### Step 1: Match Format to Your Fine-Tuning Method

The type of fine-tuning you're doing determines the shape of each example:

- **Supervised fine-tuning (SFT) / instruction tuning** — prompt-response pairs, often wrapped in a chat template with role fields (system/user/assistant).
- **Preference tuning (DPO, RLHF-style methods)** — a prompt paired with two responses: one preferred (chosen) and one rejected.
- **Classification** — raw text paired with a label from a fixed taxonomy.
- **Continued pretraining / domain adaptation** — plain unlabeled text, no roles or turns needed.
- **Embedding / retrieval fine-tuning** — pairs or triplets: an anchor, a positive match, and optionally a hard negative.

Picking the wrong shape for your method is the single most common formatting mistake — a classification dataset won't work for DPO, and vice versa.

### Step 2: Choose the Right File-Level Format

Beyond the internal structure, the file format matters for tooling compatibility:

- **JSONL (JSON Lines)** — the de facto standard for LLM fine-tuning; one JSON object per line, easy to stream and shuffle.
- **CSV/TSV** — common for classification or tabular tasks, but awkward for multi-turn conversations or nested fields.
- **Parquet** — efficient for large-scale datasets, widely supported by data processing libraries, good when volume is high.
- **WebDataset / manifest formats** — used for audio and image data, pairing a media file path with a JSON metadata record (transcript, caption, labels).

JSONL is the safest default for anything text-based unless your training framework specifically expects something else.

### Step 3: Get the Chat Template Right for Your Target Model

Instruction-tuned models are trained with a specific chat template, and fine-tuning data should follow that same template so the model doesn't have to relearn its own conversational structure. Templates differ in how they mark roles and turn boundaries — some use explicit header tags, some use bracketed instruction markers, some use a single repeated turn token. Most modern training frameworks (e.g., Hugging Face `transformers`) can apply a model's official chat template automatically from a simple role/content list, which is generally safer than hand-formatting the template yourself.

### Step 4: Handle Special Tokens and System Prompts Correctly

- Don't manually insert a model's special tokens (like end-of-turn or header markers) into the text field — let the tokenizer's chat template add them, or you risk double-tokenizing or misaligning them.
- If a model's template doesn't support a system role, decide upfront whether to fold system instructions into the first user turn or drop them.
- Keep system prompts consistent across training examples if your production use case will also use a fixed system prompt — mismatches between training-time and inference-time system prompts can hurt performance.

### Step 5: Structure Multi-Turn and Multimodal Data Properly

- For multi-turn conversations, keep the full turn history in each example rather than flattening it into a single block of text — most chat templates expect a list of role/content turns.
- For vision-language fine-tuning, pair each image (as a file path, URL, or base64 blob) with its associated text turns in the same record, following the target model's expected multimodal schema.
- For audio, a manifest-style format (one JSON record per line, each with an audio path, transcript, and duration) is standard across most speech fine-tuning frameworks.

### Step 6: Validate Before You Train

Before launching a full training run:

- Render a handful of formatted examples through the actual tokenizer and read the output — confirm turns, roles, and special tokens look correct.
- Check token length distribution against your model's context window; truncated examples can silently corrupt training signal.
- Confirm your schema is consistent across every row — a single malformed record can break a training job partway through.

### Common Mistakes to Avoid

- **Hand-writing special tokens** instead of letting the tokenizer's chat template apply them.
- **Mixing formats** across a dataset assembled from multiple sources without normalizing first.
- **Ignoring context length limits**, leading to silent truncation of long examples.
- **Flattening multi-turn conversations** into a single text blob when the model expects structured turns.
- **Skipping a tokenizer dry-run** before committing to a full training job.

### Quick Reference Checklist

- [ ] Format matches your fine-tuning method (SFT, preference tuning, classification, pretraining, embeddings)
- [ ] File-level format (JSONL, Parquet, CSV, manifest) fits your training framework
- [ ] Target model's official chat template applied via tokenizer, not hand-written
- [ ] System prompt handling decided and kept consistent
- [ ] Multi-turn and multimodal records structured as expected, not flattened
- [ ] Sample examples rendered through the tokenizer and manually checked
- [ ] Token length distribution checked against the model's context window

### Example: Data Formats by Model Family and Purpose

| Model / Purpose | Recommended Format | Notes |
|---|---|---|
| Meta Llama 3 / 4 (instruction tuning) | JSONL, Llama chat template with header-tagged turns | Apply via tokenizer's chat template rather than hand-writing header tags |
| Alibaba Qwen (instruction tuning) | JSONL, ChatML-style turns | Widely reused format — many other community models also adopt ChatML |
| Mistral / Mixtral (instruction tuning) | JSONL, bracketed instruction-style turns | Historically limited or no native system-role support; check your specific version |
| Google Gemma (instruction tuning) | JSONL, turn-tagged chat template | Distinct start/end-of-turn tokens; apply through the official tokenizer template |
| OpenAI GPT models (API fine-tuning) | JSONL with role-based message arrays | One JSON object per line, each containing a full message list |
| Encoder / classification models (BERT-style) | CSV, TSV, or JSON with text + label columns | Label taxonomy should exactly match production categories |
| Preference tuning (DPO / RLHF-style, any base model) | JSONL with prompt, chosen, and rejected fields | Chosen/rejected responses should differ meaningfully, not just cosmetically |
| Code models (Code Llama, StarCoder, etc.) | JSONL instruction/code pairs, or raw source files | Some code models expect fill-in-the-middle formatted examples for completion tasks |
| Vision-language models (LLaVA-style, Qwen-VL, Gemma vision) | JSON/JSONL pairing image path or URL with conversation turns | Follow the specific model's multimodal schema — these vary more than text-only templates |
| Speech / ASR models (Whisper, NeMo-style) | JSONL manifest with audio path, transcript, and duration | One manifest line per audio clip; duration field is often required by the training framework |
| Embedding / retrieval fine-tuning | JSON or TSV anchor/positive/negative triplets | Hard negatives (similar but incorrect matches) improve retrieval quality more than random negatives |
| Continued pretraining / domain adaptation | Plain text or JSONL with a single "text" field | No roles or turns needed — just clean, representative text from the target domain |

Chat templates and schema conventions change with new model releases, so it's worth checking the specific model's tokenizer config or documentation before finalizing a large dataset in one format.

## step 4: Data normalization and cleaning 

### Cleaning and Normalizing Datasets for Fine-Tuning

Formatting gets a dataset into the right shape; cleaning makes sure what's inside that shape is actually worth training on. This guide is written to be run, not just read — each step includes working Python you can drop into a preprocessing script, using widely available libraries. Run the steps in order: each one assumes the previous one has already been applied.

### Why Cleaning Comes Before Everything Else

A perfectly formatted dataset full of encoding errors, duplicates, and mislabeled rows will still produce a poor fine-tune — formatting only controls the shape, not the substance. Cleaning should happen before you apply a chat template or split into train/eval, since template characters and turn structure can mask the underlying junk you're trying to catch.

### Setup

```bash
pip install ftfy langdetect datasketch pydantic pandas pyarrow datasets
```

Everything below reads into and writes out of a single pandas DataFrame called `df`. Whatever format your raw data or your downstream training framework expects, convert at the edges (Step 0 in, Export at the end) and keep the cleaning logic in between format-agnostic.

### Step 0: Load Your Dataset, Whatever Format It's In

Fine-tuning datasets show up as JSONL, CSV/TSV, Parquet, a single JSON array, or a Hugging Face Hub dataset. Normalize to a DataFrame immediately so every step after this one is format-agnostic.

```python
import pandas as pd
import json

# JSONL — one JSON object per line (most common for LLM fine-tuning data)
def load_jsonl(path: str) -> pd.DataFrame:
    with open(path, encoding="utf-8") as f:
        return pd.DataFrame([json.loads(line) for line in f if line.strip()])

df = load_jsonl("raw_data.jsonl")

# CSV / TSV
# df = pd.read_csv("raw_data.csv")            # add sep="\t" for TSV

# A single JSON file containing a top-level array of records
# df = pd.read_json("raw_data.json")

# Parquet
# df = pd.read_parquet("raw_data.parquet")

# A dataset hosted on the Hugging Face Hub
# from datasets import load_dataset
# df = load_dataset("org/dataset-name", split="train").to_pandas()
```

### Handling Different Schemas (Plain Text vs. Instruction vs. Chat Format)

The cleaning functions in Steps 1–9 are written against a single `text` column, but real fine-tuning data is often structured — `instruction`/`input`/`output` triples, or a `messages` list of `{role, content}` turns. Rather than rewriting every function per schema, build one `text` view for the checks that only need to *see* the content (dedup, language ID, length, quality), while applying the actual cleaning functions (`fix_text`, `normalize_whitespace`, `redact_pii`) to every real text field so the underlying data is what gets fixed, not just the view.

```python
def flatten_for_checks(row: dict) -> str:
    if "text" in row:
        return row["text"]
    if "messages" in row:
        return " ".join(m["content"] for m in row["messages"])
    if "instruction" in row:
        return " ".join(str(row.get(k, "")) for k in ("instruction", "input", "output"))
    return ""

def clean_fields(row: dict) -> dict:
    def clean(s: str) -> str:
        return normalize_whitespace(fix_text(redact_pii(s)))

    if "messages" in row:
        row["messages"] = [{**m, "content": clean(m["content"])} for m in row["messages"]]
    else:
        for field in ("text", "instruction", "input", "output"):
            if isinstance(row.get(field), str):
                row[field] = clean(row[field])
    return row

records = df.to_dict(orient="records")
records = [clean_fields(r) for r in records]
df = pd.DataFrame(records)
df["_flat_text"] = df.apply(lambda r: flatten_for_checks(r.to_dict()), axis=1)
# Steps 3–6 and Step 9 (dedup, language, quality, contamination) should run against
# df["_flat_text"] instead of df["text"] when your schema isn't plain text.
# Drop the helper column before exporting: df = df.drop(columns="_flat_text")
```

### Step 1: Fix Encoding and Unicode Issues

Scraped or merged text often contains mojibake (double-encoded text) and inconsistent Unicode forms (the same character represented multiple ways). Fix both before anything else, since later steps like deduplication and length filtering depend on clean text.

```python
import ftfy
import unicodedata

def fix_text(text: str) -> str:
    text = ftfy.fix_text(text)                # repairs mojibake / encoding artifacts
    text = unicodedata.normalize("NFKC", text) # collapses equivalent Unicode forms
    return text

df["text"] = df["text"].astype(str).map(fix_text)
```

### Step 2: Normalize Whitespace and Strip Structural Noise

Leftover HTML tags, markdown artifacts, and irregular whitespace are common in scraped or converted data and add noise without adding signal.

```python
import re

def normalize_whitespace(text: str) -> str:
    text = re.sub(r"<[^>]+>", " ", text)        # strip HTML tags
    text = re.sub(r"\s+", " ", text)             # collapse whitespace/newlines
    text = re.sub(r"[ \t]+\n", "\n", text)       # trim trailing spaces before newlines
    return text.strip()

df["text"] = df["text"].map(normalize_whitespace)
```

### Step 3: Remove Exact Duplicates

Exact duplicates are cheap to catch and easy to miss when a dataset is assembled from multiple overlapping sources.

```python
import hashlib

def row_hash(text: str) -> str:
    return hashlib.sha256(text.strip().lower().encode("utf-8")).hexdigest()

df["hash"] = df["text"].map(row_hash)
before = len(df)
df = df.drop_duplicates(subset="hash").drop(columns="hash")
print(f"Removed {before - len(df)} exact duplicates")
```

### Step 4: Remove Near-Duplicates

Near-duplicates (same content with minor edits) inflate a dataset's apparent size and can bias training toward whatever got duplicated most. MinHash-based locality-sensitive hashing (LSH) finds these efficiently at scale.

```python
from datasketch import MinHash, MinHashLSH

def get_minhash(text: str, num_perm: int = 128) -> MinHash:
    m = MinHash(num_perm=num_perm)
    for token in set(text.lower().split()):
        m.update(token.encode("utf-8"))
    return m

lsh = MinHashLSH(threshold=0.85, num_perm=128)  # 0.85 = similarity cutoff
keep_indices = []

for idx, text in enumerate(df["text"]):
    mh = get_minhash(text)
    if not lsh.query(mh):        # no sufficiently similar row seen yet
        lsh.insert(str(idx), mh)
        keep_indices.append(idx)

print(f"Removed {len(df) - len(keep_indices)} near-duplicates")
df = df.iloc[keep_indices].reset_index(drop=True)
```

### Step 5: Filter by Language

Mixed-language contamination is common in web-scraped and multi-source datasets. Filter to your target language(s) explicitly rather than assuming a source is monolingual.

```python
from langdetect import detect, DetectorFactory
DetectorFactory.seed = 0  # deterministic results

def safe_detect(text: str) -> str:
    try:
        return detect(text)
    except Exception:
        return "unknown"

df["lang"] = df["text"].map(safe_detect)
target_languages = {"en"}  # adjust to your task
before = len(df)
df = df[df["lang"].isin(target_languages)].drop(columns="lang")
print(f"Removed {before - len(df)} rows in other languages")
```

### Step 6: Filter by Length and Quality Heuristics

Very short rows are often junk (empty scrapes, placeholder text); very long rows risk truncation later. Simple heuristics catch a surprising amount of low-quality content cheaply.

```python
def quality_score(text: str) -> bool:
    if not (10 <= len(text.split()) <= 2000):
        return False
    alpha_ratio = sum(c.isalpha() for c in text) / max(len(text), 1)
    if alpha_ratio < 0.5:                     # mostly symbols/numbers = likely junk
        return False
    words = text.lower().split()
    if words and (max(words.count(w) for w in set(words)) / len(words)) > 0.3:
        return False                          # one word dominates = likely boilerplate/spam
    return True

before = len(df)
df = df[df["text"].map(quality_score)].reset_index(drop=True)
print(f"Removed {before - len(df)} low-quality rows")
```

### Step 7: Redact or Remove PII

Personally identifiable information can leak into training data from user-generated or scraped sources. A regex pass catches common patterns; for anything sensitive or regulated, use a dedicated PII-detection library (e.g., Microsoft Presidio) rather than relying on regex alone.

```python
PII_PATTERNS = {
    "email": re.compile(r"[\w.+-]+@[\w-]+\.[a-zA-Z]{2,}"),
    "phone": re.compile(r"\b\d{3}[-.\s]?\d{3}[-.\s]?\d{4}\b"),
    "ssn":   re.compile(r"\b\d{3}-\d{2}-\d{4}\b"),
}

def redact_pii(text: str) -> str:
    for label, pattern in PII_PATTERNS.items():
        text = pattern.sub(f"[{label.upper()}_REDACTED]", text)
    return text

df["text"] = df["text"].map(redact_pii)
```

### Step 8: Validate Schema and Labels

Malformed or mislabeled rows should be caught and logged, not allowed to crash a training job partway through. A schema model makes invalid rows explicit and easy to filter.

```python
from pydantic import BaseModel, ValidationError

class Example(BaseModel):
    text: str
    label: str

valid_rows, invalid_rows = [], []
for row in df.to_dict(orient="records"):
    try:
        valid_rows.append(Example(**row).dict())
    except ValidationError as e:
        invalid_rows.append((row, str(e)))

print(f"{len(valid_rows)} valid rows, {len(invalid_rows)} rejected")
df = pd.DataFrame(valid_rows)
```

### Step 9: Check for Contamination Against Your Evaluation Set

If any training row overlaps with your evaluation set, your eval results become unreliable. Check explicitly rather than trusting that your splits don't overlap.

```python
eval_hashes = set(eval_df["text"].map(row_hash))
before = len(df)
df = df[~df["text"].map(row_hash).isin(eval_hashes)].reset_index(drop=True)
print(f"Removed {before - len(df)} rows overlapping with the eval set")
```

### Step 10: Run Final Stats as a Sanity Check

Before handing the dataset off to formatting and training, print summary statistics so any remaining issue is visible before it costs you a training run.

```python
import numpy as np

lengths = df["text"].str.split().map(len)
print(f"Final row count: {len(df)}")
print(f"Length (words) — min: {lengths.min()}, median: {lengths.median()}, max: {lengths.max()}")
if "label" in df.columns:
    print("Label distribution:")
    print(df["label"].value_counts(normalize=True))
```

### Exporting the Cleaned Dataset

Your cleaning pipeline's output format doesn't have to match its input format — this is the natural point to convert into whatever the formatting stage (chat templates, JSONL for training, etc.) expects.

```python
if "_flat_text" in df.columns:
    df = df.drop(columns="_flat_text")

# JSONL — the standard handoff format for LLM fine-tuning
with open("cleaned_data.jsonl", "w", encoding="utf-8") as f:
    for record in df.to_dict(orient="records"):
        f.write(json.dumps(record, ensure_ascii=False) + "\n")

# CSV (fine for flat/classification schemas, awkward for nested "messages" fields)
# df.to_csv("cleaned_data.csv", index=False)

# Parquet — efficient for large datasets, preserves nested structures
# df.to_parquet("cleaned_data.parquet", index=False)

# Push straight to the Hugging Face Hub
# from datasets import Dataset
# Dataset.from_pandas(df).push_to_hub("your-username/cleaned-dataset")
```

### Putting It All Together: The Full Pipeline

Each step above is written to stand alone for debugging, but in practice you'll run them in sequence in a single script. This assumes every function from Steps 1–9 (`fix_text`, `normalize_whitespace`, `row_hash`, `get_minhash`, `safe_detect`, `quality_score`, `redact_pii`, `Example`) is already defined earlier in the file.

```python
def clean_pipeline(df: pd.DataFrame, eval_df: pd.DataFrame) -> pd.DataFrame:
    print(f"Start: {len(df)} rows")

    records = [clean_fields(r) for r in df.to_dict(orient="records")]   # Steps 1–2, 7
    df = pd.DataFrame(records)
    df["_flat_text"] = df.apply(lambda r: flatten_for_checks(r.to_dict()), axis=1)

    df["hash"] = df["_flat_text"].map(row_hash)                        # Step 3
    df = df.drop_duplicates(subset="hash").drop(columns="hash")
    print(f"After exact dedup: {len(df)} rows")

    lsh = MinHashLSH(threshold=0.85, num_perm=128)                     # Step 4
    keep = []
    for idx, text in enumerate(df["_flat_text"]):
        mh = get_minhash(text)
        if not lsh.query(mh):
            lsh.insert(str(idx), mh)
            keep.append(idx)
    df = df.iloc[keep].reset_index(drop=True)
    print(f"After near-dedup: {len(df)} rows")

    df["lang"] = df["_flat_text"].map(safe_detect)                     # Step 5
    df = df[df["lang"].isin({"en"})].drop(columns="lang")
    print(f"After language filter: {len(df)} rows")

    df = df[df["_flat_text"].map(quality_score)].reset_index(drop=True) # Step 6
    print(f"After quality filter: {len(df)} rows")

    valid = []                                                          # Step 8
    for row in df.to_dict(orient="records"):
        try:
            Example(**row)
            valid.append(row)
        except Exception:
            pass
    df = pd.DataFrame(valid)
    print(f"After schema validation: {len(df)} rows")

    eval_hashes = set(eval_df["text"].map(row_hash))                    # Step 9
    df = df[~df["_flat_text"].map(row_hash).isin(eval_hashes)].reset_index(drop=True)
    print(f"After contamination check: {len(df)} rows")

    return df.drop(columns="_flat_text")

df = load_jsonl("raw_data.jsonl")
eval_df = load_jsonl("eval_data.jsonl")
df = clean_pipeline(df, eval_df)

with open("cleaned_data.jsonl", "w", encoding="utf-8") as f:
    for record in df.to_dict(orient="records"):
        f.write(json.dumps(record, ensure_ascii=False) + "\n")
```

Run this on a small sample first (e.g., `df.sample(500)`) to confirm the pipeline behaves as expected before pointing it at a full dataset — schema mismatches and encoding surprises are far cheaper to catch on 500 rows than on 5 million.

### Common Mistakes to Avoid

- **Deduplicating after formatting** — chat template tokens make near-identical rows look different to a hash function; always dedupe on raw content first.
- **Relying only on exact-match dedup** — near-duplicates from re-scraped or lightly edited sources slip through and quietly bias training.
- **Regex-only PII scrubbing for sensitive domains** — regex catches common patterns but misses many real-world PII formats; use a dedicated tool when the stakes are high.
- **Skipping the contamination check** — this is the step most likely to make your eval numbers look better than reality.
- **Cleaning in place without logging what was removed** — you lose the ability to debug why row counts dropped or diagnose upstream data issues.

### Quick Reference Checklist

- [ ] Raw data loaded into a common working format (DataFrame) regardless of source file type
- [ ] Schema identified (plain text, instruction triples, or chat turns) and cleaning applied to every real text field, not just a flattened view
- [ ] Encoding and Unicode normalization applied
- [ ] Whitespace and structural noise (HTML, markdown artifacts) stripped
- [ ] Exact duplicates removed via content hashing
- [ ] Near-duplicates removed via MinHash/LSH or embedding similarity
- [ ] Rows filtered to target language(s)
- [ ] Length and quality heuristics applied
- [ ] PII redacted or removed
- [ ] Schema and labels validated, invalid rows logged and excluded
- [ ] Training set checked for overlap with the evaluation set
- [ ] Final summary stats reviewed before handing off to formatting
- [ ] Cleaned data exported to whatever format the next stage (formatting/training) expects, not left in the input format by default

## step 5: AI fine tunning methodology and finetunning optimization 


### Fine-Tuning Methodology and Optimization Guide

This is the pipeline's final and most important stage: actually running the fine-tune. Everything from the earlier guides in this series — picking a base model, sourcing data, formatting it, cleaning it — sets you up for this step. This guide covers the methods that exist (what), how to choose between them (when), and how to configure a run that minimizes training time and hardware requirements while keeping accuracy close to what full fine-tuning would get you (how).

### Part 1: What — The Fine-Tuning Methods

| Method | What It Updates | Relative Cost | Accuracy vs. Full Fine-Tuning | Best For |
|---|---|---|---|---|
| Full fine-tuning | Every model weight | Highest | 100% (the baseline others are measured against) | Small models, or when the domain shift is large enough that adapters can't capture it |
| LoRA | Small rank-decomposition matrices injected into attention/MLP layers (~0.1–1% of total parameters) | Low | ~90–95% | The default choice for most fine-tuning tasks |
| QLoRA | LoRA adapters on top of a 4-bit quantized frozen base model | Lowest (~75% less VRAM than LoRA) | ~80–90%, often closer to full LoRA parity in practice | VRAM-constrained or consumer hardware |
| DoRA | A LoRA variant that decomposes weight updates into magnitude and direction components | Similar to LoRA | Slightly above plain LoRA | The 2026 default upgrade over plain LoRA when your framework supports it |
| rsLoRA | Rank-stabilized scaling applied to LoRA | Similar to LoRA | Prevents the instability plain LoRA shows at high rank | Any time you go above rank 64 |
| Adapters / prefix / prompt tuning | Small inserted modules or learned virtual tokens | Very low | Generally below LoRA | Legacy setups or multi-task serving where swapping small modules matters more than peak quality |
| Continued pretraining | Most or all weights, trained on raw domain text (no instruction format) | High | Not directly comparable — different objective | Injecting new domain knowledge before instruction-tuning on top |
| SFT (supervised/instruction fine-tuning) | A training objective, layered on any of the above update methods | — | — | Teaching format, tone, and task behavior via labeled examples |
| Preference tuning (DPO, ORPO, KTO) | A training objective using paired or labeled preference data | Low–medium | Improves alignment/quality beyond what SFT alone achieves | Refining response quality after SFT, when you have preference data |
| GRPO (RL fine-tuning) | A training objective using on-policy sampling against a reward signal | Medium (far cheaper than classic PPO/RLHF pipelines) | Can push past SFT's imitation ceiling on the right tasks | Reasoning, math, and code tasks with a verifiable reward function |

The rows aren't all mutually exclusive — "method" (top rows) answers *what gets updated*, while "objective" (bottom rows) answers *what you're teaching*. A typical modern pipeline is QLoRA + DoRA as the update method, SFT as the first training objective, and optionally DPO or GRPO layered on afterward.

### Part 2: When — Choosing Your Approach

#### Fine-Tuning vs. Prompt Engineering vs. RAG

| Approach | Best When | Weak Point |
|---|---|---|
| Prompt engineering | You need to iterate fast with zero training cost | Doesn't scale to consistent style/format at high volume, and eats into your context budget every request |
| RAG | Your knowledge changes frequently, or you need citations/sources | Doesn't change *how* the model responds — only *what* it has access to — and adds retrieval latency and infrastructure |
| Fine-tuning | You need consistent tone/format, domain-specific terminology, or a smaller/faster model that doesn't need retrieval | Requires genuinely good training data (roughly 500+ examples as a practical floor, ideally 1,000+), and knowledge baked into weights goes stale as fast as the world does |

Fine-tuning teaches a model *how* to respond; RAG gives it *what* to say. In practice these are complementary rather than competing — many production systems fine-tune for behavior and pair it with RAG for facts.

#### Full Fine-Tuning vs. PEFT (LoRA/QLoRA/DoRA)

Default to a PEFT method (LoRA, QLoRA, or DoRA) unless one of these applies:

- The model is small enough (roughly under 3B parameters) that full fine-tuning is cheap anyway.
- You've already tuned LoRA rank and target modules and are still short of your accuracy bar.
- You're doing continued pretraining to inject a large volume of new domain knowledge rather than adapting existing behavior.

For nearly everything else, PEFT recovers 90%+ of full fine-tuning's quality at a small fraction of the compute — that trade is almost always worth taking.

#### LoRA vs. QLoRA vs. DoRA vs. rsLoRA

- **Enough VRAM for full-precision LoRA?** Use LoRA (or DoRA) directly — you get a slightly higher quality ceiling than quantized training.
- **VRAM-constrained or on consumer hardware?** Use QLoRA — the quality gap versus LoRA has narrowed significantly and is often small in practice.
- **Rank above 64?** Add rsLoRA to prevent the gradient instability plain LoRA shows at high rank.
- **Default starting point in 2026:** DoRA, rank 16, on all linear layers, with rsLoRA enabled if you raise the rank.

#### SFT vs. Preference Tuning vs. GRPO

- **Always start with SFT.** It establishes the base format and behavior you want, and nothing downstream works well without it.
- **Add DPO/ORPO/KTO after SFT** when responses are technically on-task but not in the *preferred* style, and you have (or can generate) paired preference data.
- **Reach for GRPO** when the task has a verifiable reward signal — a math answer that's checkably right or wrong, code that passes or fails tests, a tool call that succeeds or doesn't. It's the 2026-era default for RL-based fine-tuning because it needs far less compute than classic PPO/RLHF pipelines.

### Part 3: Hardware and Method Sizing Guide

Figures are ballpark and will shift with sequence length, batch size, and optimizer choice — treat them as a starting point for planning, not a guarantee.

| Model Size | Method | Approx. Min VRAM | Example Hardware | Notes |
|---|---|---|---|---|
| ≤3B | QLoRA | 4–6 GB | Laptop GPU / RTX 3060 | Runs on genuinely modest hardware |
| 7B–8B | QLoRA | 8–12 GB | RTX 4070 Ti | The most common "prosumer" fine-tuning target |
| 7B–8B | LoRA (no quantization) | 16–24 GB | RTX 4090 | Slightly higher ceiling than QLoRA if you have the VRAM |
| 13B–14B | QLoRA | 12–16 GB | RTX 4080 / 4090 | Still comfortably single-GPU |
| 30B–34B | QLoRA | ~24 GB | RTX 4090 (tight) or A6000 | Lean on gradient checkpointing and a small batch size |
| 70B | QLoRA | ~46 GB | Single A100 (80 GB) | The practical ceiling for single-GPU fine-tuning |
| 70B | LoRA (no quantization) | 140 GB+ | Multi-GPU (2–4× A100) | Consider DeepSpeed ZeRO-3 or FSDP for sharding |
| 100B+ / large MoE | QLoRA + multi-GPU | 48 GB+ per shard | Multi-GPU cluster | Only fine-tune at this scale if adapters on a hosted model genuinely can't do the job |

### Part 4: How — Step-by-Step Optimized Fine-Tuning

`donot run this yet look at step 6 before doing anything`

### Setup

```bash
pip install unsloth trl peft bitsandbytes accelerate transformers datasets
```

Unsloth is the practical 2026 default for single-GPU LoRA/QLoRA work — it rewrites the underlying training kernels for roughly 2x the speed and 60–70% less VRAM than a vanilla Transformers + PEFT setup, using the same training code you'd otherwise write.

### Step 1: Load and Quantize the Base Model

```python
from unsloth import FastLanguageModel

model, tokenizer = FastLanguageModel.from_pretrained(
    model_name="unsloth/llama-3.1-8b-bnb-4bit",  # pre-quantized checkpoint
    max_seq_length=2048,   # match this to your actual data length, not a round number
    load_in_4bit=True,     # this is what makes it QLoRA
    dtype=None,            # let Unsloth pick the best dtype for your GPU
)
```

### Step 2: Attach LoRA/DoRA Adapters

```python
model = FastLanguageModel.get_peft_model(
    model,
    r=16,                       # start here; raise only if you're underfitting
    lora_alpha=16,              # keep alpha == r as a sane default
    lora_dropout=0,             # 0 is fine and slightly faster for most fine-tunes
    target_modules=[
        "q_proj", "k_proj", "v_proj", "o_proj",
        "gate_proj", "up_proj", "down_proj",
    ],
    use_dora=True,              # 2026 default upgrade over plain LoRA
    use_rslora=False,           # flip to True if you raise r above ~64
    use_gradient_checkpointing="unsloth",  # framework-specific flag — don't skip this
)
```

### Step 3: Format Your Data

Apply the target model's chat template (see the dataset-formatting guide in this series for the full breakdown by model family):

```python
def format_example(example):
    return {"text": tokenizer.apply_chat_template(
        example["messages"], tokenize=False, add_generation_prompt=False
    )}

dataset = dataset.map(format_example)
```

### Step 4: Configure Training for Speed

```python
from trl import SFTTrainer, SFTConfig

args = SFTConfig(
    output_dir="out",
    per_device_train_batch_size=2,
    gradient_accumulation_steps=8,   # effective batch size = 2 * 8 = 16
    gradient_checkpointing=True,     # trades ~20-30% speed for large memory savings
    bf16=True,                       # use fp16=True instead only on older GPUs without bf16 support
    optim="paged_adamw_8bit",        # paged optimizer avoids memory spikes under QLoRA
    packing=True,                    # pack short examples together instead of padding
    max_seq_length=2048,             # match Step 1
    learning_rate=2e-4,
    lr_scheduler_type="cosine",
    warmup_ratio=0.03,
    num_train_epochs=3,
    logging_steps=10,
    eval_strategy="steps",
    eval_steps=50,
    load_best_model_at_end=True,
    metric_for_best_model="eval_loss",
)
```

### Step 5: Train

```python
trainer = SFTTrainer(
    model=model,
    tokenizer=tokenizer,
    train_dataset=dataset["train"],
    eval_dataset=dataset["validation"],
    args=args,
)
trainer.train()
```

### Step 6: Merge Adapters and Export

```python
# Merge LoRA weights into the base model for a standalone checkpoint
model.save_pretrained_merged("my-fine-tuned-model", tokenizer, save_method="merged_16bit")

# Or export straight to GGUF for local inference with llama.cpp / Ollama
model.save_pretrained_gguf("my-fine-tuned-gguf", tokenizer, quantization_method="q4_k_m")
```

### Step 7: Evaluate Against the Base Model

Run the same held-out prompts through both the base model and your fine-tune, side by side, in addition to whatever quantitative metric fits your task. Loss curves tell you training worked; a side-by-side comparison tells you it worked for your actual use case.

### Part 5: Optimization Techniques Cheat Sheet

| Technique | Reduces | Trade-off | When to Use |
|---|---|---|---|
| QLoRA (4-bit NF4 + double quantization) | VRAM, ~75% vs. LoRA | Small accuracy dip vs. full-precision LoRA | Whenever VRAM is the constraint |
| LoRA rank tuning | Trainable parameters / compute | Lower rank trains faster but caps capacity | Start at r=16, raise only if underfitting |
| DoRA | — (quality, not compute) | Marginally more compute than plain LoRA | Default upgrade unless your framework lacks support |
| rsLoRA | Gradient instability at high rank | None meaningful — it's a scaling fix | Any time rank exceeds ~64 |
| Gradient checkpointing | Activation memory | ~20–30% slower per step | Almost always worth it when memory-bound |
| Mixed precision (bf16) | Memory and step time | Effectively none on Ampere-generation GPUs or newer | Always, if your GPU supports bf16 |
| Flash Attention 2 | Memory and speed on long sequences | Requires a compatible GPU and library version | Long-context fine-tuning |
| Paged 8-bit optimizer | Optimizer state memory and memory spikes | None significant | QLoRA runs, especially with longer sequences |
| Sequence/sample packing | Wasted compute on padding | Slight risk of cross-example leakage if attention masking isn't handled correctly | Datasets with many short, variable-length examples |
| Gradient accumulation | Enables a larger effective batch on a small GPU | Slower wall-clock time per effective batch | Whenever your target batch size doesn't fit in VRAM |
| Smaller, higher-quality dataset | Total training steps and time | Requires more upfront curation effort (see the cleaning guide in this series) | Always prioritize over adding more raw, unfiltered data |
| Early stopping on validation loss | Wasted training steps and overfitting risk | Requires a genuine held-out validation set | Every run |
| Unsloth's custom kernels | Overall step time and VRAM (~2x speed, ~60–70% less VRAM) | Only covers supported model families; don't mix patched and unpatched HF utilities | The default choice for single-GPU LoRA/QLoRA work in 2026 |
| DeepSpeed ZeRO / FSDP | Memory across multiple GPUs via weight sharding | Added setup complexity and communication overhead | Full fine-tuning or very large models spanning multiple GPUs |
| GRPO instead of PPO/RLHF | VRAM for RL fine-tuning (~80% less than classic RLHF pipelines) | Needs a well-defined, verifiable reward function | RL fine-tuning on reasoning, math, or code tasks |

### Common Mistakes to Avoid

- **Setting LoRA rank high without enabling rsLoRA** — this is a well-documented cause of gradient instability, not a rare edge case.
- **Skipping the framework-specific gradient checkpointing flag** (e.g., `use_gradient_checkpointing="unsloth"`) — using the generic version loses much of the speedup.
- **Mixing a patched model (Unsloth) with non-patched Hugging Face utilities** — this causes silent miscompiles that are hard to debug after the fact.
- **Over-quantizing at export time** (e.g., a very aggressive GGUF quantization level) just to save a little disk space — this destroys quality you spent compute earning during training.
- **Jumping straight to full fine-tuning "to be safe"** — for the overwhelming majority of tasks, PEFT gets you 90%+ of the way there for a fraction of the cost.
- **Setting max sequence length far above your actual data length** — this wastes compute padding short examples for no benefit.
- **Training for too many epochs on a small dataset** — watch the validation loss; if it starts climbing while training loss keeps falling, you've moved from learning to memorizing.

### Quick Reference Checklist

- [ ] Confirmed fine-tuning — not prompt engineering or RAG — is the right tool for this task
- [ ] Chose a PEFT method (LoRA/QLoRA/DoRA) over full fine-tuning unless a specific reason says otherwise
- [ ] Hardware matched to model size and method using the sizing table
- [ ] LoRA config set with sane defaults (r=16, alpha=16, all major linear layers targeted, DoRA/rsLoRA as appropriate)
- [ ] Quantization (QLoRA) applied if VRAM-constrained
- [ ] Gradient checkpointing, mixed precision, and a paged optimizer all enabled
- [ ] Sequence length matched to actual data length, with packing enabled for short/variable-length examples
- [ ] Learning rate and cosine schedule with warmup configured
- [ ] Validation set held out and early stopping configured
- [ ] Adapters merged and exported to the format your deployment target needs
- [ ] Fine-tuned model evaluated head-to-head against the base model, not just judged by loss curves

## step 6: Using Cloud (kaggle) + huggies face repo

### Running the Pipeline on Kaggle and Publishing to Hugging Face

Everything else in this series — picking a base model, sourcing and cleaning data, formatting it, configuring an optimized fine-tune — needs somewhere to actually run. This guide covers doing all of it inside a free Kaggle notebook, end to end, and shipping the finished model to the Hugging Face Hub. It's written so each cell can be copy-pasted directly into a Kaggle notebook in order.

### What Kaggle Actually Gives You

| Resource | What You Get |
|---|---|
| GPU options | A single P100 (16 GB VRAM) or dual T4 (16 GB each, 32 GB combined but not automatically pooled — most single-GPU training scripts only see one T4 unless you explicitly use both) |
| TPU option | TPU v5e-8, useful mainly for JAX/TensorFlow workflows rather than typical PyTorch fine-tuning stacks |
| Weekly quota | Roughly 30 GPU-hours, resetting weekly |
| Session limit | Up to ~12 hours per GPU session (shorter for TPU sessions) |
| Storage | 20 GB persistent working directory per notebook |
| Cost | Free, no credit card required |

For context against the hardware sizing table from the optimization guide in this series: a single T4 or P100 comfortably handles **QLoRA fine-tuning on 7B–8B models**, which covers the large majority of practical fine-tuning projects. That's the sweet spot this guide is built around.

### Step 1: Verify Your Account and Enable GPU + Internet

Both GPU access and internet access on Kaggle require one-time phone verification.

1. Go to your notebook's **Settings** panel (right sidebar).
2. If you see a prompt like "Get phone verified," follow it — this unlocks both GPU and internet access.
3. Under **Accelerator**, select **GPU T4 x2** or **GPU P100** (either works for the 7B/8B QLoRA range).
4. Under **Internet**, switch it **On**. This is off by default and will restart your session when toggled — installing packages or pushing to Hugging Face will silently fail without it.

### Step 2: Create Your Notebook and Bring In Your Data

- **New Notebook** → give it a name.
- To attach a dataset you've already uploaded to Kaggle (e.g., your cleaned, formatted JSONL from the earlier steps in this series): click **Add Input** in the right sidebar, search for your dataset, and add it. It will appear read-only at `/kaggle/input/<your-dataset-name>/`.
- To pull a dataset directly from the Hugging Face Hub instead (with internet on):

```python
from datasets import load_dataset
dataset = load_dataset("org/dataset-name")
```

- If you haven't uploaded your own dataset to Kaggle yet: on the Kaggle homepage, go to **Datasets → New Dataset**, upload your cleaned JSONL file(s) from earlier in this pipeline, and set visibility to private if it's not meant to be public.

### Step 3: Install the Fine-Tuning Stack

Kaggle notebooks come with `torch` and `transformers` preinstalled, but usually not the fine-tuning-specific libraries — and sometimes at versions Unsloth wants to upgrade.

```python
!pip install -q unsloth trl peft bitsandbytes huggingface_hub
```

If you hit an import error right after installing (a common one: a package that got upgraded but the old version is still loaded in memory), restart the kernel (**Run → Restart Session**) and re-run from the top — don't just re-run the install cell.

### Step 4: Store Your Hugging Face Token Securely

Don't paste your Hugging Face token directly into a notebook cell — Kaggle notebooks are easy to accidentally make public or fork. Use Kaggle Secrets instead.

1. Get a **write**-scoped token from your Hugging Face account settings (Settings → Access Tokens).
2. In your Kaggle notebook, go to **Add-ons → Secrets**, add a new secret named `HF_TOKEN`, and paste your token as the value.
3. Retrieve it in code:

```python
from kaggle_secrets import UserSecretsClient
hf_token = UserSecretsClient().get_secret("HF_TOKEN")

from huggingface_hub import login
login(token=hf_token)
```

### Step 5: Run the Optimized Fine-Tuning Pipeline

This is the same configuration from the fine-tuning optimization guide in this series, adapted to Kaggle's paths and single-GPU environment.

```python
from unsloth import FastLanguageModel

model, tokenizer = FastLanguageModel.from_pretrained(
    model_name="unsloth/llama-3.1-8b-bnb-4bit",
    max_seq_length=2048,
    load_in_4bit=True,
    dtype=None,
)

model = FastLanguageModel.get_peft_model(
    model,
    r=16,
    lora_alpha=16,
    lora_dropout=0,
    target_modules=["q_proj", "k_proj", "v_proj", "o_proj", "gate_proj", "up_proj", "down_proj"],
    use_dora=True,
    use_gradient_checkpointing="unsloth",
)

def format_example(example):
    return {"text": tokenizer.apply_chat_template(
        example["messages"], tokenize=False, add_generation_prompt=False
    )}

dataset = load_dataset("json", data_files="/kaggle/input/your-dataset-name/cleaned_data.jsonl")
dataset = dataset.map(format_example)
```

### Step 6: Checkpoint Often — Kaggle Sessions Are Not Guaranteed to Finish

Kaggle sessions can hit the ~12-hour cap or your weekly 30-hour quota mid-run. Configure the trainer to save checkpoints regularly to `/kaggle/working` so a cut-off session doesn't cost you the whole run.

```python
from trl import SFTTrainer, SFTConfig

args = SFTConfig(
    output_dir="/kaggle/working/checkpoints",
    per_device_train_batch_size=2,
    gradient_accumulation_steps=8,
    gradient_checkpointing=True,
    bf16=True,
    optim="paged_adamw_8bit",
    packing=True,
    max_seq_length=2048,
    learning_rate=2e-4,
    lr_scheduler_type="cosine",
    warmup_ratio=0.03,
    num_train_epochs=3,
    save_strategy="steps",
    save_steps=50,
    save_total_limit=2,       # keep disk usage under the 20 GB working-directory limit
    logging_steps=10,
)

trainer = SFTTrainer(model=model, tokenizer=tokenizer, train_dataset=dataset["train"], args=args)
trainer.train(resume_from_checkpoint=True)  # safe even on a fresh run — it just finds nothing to resume
```

For a run you expect to take several hours, use **Save Version → Save & Run All (Commit)** rather than staying in the interactive editor. A committed run executes in the background as long as your GPU/internet settings are applied to it, and its output persists even if you close the browser tab — the interactive editor session does not survive a closed tab in the same way.

### Step 7: Merge Adapters and Export

```python
model.save_pretrained_merged("/kaggle/working/merged-model", tokenizer, save_method="merged_16bit")

# Optional: also export a quantized GGUF for local inference later
model.save_pretrained_gguf("/kaggle/working/gguf-model", tokenizer, quantization_method="q4_k_m")
```

Keep an eye on the 20 GB working-directory limit — a merged 8B model in fp16 is roughly 16 GB by itself, so exporting both the merged model and a GGUF version in the same run can push you close to the ceiling. If space is tight, push the merged model to the Hub first, then delete it locally before exporting the GGUF.

### Step 8: Push to the Hugging Face Hub

```python
model.push_to_hub("your-username/your-model-name", token=hf_token)
tokenizer.push_to_hub("your-username/your-model-name", token=hf_token)
```

For the GGUF file, or anything not natively handled by `push_to_hub`, upload it directly:

```python
from huggingface_hub import HfApi

api = HfApi()
api.upload_file(
    path_or_fileobj="/kaggle/working/gguf-model/model.Q4_K_M.gguf",
    path_in_repo="model.Q4_K_M.gguf",
    repo_id="your-username/your-model-name",
    token=hf_token,
)
```

Add a minimal model card so the repo isn't a bare set of files:

```python
readme = """---
license: apache-2.0
base_model: unsloth/llama-3.1-8b-bnb-4bit
tags:
- fine-tuned
- lora
---

# Your Model Name

Fine-tuned from Llama 3.1 8B using QLoRA + DoRA on a custom dataset.
Trained on a Kaggle T4/P100 GPU.
"""

api.upload_file(
    path_or_fileobj=readme.encode(),
    path_in_repo="README.md",
    repo_id="your-username/your-model-name",
    token=hf_token,
)
```

### Step 9: Verify the Upload

Before closing the notebook, pull the model back down fresh to confirm the push actually worked end to end.

```python
from transformers import AutoModelForCausalLM, AutoTokenizer

check_model = AutoModelForCausalLM.from_pretrained("your-username/your-model-name")
check_tokenizer = AutoTokenizer.from_pretrained("your-username/your-model-name")
print("Reload succeeded:", check_model.config.model_type)
```

### Kaggle-Specific Gotchas

- **Forgetting to toggle Internet on** — package installs and Hub pushes fail silently or with a confusing connectivity error; this is the single most common Kaggle-specific issue.
- **Not restarting the kernel after installing packages** — a newly installed version can conflict with one already loaded in memory; restart and re-run from the top rather than just re-running the install cell.
- **Assuming dual T4 means 32 GB for one model** — most single-process training scripts only use one of the two T4s unless you explicitly set up multi-GPU training; plan VRAM around a single 16 GB card unless you've configured otherwise.
- **Losing a run to the session cap** — always set `save_strategy="steps"` with a reasonable `save_steps`, and prefer committing long runs rather than relying on the interactive editor staying open.
- **Filling the 20 GB working directory** — merged fp16 models are large; push to the Hub and clean up local copies between export steps rather than keeping every format on disk at once.
- **Hardcoding your Hugging Face token in a cell** — use Kaggle Secrets; a forked or accidentally-public notebook with a hardcoded token is a real, common leak.
- **Pushing only the merged model when the adapter would do** — if you don't need a standalone merged checkpoint, pushing just the LoRA adapter (a few hundred MB instead of ~16 GB) is faster and cheaper to store, and can be loaded on top of the base model at inference time.

### Quick Reference Checklist

- [ ] Kaggle account phone-verified for GPU and internet access
- [ ] Accelerator set to GPU (T4 x2 or P100), Internet switched on
- [ ] Dataset attached via Kaggle Datasets or pulled directly from the Hugging Face Hub
- [ ] Fine-tuning stack installed, kernel restarted after install
- [ ] Hugging Face token stored in Kaggle Secrets, not hardcoded
- [ ] Training configured with frequent checkpointing (`save_strategy="steps"`)
- [ ] Long runs launched via Save & Run All (Commit) rather than left in the interactive editor
- [ ] Working directory usage watched against the 20 GB limit
- [ ] Adapters merged/exported, then pushed to the Hugging Face Hub with a model card
- [ ] Upload verified by reloading the model fresh from the Hub before closing the notebook

## Step 7: Running the model locally using ollama

### Running Your Fine-Tuned Model Locally with Ollama

This closes the loop on the series: the model was trained on Kaggle and pushed to Hugging Face, and this step gets it running entirely on your own machine, no cloud dependency, no API cost. Since the earlier Kaggle guide already exported a GGUF file and pushed it to your Hugging Face repo, most of the work here is already done — Ollama can pull it directly.

### Step 1: Install Ollama

```bash
# macOS
brew install ollama

# Linux
curl -fsSL https://ollama.com/install.sh | sh

# Windows
# Download the installer from https://ollama.com/download
```

Verify it's working:

```bash
ollama --version
```

Ollama runs a local background service (`ollama serve`) that most installers start automatically. If a command says it can't connect, run `ollama serve` in a separate terminal first.

### Step 2: The Fast Path — Run Directly From Your Hugging Face Repo

If you followed the Kaggle guide in this series, your GGUF file is already sitting in your Hugging Face repo. Ollama can pull and run it with one command, no local download or Modelfile needed:

```bash
ollama run hf.co/your-username/your-model-name
```

By default this pulls the `Q4_K_M` quantization if your repo has one (which it does, if you followed the earlier export step). To pick a different quantization level from the same repo:

```bash
ollama run hf.co/your-username/your-model-name:Q5_K_M
```

**If your repo is private**, set your Hugging Face token as an environment variable before running so Ollama can authenticate:

```bash
export HF_TOKEN=your_token_here
ollama run hf.co/your-username/your-model-name
```

This is the easiest path and the one worth trying first — GGUF files exported the way the earlier guide did (via `save_pretrained_gguf`) carry their chat template as embedded metadata, so Ollama typically applies the correct template automatically.

### Step 3: The Manual Path — Build From a Local GGUF File

Use this if you have a GGUF file locally instead of on the Hub, or if you want full control over the template, stop tokens, or a baked-in system prompt.

```bash
# If you don't already have it locally, download it from your Hugging Face repo:
huggingface-cli download your-username/your-model-name model.Q4_K_M.gguf --local-dir ./model
```

Create a file named `Modelfile` (no extension) in the same directory:

```
FROM ./model/model.Q4_K_M.gguf

TEMPLATE """{{ if .System }}<|start_header_id|>system<|end_header_id|>

{{ .System }}<|eot_id|>{{ end }}{{ if .Prompt }}<|start_header_id|>user<|end_header_id|>

{{ .Prompt }}<|eot_id|>{{ end }}<|start_header_id|>assistant<|end_header_id|>

{{ .Response }}<|eot_id|>"""

PARAMETER stop "<|start_header_id|>"
PARAMETER stop "<|end_header_id|>"
PARAMETER stop "<|eot_id|>"
PARAMETER temperature 0.7

SYSTEM """You are a helpful assistant fine-tuned for [your task]."""
```

Then build and run it:

```bash
ollama create my-model -f Modelfile
ollama run my-model
```

### Step 4: Get the Template Right (This Is the Part That Breaks)

The `TEMPLATE` block in a Modelfile must match the chat template your model was actually fine-tuned with — the one covered in the dataset-formatting guide earlier in this series. A mismatched template is the single most common cause of a locally-running model producing garbled, repetitive, or oddly-truncated output despite training having gone fine.

- The example above matches Llama 3.x's header-tag format.
- For a Qwen-family fine-tune, swap in the ChatML tags (`<|im_start|>`/`<|im_end|>`) instead.
- For Mistral or Gemma, use their respective bracketed or turn-tagged formats.

If you went through Step 2 (pulling directly from your Hugging Face repo), you generally don't need to write this yourself — Ollama reads the template embedded in the GGUF's metadata. Step 3's manual Modelfile is mainly for when that metadata is missing, wrong, or you deliberately want to override it (e.g., baking in a fixed system prompt).

### Step 5: Chat With It and Compare Against the Base Model

```bash
# Interactive chat
ollama run my-model

# One-off prompt, non-interactive
ollama run my-model "Summarize this in two sentences: ..."
```

Run the same prompts through the unmodified base model (`ollama run llama3.1:8b`, for instance) and compare side by side — this is the same head-to-head check recommended at the end of the fine-tuning optimization guide, just now runnable locally instead of in a notebook.

### Step 6: Use It From Code via the Local API

Ollama exposes a REST API on `localhost:11434`, so you can call your fine-tuned model from a script without going through the CLI:

```python
import requests

response = requests.post("http://localhost:11434/api/chat", json={
    "model": "my-model",
    "messages": [{"role": "user", "content": "Hello!"}],
    "stream": False,
})
print(response.json()["message"]["content"])
```

This is the same interface used to build local apps, VS Code extensions, or anything else you want talking to your model without hitting a cloud API.

### Step 7: Managing Models Locally

```bash
ollama list                        # see everything installed locally
ollama show my-model --modelfile   # inspect the Modelfile Ollama actually used, including any auto-detected template
ollama cp my-model my-model-backup # duplicate under a new name
ollama rm my-model                 # free up disk space
```

### Choosing a Quantization Level (Disk and RAM Trade-off)

| Quantization | Approx. Size (8B model) | Approx. RAM Needed | Notes |
|---|---|---|---|
| Q8_0 | ~8.5 GB | ~10 GB+ | Closest to full precision; use if you have the RAM and want minimal quality loss |
| Q5_K_M | ~5.5 GB | ~7 GB+ | Good middle ground for most everyday use |
| Q4_K_M | ~4.5–5 GB | ~6 GB+ | The default and generally the best speed/quality trade-off |
| Q3_K_M | ~3.5 GB | ~5 GB+ | Noticeable quality drop; only for tight hardware |
| Q2_K | ~3 GB | ~4 GB+ | Quality degrades enough that it's rarely worth it — the same warning from the optimization guide applies here |

As a rule of thumb, plan for total system RAM around 1.3–1.5x the GGUF file size to leave headroom for context and the OS.

### Common Mistakes to Avoid

- **Assuming any GGUF "just works"** — a mismatched chat template is the most common cause of broken-looking output from a technically-successful fine-tune.
- **Forgetting stop tokens in a manual Modelfile** — without them, the model doesn't know when to stop generating and will run on past where it should end a turn.
- **Grabbing the most aggressive quantization to save space** — Q2_K in particular tends to noticeably hurt quality; Q4_K_M is the sane default for most use cases.
- **Not enough system RAM for the chosen quantization** — this causes heavy swapping and very slow generation rather than a clean out-of-memory error, which can look like a bug when it's actually a sizing issue.
- **Pointing `FROM` at the wrong relative path** in a manual Modelfile — `ollama create` fails immediately if the GGUF path is wrong, so check the path before assuming something deeper is broken.
- **Trying to pull a private Hugging Face repo without setting `HF_TOKEN`** — this fails silently as an auth/not-found error rather than a clear permissions message.

### Quick Reference Checklist

- [ ] Ollama installed and `ollama serve` running
- [ ] Tried the direct Hugging Face path (`ollama run hf.co/username/repo`) before building a manual Modelfile
- [ ] For private repos, `HF_TOKEN` set as an environment variable
- [ ] If building manually, `TEMPLATE` matches the exact chat format used during fine-tuning
- [ ] Stop tokens included in the Modelfile if built manually
- [ ] Quantization level chosen based on available RAM, not just smallest file size
- [ ] Output compared side by side against the base model to confirm the fine-tune is actually doing its job
- [ ] Local REST API tested if you plan to call the model from code rather than the CLI

## Step 8: Benchmarking the generate model 


### Benchmarking Your Fine-Tuned Model

Training converging and the model running locally doesn't tell you whether the fine-tune actually worked. This guide covers the methods that exist for measuring that (what), how to pick between them based on your model type, purpose, and stage in the process (when), and how to actually run each one (how).

### Part 1: What — The Benchmarking Methods

| Method | What It Measures | Format | Tooling |
|---|---|---|---|
| Standard academic benchmarks (MMLU, HellaSwag, TruthfulQA, GSM8K, BBH, HumanEval) | General knowledge and reasoning capability | Fixed question sets with known answers | `lm-evaluation-harness` |
| Task-specific automatic metrics (accuracy, F1, exact match, BLEU, ROUGE, pass@k) | Correctness against a reference answer for your specific task | Reference-based scoring | Hugging Face `evaluate` library |
| Perplexity / held-out loss | How well the model fits held-out data, independent of task framing | A single numeric score | Custom script, a few lines of PyTorch |
| LLM-as-judge | Open-ended quality — helpfulness, tone, instruction-following, style | Pairwise comparison or rubric scoring by a strong external model | Custom API script, or a framework like DeepEval/MLflow |
| Human evaluation | Ground-truth judgment, especially for subjective or high-stakes quality | Blind side-by-side rating or Likert scale | Spreadsheet or lightweight review tool |
| Regression / behavior-preservation testing | Whether fine-tuning broke something the base model used to do well | Before/after comparison on general benchmarks | `lm-evaluation-harness` run twice |
| Safety and alignment evaluation | Toxicity, bias, refusal correctness | Rubric-based scoring, often LLM-as-judge | Custom rubric or DeepEval's safety metrics |
| Efficiency benchmarking (latency, throughput, memory) | Whether the optimization work actually paid off in real terms | Timed inference runs | Custom timing script against your deployment target |

### Part 2: When — Choosing Methods

#### By Model Type

| Model Type | Primary Methods |
|---|---|
| Chat / instruction-following | LLM-as-judge for open-ended quality, general benchmarks for regression |
| Classification / encoder-based | Accuracy, F1, precision/recall, confusion matrix |
| Code generation | pass@k against executable test cases (HumanEval/MBPP-style) |
| Summarization / translation | ROUGE/BLEU/BERTScore, paired with human or LLM-judge review for fluency |
| Embedding / retrieval | Recall@k, MRR, nDCG against a labeled retrieval set |
| RL/GRPO-tuned reasoning models | Task-specific verifiable accuracy — is the math answer right, does the code pass its tests |

#### By Purpose

| Purpose | Recommended Depth |
|---|---|
| Personal / hobby project | A quick LLM-as-judge head-to-head on 20–30 examples plus a manual spot check — a full academic suite is overkill |
| Production deployment | Regression check + task-specific metrics + LLM-as-judge + a human review pass before launch, plus ongoing monitoring afterward |
| Research / publication | Full benchmark suite via `lm-evaluation-harness` for comparability with published numbers, run under the harness's default prompt templates so scores aren't apples-to-oranges |
| Safety-critical (medical, legal, financial) | Human review is not optional here — automated methods can de-risk iteration, but a person needs to sign off before anything ships |

#### By Process Stage

| Stage | What to Run |
|---|---|
| During training (checkpoint selection) | Validation loss / perplexity — cheap enough to compute every few hundred steps |
| Immediately after training | A regression check against general benchmarks, plus your task-specific metric |
| Before deployment | LLM-as-judge head-to-head, a human spot-check, and an efficiency benchmark on your actual deployment target |
| Post-deployment | Ongoing sampled review (human or LLM-judge) of live outputs, watching for drift over time |

### Part 3: How — Running Each Method

#### Standard Benchmarks and Regression Checks

```bash
pip install lm-eval
```

```bash
lm_eval --model hf \
    --model_args pretrained=your-username/your-model-name \
    --tasks mmlu,hellaswag,gsm8k,truthfulqa_mc2 \
    --device cuda:0 \
    --batch_size 8
```

Run the identical command against the unmodified base model and compare. As a rule of thumb, a drop of more than roughly 2–3 points on a general benchmark like MMLU signals that fine-tuning has degraded general knowledge — a sign of catastrophic forgetting worth investigating (lower the learning rate, reduce epochs, or check whether your dataset is too narrow). CPU evaluation works here but is slow — plan for hours even on a small model, so use a GPU session if you have one available from earlier in this pipeline.

#### Task-Specific Automatic Metrics

```python
import evaluate

rouge = evaluate.load("rouge")
results = rouge.compute(predictions=predictions, references=references)

# Classification
f1 = evaluate.load("f1")
f1_score = f1.compute(predictions=predictions, references=labels, average="macro")

# Extractive QA
squad_metric = evaluate.load("squad")
em_f1 = squad_metric.compute(predictions=predictions, references=references)
```

50 examples give a rough estimate; 200+ gives a stable, trustworthy result for these reference-based metrics.

#### Perplexity on a Held-Out Set

```python
import torch

def compute_perplexity(model, tokenizer, texts):
    losses = []
    for text in texts:
        inputs = tokenizer(text, return_tensors="pt", truncation=True, max_length=2048).to(model.device)
        with torch.no_grad():
            loss = model(**inputs, labels=inputs["input_ids"]).loss
        losses.append(loss.item())
    return torch.exp(torch.tensor(losses).mean()).item()
```

Lower is better, but treat perplexity as a sanity check, not a quality score — a model can have great perplexity on your data and still produce responses that miss what the task actually needs.

#### LLM-as-Judge

```python
import json
from anthropic import Anthropic

client = Anthropic()

JUDGE_PROMPT = """You are evaluating two AI responses to the same prompt.

Prompt: {prompt}

Response A: {response_a}
Response B: {response_b}

First, analyze each response against these criteria: accuracy, helpfulness, and adherence to the expected format.
Then output a JSON object: {{"reasoning": "...", "winner": "A" or "B" or "tie"}}"""

def judge_pair(prompt, response_a, response_b):
    message = client.messages.create(
        model="claude-sonnet-4-6",
        max_tokens=500,
        messages=[{"role": "user", "content": JUDGE_PROMPT.format(
            prompt=prompt, response_a=response_a, response_b=response_b
        )}],
    )
    return json.loads(message.content[0].text)
```

A few practices that matter here more than the code itself:

- **Ask for reasoning before the verdict** (as above) — this measurably reduces arbitrary scoring and gives you something to inspect when a result looks off.
- **Randomize which response is A and which is B** across your test set — judges show a measurable position bias toward whichever slot comes first.
- **Include 2–4 calibration examples** in the prompt (an ideal response and a flawed one, each with a rationale) if your criteria are at all subjective.
- **Cross-check with a small human sample.** 30 examples is usually enough for a directional signal — don't trust the judge blindly, especially early on.
- **Use a judge from a different provider than any model involved in generating your training data**, to avoid self-preference bias skewing results in the fine-tune's favor.

#### Human Evaluation

You don't need a platform for this. A spreadsheet works:

1. Sample 20–50 prompts covering both typical and edge cases.
2. Generate a response from both the base and fine-tuned model for each.
3. Present them **blind and randomized** (reviewer shouldn't know which is which) in two columns.
4. Have one or more reviewers mark a preference or a 1–5 quality score per response.
5. Tally win rate or mean score difference.

This is the step that's easy to skip and the one most worth keeping — automated methods de-risk iteration, but a human pass before anything ships is what actually catches the failures that don't show up in a metric.

#### Efficiency Benchmarking

This is the step that verifies the optimization work from earlier in this series actually paid off in real terms, not just in theory.

```python
import time

def benchmark_inference(model_call, prompt, num_runs=10):
    times, token_counts = [], []
    for _ in range(num_runs):
        start = time.time()
        response = model_call(prompt)
        elapsed = time.time() - start
        times.append(elapsed)
        token_counts.append(len(response.split()))  # swap for a real tokenizer count if precision matters

    avg_time = sum(times) / len(times)
    avg_tokens_per_sec = sum(token_counts) / sum(times)
    print(f"Avg latency: {avg_time:.2f}s | Avg throughput: {avg_tokens_per_sec:.1f} tokens/sec")
```

Run this against each quantization level you're considering (see the quantization table in the Ollama guide) to see the actual speed/quality trade-off on your own hardware rather than assuming it from general figures — real numbers on your machine beat published estimates every time.

### Part 4: A Minimal Complete Benchmark Suite

If you don't want to reason through every option above, this default sequence covers the large majority of practical fine-tuning projects:

1. Run `lm-evaluation-harness` on 2–3 general benchmarks relevant to your domain, on both the base and fine-tuned model, as a regression check.
2. Run your own held-out task-specific eval set through the metric that matches your task type (accuracy, F1, ROUGE, pass@k, etc.).
3. Run an LLM-as-judge pairwise comparison, base vs. fine-tuned, on 30–50 examples.
4. Benchmark inference speed and memory on your actual deployment target (e.g., your local Ollama setup).
5. Do a manual human spot-check of ~20 examples, weighted toward edge cases rather than the easy majority.

### Common Mistakes to Avoid

- **Watching only the training loss curve** and treating a converged loss as proof the model is good — loss tells you training worked, not that the model does what you need.
- **Trusting LLM-as-judge with no human cross-check**, or letting a model judge outputs from a system that helped generate its own training data — both invite bias you won't notice without a sanity check.
- **Skipping the regression check** against the base model — this is exactly how catastrophic forgetting slips into production unnoticed.
- **Comparing your scores directly against published leaderboard numbers** when your prompt template or answer-extraction logic differs from how the benchmark's authors ran it — the numbers aren't comparable unless the evaluation methodology matches exactly.
- **Skipping human review before shipping anything safety-critical** — automated methods are for de-risking iteration speed, not for replacing a person's judgment on the final call.
- **Evaluating only on data that looks like your training distribution** and never on genuinely out-of-domain or adversarial inputs.
- **Doing all the optimization work from earlier in this series and never actually measuring the speed gain** — verify it with real timing numbers on your deployment target instead of assuming the theory held.

### Quick Reference Checklist

- [ ] Benchmark methods chosen based on model type, purpose, and process stage — not just whatever's easiest to run
- [ ] Regression check run against the base model on general benchmarks
- [ ] Task-specific automatic metric computed on a held-out set matching your actual use case
- [ ] LLM-as-judge comparison run with randomized ordering and reasoning-before-verdict prompting
- [ ] LLM-as-judge results cross-checked against a small human sample
- [ ] Human review pass completed before anything ships, especially for subjective or high-stakes tasks
- [ ] Efficiency (latency/throughput/memory) benchmarked on the real deployment target, not assumed from theory
- [ ] Results compared like-for-like — same prompt template and evaluation logic as any published numbers you're comparing against

## step 9: Embedding the model into the application 

# Deploying a Trained AI Model into Applications

A practical, no-fluff guide to taking a trained model and getting it running inside real applications — APIs, Docker containers, mobile, browser, and serverless.

## The Core Idea

Once a model is trained, "embedding" it usually means one of these patterns:

- **Wrap it in an API** — model runs on a server, apps talk to it over HTTP/gRPC.
- **Bundle it inside the app itself** — model ships with the binary (mobile, desktop, browser).
- **Run it on the edge** — model runs on-device (IoT, embedded hardware).
- **Serverless inference** — model runs only when triggered (Lambda, Cloud Run).

Pick based on latency needs, privacy constraints, model size, and how often it needs updating.

## Step 1: Export the Model in a Portable Format

Don't ship your raw training checkpoint. Export it to a format built for inference.

- PyTorch → `TorchScript` (`.pt`) or `ONNX` (`.onnx`)
- TensorFlow/Keras → `SavedModel` or `.tflite` (mobile/edge)
- scikit-learn → `joblib` / `pickle` (fine for small classical models)
- Any framework → `ONNX` is the best universal choice if you need cross-platform support

```python
# PyTorch -> ONNX example
import torch
model.eval()
dummy_input = torch.randn(1, 3, 224, 224)
torch.onnx.export(model, dummy_input, "model.onnx", opset_version=17)
```

## Step 2: Wrap It in an Inference API

Most applications talk to a model through a small web service, not the model directly.

```python
# app.py - FastAPI example
from fastapi import FastAPI
import onnxruntime as ort
import numpy as np

app = FastAPI()
session = ort.InferenceSession("model.onnx")

@app.post("/predict")
def predict(data: dict):
    input_array = np.array(data["input"], dtype=np.float32)
    outputs = session.run(None, {"input": input_array})
    return {"prediction": outputs[0].tolist()}
```

Run it locally first:

```bash
pip install fastapi uvicorn onnxruntime numpy
uvicorn app:app --host 0.0.0.0 --port 8000
```

This is now a standalone service any application (web, mobile, desktop) can call over HTTP.

## Step 3: Containerize It with Docker

This is the part that makes it portable across machines, teams, and cloud providers.

```dockerfile
# Dockerfile
FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py model.onnx ./

EXPOSE 8000
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8000"]
```

```text
# requirements.txt
fastapi
uvicorn
onnxruntime
numpy
```

Build and run:

```bash
docker build -t my-model-api .
docker run -p 8000:8000 my-model-api
```

Now any application — a web frontend, a mobile app, another backend service — just calls `http://<host>:8000/predict`.

## Step 4: Multi-Container Setups (Model + App)

If your model service backs a full application (frontend + backend + model), use `docker-compose` so everything spins up together.

```yaml
# docker-compose.yml
version: "3.9"
services:
  model-api:
    build: ./model-service
    ports:
      - "8000:8000"
  backend:
    build: ./backend
    ports:
      - "5000:5000"
    depends_on:
      - model-api
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend
```

```bash
docker-compose up --build
```

The backend calls `model-api` internally by service name (`http://model-api:8000/predict`), no need to expose it publicly.

## Step 5: Choosing an Integration Pattern per App Type

Different applications need the model embedded differently:

**Web application (React/Vue/etc. + backend)**
- Model runs as its own Dockerized microservice.
- Frontend never talks to the model directly — always through your backend, for auth and rate limiting.

**Mobile application (iOS/Android)**
- Convert to `.tflite` (Android) or `Core ML` (iOS) and bundle inside the app for offline/low-latency use.
- Or call a hosted API if the model is too large or updates frequently.

**Browser-only (no backend)**
- Use `ONNX Runtime Web` or `TensorFlow.js` to run inference client-side in JavaScript.
- Good for privacy-sensitive or offline-first apps; limited by model size and browser compute.

**Desktop application**
- Embed via `ONNX Runtime` C++/Python/C# bindings directly in the app process.
- Avoids network calls entirely; good for local tools.

**Serverless (AWS Lambda, Google Cloud Run, Azure Functions)**
- Package the same Docker image and deploy as a container-based function.
- Best for spiky, infrequent traffic — you don't pay for idle compute.
- Watch cold-start time if the model is large.

**Edge/IoT devices**
- Use quantized models (`.tflite`, `ONNX quantized`, `NVIDIA TensorRT`) to fit memory/compute limits.
- Runs fully on-device, no network dependency.

## Step 6: Production Concerns (Don't Skip These)

- **Batching**: batch incoming requests where possible to use GPU/CPU efficiently.
- **Model versioning**: tag Docker images with model version (`my-model-api:v1.2`), never `latest` in production.
- **Health checks**: add a `/health` endpoint so orchestrators (Kubernetes, ECS) know when the container is ready.
- **Monitoring**: log prediction latency, input shapes, and error rates — models fail silently more often than code does.
- **GPU support**: if you need GPU inference in Docker, use an `nvidia/cuda` base image and run with `--gpus all` (requires NVIDIA Container Toolkit on the host).
- **Scaling**: for real traffic, put the container behind Kubernetes or a managed service (SageMaker, Vertex AI, Azure ML) rather than a single `docker run`.

## Minimal Decision Cheat Sheet

- Need it callable from anywhere → **Dockerized API**
- Need it offline on a phone → **TFLite / Core ML bundled in app**
- Need it in-browser, no server → **ONNX Runtime Web / TensorFlow.js**
- Need it cheap for occasional use → **Serverless container**
- Need it on a sensor/microcontroller → **Quantized model, on-device runtime**
