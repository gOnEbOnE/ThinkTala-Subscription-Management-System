import sys, json
from graphify.build import build_from_json
from graphify.cluster import cluster, score_all
from graphify.analyze import god_nodes, surprising_connections, suggest_questions
from graphify.report import generate
from graphify.export import to_json
from pathlib import Path

# Since we have only code files (no semantic extraction needed), just use AST
ast = json.loads(Path('management/.graphify_ast.json').read_text())
extraction = ast  # Direct copy for code-only corpus

extraction['input_tokens'] = ast.get('input_tokens', 0)
extraction['output_tokens'] = ast.get('output_tokens', 0)

Path('management/.graphify_extract.json').write_text(json.dumps(extraction, indent=2))

detection = json.loads(Path('management/.graphify_detect.json').read_text())

# Build and analyze
G = build_from_json(extraction)
communities = cluster(G)
cohesion = score_all(G, communities)
tokens = {'input': extraction.get('input_tokens', 0), 'output': extraction.get('output_tokens', 0)}
gods = god_nodes(G)
surprises = surprising_connections(G, communities)
labels = {cid: f'Community {cid}' for cid in communities}
questions = suggest_questions(G, communities, labels)

report = generate(G, communities, cohesion, labels, gods, surprises, detection, tokens, 'management')
Path('management/graphify-out').mkdir(parents=True, exist_ok=True)
Path('management/graphify-out/GRAPH_REPORT.md').write_text(report, encoding='utf-8')
to_json(G, communities, 'management/graphify-out/graph.json')

analysis = {
    'communities': {str(k): v for k, v in communities.items()},
    'cohesion': {str(k): v for k, v in cohesion.items()},
    'gods': gods,
    'surprises': surprises,
    'questions': questions,
}
Path('management/.graphify_analysis.json').write_text(json.dumps(analysis, indent=2))

print(f'Graph: {G.number_of_nodes()} nodes, {G.number_of_edges()} edges, {len(communities)} communities')
print(f'God nodes: {len(gods)}')
print()
print('Communities:')
for cid in sorted(communities.keys()):
    nodes_in_community = len(communities[cid])
    print(f'  {cid}: {nodes_in_community} nodes')
