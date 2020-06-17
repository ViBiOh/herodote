import algoliasearch from "algoliasearch/lite";

let index;

export function init(config) {
  if (!config || !config.ALGOLIA_APP || !config.ALGOLIA_KEY) {
    global.console.error("[algolia] config not provided");
    return;
  }

  const client = algoliasearch(config.ALGOLIA_APP, config.ALGOLIA_KEY);
  index = client.initIndex(config.ALGOLIA_INDEX);
}

export async function search(query, options = {}) {
  if (!index) {
    global.console.error("[algolia] index not initialized");
    return [];
  }

  const output = await index.search(query, options);
  return output.hits;
}
