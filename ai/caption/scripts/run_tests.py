import os
import json
import time
import pandas as pd
import random

from pathlib import Path
from dotenv import load_dotenv
from google import genai
from google.genai import errors, types


# ==========================================
# 1. INITIAL CONFIGURATION
# ==========================================

BASE_DIR = Path(__file__).resolve().parent.parent

load_dotenv(BASE_DIR / ".env")

API_KEY = os.getenv("GEMINI_API_KEY")

# Gunakan ID model yang sama dengan konfigurasi
# model yang tersedia di Google AI Studio.
MODEL_NAME = "gemini-3.5-flash-lite"

PROMPT_PATH = BASE_DIR / "prompts" / "system_instructions_v2.md"
DATASET_PATH = BASE_DIR / "dataset" / "test_cases.json"
OUTPUT_PATH = BASE_DIR / "evaluation" / "model_testing_results.xlsx"

if not API_KEY:
    raise ValueError("GEMINI_API_KEY belum ditemukan di file .env")

client = genai.Client(api_key=API_KEY)


RETRYABLE_STATUS_CODES = {429, 500, 502, 503, 504}
def call_with_retry(fn, max_retries=6, base_delay=6, max_delay=60):
    for attempt in range(max_retries):
        try:
            return fn()
        except errors.APIError as e:
            code = getattr(e, "code", None)
            if code in RETRYABLE_STATUS_CODES and attempt < max_retries - 1:
                # Exponential backoff + jitter agar tidak menabrak antrean serentak
                exp_delay = min(max_delay, base_delay * (2**attempt))
                sleep_time = exp_delay + random.uniform(0.5, 2.5)
                print(
                    f"  [Retry {attempt + 1}/{max_retries - 1}] "
                    f"HTTP {code} (high demand/rate limit). Menunggu {sleep_time:.1f}s..."
                )
                time.sleep(sleep_time)
                continue
            raise


# ==========================================
# 2. LOAD SYSTEM INSTRUCTIONS & DATASET
# ==========================================

with open(PROMPT_PATH, "r", encoding="utf-8") as file:
    system_instruction = file.read()

with open(DATASET_PATH, "r", encoding="utf-8") as file:
    test_cases = json.load(file)

print(f"Total test cases: {len(test_cases)}")
print(f"Model: {MODEL_NAME}")
print("=" * 50)


# ==========================================
# 3. PREPARE TESTING
# ==========================================

results = []

for index, test in enumerate(test_cases, start=1):

    test_id = test["id"]
    scenario = test["scenario"]
    product_input = test["input"]
    expected = test["expected_behavior"]

    print(f"[{index}/{len(test_cases)}] Testing {test_id}...")

    try:

        # Send product input to Gemini with retry on rate limit (429)
        response = call_with_retry(
            lambda: client.models.generate_content(
                model=MODEL_NAME,
                contents=json.dumps(
                    product_input,
                    ensure_ascii=False
                ),
                config=types.GenerateContentConfig(
                    system_instruction=system_instruction,
                    temperature=1.0,
                    max_output_tokens=2048,
                    response_mime_type="application/json",
                    thinking_config=types.ThinkingConfig(
                        thinking_level="MEDIUM"
                ),
                )
            )
        )

        raw_output = response.text or ""

        # Check JSON validity
        try:
            parsed_output = json.loads(raw_output)
            json_valid = True

        except (json.JSONDecodeError, TypeError):
            parsed_output = None
            json_valid = False

        # Check output structure
        if json_valid:
            description = parsed_output.get("description")
            hashtags = parsed_output.get("hashtags")

            structure_valid = (
                isinstance(description, str)
                and isinstance(hashtags, list)
                and all(isinstance(tag, str) for tag in hashtags)
                and len(hashtags) >= 3
                and len(hashtags) <= 7
                and set(parsed_output.keys()) == {
                    "description", "hashtags"
                }
            )

        else:
            description = None
            hashtags = None
            structure_valid = False

        # Save test result
        results.append({
            "Test ID": test_id,
            "Scenario": scenario,
            "Input": json.dumps(
                product_input,
                ensure_ascii=False
            ),
            "Expected Behavior": expected,
            "AI Output": raw_output,
            "Description": description,
            "Hashtags": (
                ", ".join(hashtags)
                if isinstance(hashtags, list)
                else None
            ),
            "JSON Valid": "PASS" if json_valid else "FAIL",
            "Structure Valid": (
                "PASS" if structure_valid else "FAIL"
            ),
            "Description Quality": "",
            "Factual Accuracy": "",
            "Hashtag Relevance": "",
            "Evaluation Notes": "",
            "Error": ""
        })

        print(
            f"  JSON: {'PASS' if json_valid else 'FAIL'} | "
            f"Structure: {'PASS' if structure_valid else 'FAIL'}"
        )

    except Exception as error:

        print(f"  ERROR: {error}")

        results.append({
            "Test ID": test_id,
            "Scenario": scenario,
            "Input": json.dumps(product_input, ensure_ascii=False),
            "Expected Behavior": expected,
            "AI Output": "",
            "Description": "",
            "Hashtags": "",
            "JSON Valid": "ERROR",
            "Structure Valid": "ERROR",
            "Description Quality": "",
            "Factual Accuracy": "",
            "Hashtag Relevance": "",
            "Evaluation Notes": "",
            "Error": str(error)
        })

    # Avoid sending requests too quickly (~12 requests/minute)
    if index < len(test_cases):
        time.sleep(5)


# ==========================================
# 4. EXPORT RESULTS TO EXCEL
# ==========================================

df = pd.DataFrame(results)

OUTPUT_PATH.parent.mkdir(
    parents=True,
    exist_ok=True
)

df.to_excel(OUTPUT_PATH, index=False)

# Summary
total = len(results)

json_pass = sum(
    row["JSON Valid"] == "PASS"
    for row in results
)

structure_pass = sum(
    row["Structure Valid"] == "PASS"
    for row in results
)

print("\n" + "=" * 50)
print("TESTING COMPLETED")
print("=" * 50)

print(f"Total Test Cases : {total}")
print(f"Valid JSON       : {json_pass}/{total}")
print(f"Valid Structure  : {structure_pass}/{total}")

print(f"\nResults saved to: {OUTPUT_PATH}")

client.close()