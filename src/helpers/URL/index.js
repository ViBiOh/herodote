/**
 * Retrieve search params in URL into object.
 * @return {Object} Search params
 */
export function search() {
  const params = {};

  window.location.search.replace(
    /([^?&=]+)(?:=([^?&=]*))?/g,
    (match, key, value) => {
      params[key] =
        typeof value === 'undefined' ? true : decodeURIComponent(value);
    },
  );

  return params;
}

/**
 * Push givens params into URL
 * @param  {Object} params Search params
 */
export function push(params) {
  window.history.pushState(
    null,
    null,
    `${window.location.pathname}${encode(params)}${window.location.hash}`,
  );
}

/**
 * Encode given params into the query string
 * @param  {Object} params Wanted params as object
 * @return {String}        Query string to append
 */
export function encode(params) {
  const encoded = Object.entries(params)
    .filter(([, value]) => Boolean(value))
    .map(([key, value]) => `${key}=${encodeURIComponent(value)}`);

  let searchStr = '';
  if (encoded.length) {
    searchStr = `?${encoded.join('&')}`;
  }

  return searchStr;
}
