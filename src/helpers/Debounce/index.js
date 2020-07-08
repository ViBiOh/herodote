/**
 * Debounce call of given func
 * @param  {Number} duration Debounce duration
 * @return {Function} Debounced function call
 */
export default function (duration = 300) {
  let timeout;

  return function (fn, ...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => fn(...args), duration);
  };
}
