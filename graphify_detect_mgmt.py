import sys, json
from graphify.detect import detect
from pathlib import Path

result = detect(Path('management'))
with open('management/.graphify_detect.json', 'w') as f:
    json.dump(result, f, indent=2)

print(f'Corpus: {result["total_files"]} files · ~{result["total_words"]:,} words')
for k, v in result['files'].items():
    if len(v) > 0:
        print(f'  {k}: {len(v)} files')
