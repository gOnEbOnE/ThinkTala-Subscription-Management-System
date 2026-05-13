import json
from pathlib import Path
from graphify.extract import collect_files, extract
from graphify.cache import check_semantic_cache

raw = Path('.graphify_detect.json').read_bytes()
for enc in ('utf-16','utf-8-sig','utf-8'):
    try:
        detect = json.loads(raw.decode(enc))
        break
    except Exception:
        detect = None
if detect is None:
    raise SystemExit('Cannot decode detect file')

files = detect.get('files', {})
code_files = []
for f in files.get('code', []):
    p = Path(f)
    code_files.extend(collect_files(p) if p.is_dir() else [p])
if code_files:
    ast = extract(code_files)
else:
    ast = {'nodes':[],'edges':[],'input_tokens':0,'output_tokens':0}
Path('.graphify_ast.json').write_text(json.dumps(ast, indent=2), encoding='utf-8')

non_code = files.get('document', []) + files.get('paper', []) + files.get('image', []) + files.get('video', [])
cached_nodes, cached_edges, cached_hyperedges, uncached = check_semantic_cache(non_code)
if cached_nodes or cached_edges or cached_hyperedges:
    Path('.graphify_cached.json').write_text(json.dumps({'nodes':cached_nodes,'edges':cached_edges,'hyperedges':cached_hyperedges}, indent=2), encoding='utf-8')
Path('.graphify_uncached.txt').write_text('\n'.join(uncached), encoding='utf-8')

print('AST nodes={} edges={}'.format(len(ast.get('nodes',[])), len(ast.get('edges',[]))))
print('cache_hits={} uncached_non_code={}'.format(len(non_code)-len(uncached), len(uncached)))
