import algoliasearch from 'algoliasearch/lite';

let index;

/**
 * Init algolia index
 * @param  {Object} config Configuration of app
 */
export function init(config) {
  if (!config || !config.ALGOLIA_APP || !config.ALGOLIA_KEY) {
    global.console.warn('[algolia] config not provided');
    return;
  }

  const client = algoliasearch(config.ALGOLIA_APP, config.ALGOLIA_KEY);
  index = client.initIndex(config.ALGOLIA_INDEX);
}

/**
 * Enabled tell if Algolia is enabled
 * @return {Boolean} True is enabled, false otherwise
 */
export function enabled() {
  return Boolean(index);
}

/**
 * Perform algolia search
 * @param  {string} query   Query searched
 * @param  {Object} options Search options
 * @return {Object}         Algolia reponse
 */
export async function search(query, filters = [], page = 0) {
  if (!index) {
    throw new Error('[algolia] index not initialized');
  }

  let filtersValue = '';
  if (filters.length) {
    filtersValue = `(${Object.values(
      filters.reduce((previous, current) => {
        const parts = current.split(':');
        if (parts.length < 1) {
          return previous;
        }

        const previousValue = previous[parts[0]];

        if (previousValue) {
          previous[parts[0]] = `${previousValue} OR ${current}`;
        } else {
          previous[parts[0]] = current;
        }

        return previous;
      }, {}),
    ).join(') AND (')})`;
  }

  return await index.search(query, { filters: filtersValue, page });
}

/**
 * Perform algolia facets search
 * @param  {String} name Name of facet
 * @return {Object}      Algolia reponse
 */
export async function facets(name, query) {
  if (!index) {
    throw new Error('[algolia] index not initialized');
  }

  return await index.searchForFacetValues(name, query, {
    maxFacetHits: 100,
    sortFacetValuesBy: 'alpha',
  });
}
