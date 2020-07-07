import React from 'react';
import PropTypes from 'prop-types';
import './index.css';

/**
 * Filters Functional Component.
 */
export default function Filters({ name, values, onChange, selected }) {
  return (
    <div className="filter">
      <h3 className="no-padding no-margin">{name}</h3>
      <ol className="filter__values no-padding no-margin">
        {values.map(({ value }) => {
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
              <label htmlFor={id} className="filter__label">
                {value}
              </label>
            </li>
          );
        })}
      </ol>
    </div>
  );
}

Filters.displayName = 'Filters';

Filters.propTypes = {
  name: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  selected: PropTypes.arrayOf(PropTypes.string).isRequired,
  values: PropTypes.arrayOf(
    PropTypes.shape({
      value: PropTypes.string.isRequired,
    }),
  ).isRequired,
};
