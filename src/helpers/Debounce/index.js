import { useState, useEffect } from 'react';

/**
 * Debounce call to an effect
 * @param {Number}   duration     Debounce duration
 * @param {Function} fn           Debounced function
 * @param {Array}    dependencies List of effect's dependencies
 */
export function Debounce(duration = 300, fn, dependencies) {
  const [timeout, saveTimeout] = useState();

  useEffect(() => {
    clearTimeout(timeout);
    saveTimeout(setTimeout(fn, duration));
  }, [duration, ...dependencies]); // eslint-disable-line react-hooks/exhaustive-deps
}
