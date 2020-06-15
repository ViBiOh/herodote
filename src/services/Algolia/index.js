import algoliasearch from "algoliasearch/lite";

let index;

export function init(config) {
  if (!config || !config.ALGOLIA_APPLICATION_ID || !config.ALGOLIA_API_KEY) {
    global.console.error("[algolia] config not provided");
  }

  const client = algoliasearch(
    config.ALGOLIA_APPLICATION_ID,
    config.ALGOLIA_API_KEY
  );
  index = client.initIndex(config.ALGOLIA_INDEX);
}

export async function search(query, options = {}) {
  if (!index) {
    global.console.error("[algolia] index not initialized");
  }

  const output = await index.search(query, options);
  return output.hits;
}
