import funtch from 'funtch';

let funtcher;

export function init(config) {
  if (!config || !config.HERODOTE_API) {
    global.console.warn('[api] config not provided');
    return;
  }

  funtcher = funtch.withDefault({
    baseURL: config.HERODOTE_API,
  });
}

/**
 * Enabled tell if backend with API is enabled
 * @return {Boolean} True is enabled, false otherwise
 */
export function enabled() {
  return Boolean(funtcher);
}

/**
 * Perform herodote backend search
 * @param  {string} query   Query searched
 * @param  {Object} options Search options
 * @return {Object}         API response
 */
export async function search(query, filters = [], dates = {}, page = 0) {
  if (!funtcher) {
    throw new Error('[api] config not initialized');
  }

  const filtersValue = filters
    .map((value) => {
      const parts = value.split(':', 2);
      if (parts.length !== 2) {
        return '';
      }

      return `${parts[0]}=${encodeURIComponent(parts[1])}`;
    })
    .filter(Boolean)
    .join('&');

  const datesValue = Object.entries(dates)
    .filter(([, value]) => Boolean(value))
    .map(([key, value]) => `${key}=${encodeURIComponent(value)}`)
    .join('&');

  let params = '';
  if (filtersValue) {
    params += `&${filtersValue}`;
  }
  if (datesValue) {
    params += `&${datesValue}`;
  }

  return await funtcher.get(
    `/commits?page=${encodeURIComponent(page + 1)}&q=${encodeURIComponent(
      query,
    )}${params}`,
  );
}

/**
 * Perform API filters search
 * @param  {String} name Name of facet
 * @return {Object}      API response
 */
export async function filters(name) {
  if (!funtcher) {
    throw new Error('[api] config not initialized');
  }

  return await funtcher.get(`/filters?name=${encodeURIComponent(name)}`);
}
