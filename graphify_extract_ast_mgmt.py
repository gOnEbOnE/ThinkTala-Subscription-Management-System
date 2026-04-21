import sys, json
from graphify.extract import collect_files, extract
from pathlib import Path

# AST Extraction
code_files = []
detect = json.loads(Path('management/.graphify_detect.json').read_text())
for f in detect.get('files', {}).get('code', []):
    code_files.extend(collect_files(Path(f)) if Path(f).is_dir() else [Path(f)])

if code_files:
    print(f'Extracting AST from {len(code_files)} code file(s)...')
    result = extract(code_files)
    Path('management/.graphify_ast.json').write_text(json.dumps(result, indent=2))
    print(f'AST: {len(result["nodes"])} nodes, {len(result["edges"])} edges')
else:
    Path('management/.graphify_ast.json').write_text(json.dumps({'nodes':[],'edges':[],'input_tokens':0,'output_tokens':0}))
    print('No code files - skipping AST extraction')
