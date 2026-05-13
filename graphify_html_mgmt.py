import json
from pathlib import Path
from graphify.build import build_from_json
from graphify.export import to_html

# Load analysis data
extraction = json.loads(Path('management/.graphify_extract.json').read_text())
analysis = json.loads(Path('management/.graphify_analysis.json').read_text())
labels_raw = json.loads(Path('management/.graphify_labels.json').read_text()) if Path('management/.graphify_labels.json').exists() else {}

G = build_from_json(extraction)
communities = {int(k): v for k, v in analysis['communities'].items()}
labels = {int(k): v for k, v in labels_raw.items()} if labels_raw else {int(k): f'Community {k}' for k in communities.keys()}

if G.number_of_nodes() > 5000:
    print(f'Graph has {G.number_of_nodes()} nodes - too large for HTML viz.')
else:
    to_html(G, communities, 'management/graphify-out/graph.html', community_labels=labels or None)
    print('graph.html written')
