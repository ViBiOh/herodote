import React, { useState, useRef } from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import { useOnClickOutside } from 'helpers/Hooks';
import './index.css';

/**
 * Filter Functional Component.
 */
export default function Filter({ name, values, onChange, selected }) {
  const ref = useRef();
  const [opened, toggle] = useState(false);
  useOnClickOutside(ref, () => toggle(false));

  const count = selected.filter((e) => e.startsWith(name)).length;
  let buttonLabel = name;
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
        {values.map((value) => {
          const id = `${name}:${value}`;
          return (
            <li key={id}>
              <input
                id={id}
                type="checkbox"
                value={value}
                onChange={(e) => onChange(e, name, value)}
                checked={selected.includes(id)}
              />
              <label htmlFor={id} className="filter__label ellipsis">
                {value}
              </label>
            </li>
          );
        })}
      </ol>
    </span>
  );
}

Filter.displayName = 'Filter';

Filter.propTypes = {
  name: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  selected: PropTypes.arrayOf(PropTypes.string).isRequired,
  values: PropTypes.arrayOf(PropTypes.string).isRequired,
};
