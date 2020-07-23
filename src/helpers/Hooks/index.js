import { useEffect, useState } from 'react';

/**
 * Trigger on click outside of ref
 * @param  {Object}   ref     Ref of the object, obtained with `useRef`
 * @param  {Function} handler Function to call when click is outide
 */
export function useOnClickOutside(ref, handler) {
  useEffect(() => {
    function listener(event) {
      if (!ref.current || ref.current.contains(event.target)) {
        return;
      }

      handler(event);
    }

    document.addEventListener('mousedown', listener);
    document.addEventListener('touchstart', listener);

    return () => {
      document.removeEventListener('mousedown', listener);
      document.removeEventListener('touchstart', listener);
    };
  }, [ref, handler]);
}

/**
 * Debounce call to an effect
 * @param {Function} fn           Debounced function
 * @param {Array}    dependencies List of effect's dependencies
 * @param {Number}   duration     Debounce duration
 */
export function useDebounce(fn, dependencies, duration = 300) {
  const [timeout, saveTimeout] = useState();

  useEffect(() => {
    clearTimeout(timeout);
    saveTimeout(setTimeout(fn, duration));
  }, [duration, ...dependencies]); // eslint-disable-line react-hooks/exhaustive-deps
}
