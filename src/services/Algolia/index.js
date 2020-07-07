import algoliasearch from 'algoliasearch/lite';

let index;

/**
 * Init algolia index
 * @param  {Object} config Configuration of app
 */
export function init(config) {
  if (!config || !config.ALGOLIA_APP || !config.ALGOLIA_KEY) {
    global.console.error('[algolia] config not provided');
    return;
  }

  const client = algoliasearch(config.ALGOLIA_APP, config.ALGOLIA_KEY);
  index = client.initIndex(config.ALGOLIA_INDEX);
}

/**
 * Perform algolia search
 * @param  {string} query   Query searched
 * @param  {Object} options Search options
 * @return {Object}         Algolia reponse
 */
export async function search(query, options = {}) {
  if (!index) {
    global.console.error('[algolia] index not initialized');
    return;
  }

  return await index.search(query, options);
}

/**
 * Perform algolia facets search
 * @param  {String} name Name of facet
 * @return {Object}      Algolia reponse
 */
export async function facets(name) {
  if (!index) {
    global.console.error('[algolia] index not initialized');
    return;
  }

  return await index.searchForFacetValues(name, '', {
    maxFacetHits: 100,
  });
}
