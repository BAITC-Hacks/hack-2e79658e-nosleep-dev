# Repository instructions for AI agents

## Brand

Read `docs/BRAND.md` before changing product naming, marketing copy, metadata,
titles, descriptions, or visual identity. The canonical product spelling is
**BesSheshim**.

## Knowledge graph

Treat `graphify-out/graph.json` as the repository knowledge graph and keep it in
sync with meaningful code or documentation changes.

Before answering repository-wide questions or making cross-cutting changes:

1. Read `graphify-out/GRAPH_REPORT.md`.
2. Query the existing graph with `graphify query "<question>"` when it can help
   locate concepts, dependencies, or relevant files.
3. Use source files as the final authority; graph relationships may include
   explicitly marked inferences.

After meaningful code or documentation changes, update the graph with the
graphify skill using `/graphify . --update`. If no graph exists yet, build it
with `/graphify .`. Commit the refreshed `graphify-out/graph.json`,
`graphify-out/GRAPH_REPORT.md`, and `graphify-out/graph.html` with the change.
