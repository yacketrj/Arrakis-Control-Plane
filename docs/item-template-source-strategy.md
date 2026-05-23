# Item Template Source Strategy

## Decision

Use a hybrid item source:

1. Pull live item template identifiers from the database on connect and on an operator-controlled refresh cadence.
2. Keep `item-data.json` as a curated metadata fallback for display names, stack defaults, volume defaults, and templates that have not yet appeared in the database.
3. Serve the frontend from an in-memory merged cache instead of querying the database on every search keystroke.

This is more reliable than JSON alone and more economical than making every UI search query the database.

## Why not JSON only

A JSON-only source is cheap at runtime, but it can drift from the live server. It misses templates that appear after game updates, mod changes, server migrations, or player/import activity unless the JSON file is manually regenerated.

JSON is still useful for:

- Friendly display names.
- Stack maximum hints.
- Volume hints.
- Curated aliases.
- New or uncommon templates that have not appeared in the database yet.

## Why not query per UI search

Querying the database on every typeahead change is not recommended. The item table can grow large, and repeated `DISTINCT`, `ILIKE`, or trigram searches can add avoidable load while an operator is typing.

The more efficient model is:

- Refresh once at backend start/connect.
- Refresh on explicit operator request.
- Optionally refresh on a low-frequency timer.
- Search/filter the merged list in memory in the frontend or backend.

## Recommended database query

A common table expression is useful for readability and for combining live database templates with curated JSON/preset sources, but the CTE itself is not the performance optimization. Performance comes from using the query sparingly and caching the result.

Example live template query:

```sql
WITH live_templates AS (
  SELECT DISTINCT i.template_id
  FROM dune.items i
  WHERE i.template_id IS NOT NULL
    AND i.template_id <> ''
), observed_stats AS (
  SELECT
    i.template_id,
    MAX(i.stack_size) AS observed_max_stack,
    MAX(i.volume_override) FILTER (WHERE i.volume_override IS NOT NULL) AS observed_volume,
    COUNT(*) AS observed_count
  FROM dune.items i
  WHERE i.template_id IS NOT NULL
    AND i.template_id <> ''
  GROUP BY i.template_id
)
SELECT
  lt.template_id,
  COALESCE(os.observed_max_stack, 1) AS observed_max_stack,
  os.observed_volume,
  COALESCE(os.observed_count, 0) AS observed_count
FROM live_templates lt
LEFT JOIN observed_stats os ON os.template_id = lt.template_id
ORDER BY lt.template_id;
```

For a lightweight identifier-only refresh:

```sql
SELECT DISTINCT template_id
FROM dune.items
WHERE template_id IS NOT NULL
  AND template_id <> ''
ORDER BY template_id;
```

## Reliability and performance assessment

### Reliability

Database-backed discovery is reliable for templates that have already existed in the live database. It is not a complete catalog of every possible game template, because templates that have never been created on the server will not be represented in `dune.items`.

The hybrid approach covers both cases:

- Database provides current observed templates.
- JSON provides curated known templates and metadata.

### Performance

A scheduled or manual `SELECT DISTINCT template_id FROM dune.items` is economical when run infrequently. It is less economical if run repeatedly for every UI filter operation.

Recommended refresh behavior:

- At backend startup after DB connection.
- After `/api/v1/reconnect`.
- Optional manual refresh endpoint.
- Optional timer no more frequent than every 15 to 60 minutes for long-running admin sessions.

### Indexing

Do not create indexes automatically from the admin tool. If the table grows very large and refresh becomes slow, operators can evaluate a database-side index such as:

```sql
CREATE INDEX CONCURRENTLY IF NOT EXISTS items_template_id_idx
ON dune.items (template_id);
```

This should be an operator-controlled migration, not an automatic app startup action.

## Recommendation for Dune Admin

Keep the current merged-template model and evolve it into a cached hybrid provider:

- `dbItemTemplates` should be refreshed from the database on connect/reconnect.
- `itemData.Names` and `itemData.Items` should continue to merge into that list.
- The frontend should consume `/api/v1/players/templates` from the backend cache.
- Avoid direct browser-to-database template search.

This gives the best balance of correctness, cost, speed, and operational safety.
