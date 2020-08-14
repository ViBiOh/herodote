import React, { useState, useRef } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { useOnClickOutside } from 'helpers/Hooks';
import './index.css';

/**
 * Filter Functional Component.
 */
export default function Filter({ label, count, children }) {
  const ref = useRef();
  const [opened, toggle] = useState(false);
  useOnClickOutside(ref, () => toggle(false));

  let buttonLabel = label;
  if (count) {
    buttonLabel += ` (${count})`;
  }

  return (
    <span className="filter" ref={ref}>
      <button className="button bg-grey" onClick={() => toggle(!opened)}>
        {buttonLabel}
      </button>
      <ol
        className={classnames(
          'filter__values',
          'padding',
          'no-margin',
          'bg-grey',
          {
            'filer__values-active': opened,
          },
        )}
      >
        {children}
      </ol>
    </span>
  );
}

Filter.displayName = 'Filter';

Filter.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.arrayOf(PropTypes.node),
    PropTypes.node,
  ]),
  count: PropTypes.number.isRequired,
  label: PropTypes.string.isRequired,
};
