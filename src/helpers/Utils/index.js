const doubleSpaces = /\s{2,}/gim;

/**
 * Toggle filter value in query
 * @param  {String} query   Current query
 * @param  {String} name    Filter's name
 * @param  {String} value   Filter's value
 * @param  {Boolean} unique Indicate if key is unique
 * @return {String}       New query
 */
export default function toggleFilter(query, name, value, unique) {
  const filterValue = `${name}:${value}`;
  if (
    query.indexOf(`${name}:`) === -1 ||
    (!unique && query.indexOf(filterValue) === -1)
  ) {
    return `${query} ${filterValue}`.trim();
  }

  if (unique) {
    if (value) {
      return query.replace(new RegExp(`${name}:[^\\s]+`), filterValue);
    }
    return query.replace(new RegExp(`${name}:[^\\s]+`), '').trim();
  }

  return query.replace(filterValue, '').replace(doubleSpaces, ' ').trim();
}
