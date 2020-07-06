import { createBrowserHistory } from "history";

/**
 * BrowserHistory.
 */
export const history = createBrowserHistory();

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
        typeof value === "undefined" ? true : decodeURIComponent(value);
    }
  );

  return params;
}

/**
 * Push givens params into URL
 * @param  {Object} params Search params
 */
export function push(params) {
  const encoded = Object.entries(params)
    .filter(([, value]) => Boolean(value))
    .map(([key, value]) => `${key}=${encodeURIComponent(value)}`);

  let query = "";
  if (encoded.length) {
    query = `?${encoded.join("&")}`;
  }

  history.push(`/${query}`);
}
