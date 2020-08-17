import React from 'react';
import PropTypes from 'prop-types';
import Filter from 'components/Filter';

/**
 * ListFilter Functional Component.
 */
export default function ListFilter({ name, values, onChange, selected }) {
  return (
    <Filter
      label={name}
      count={selected.filter((e) => e.startsWith(name)).length}
    >
      {values.map((value) => {
        const id = `${name}:${value}`;
        return (
          <li key={id}>
            <input
              id={id}
              type="checkbox"
              value={value}
              onChange={() => onChange(value)}
              checked={selected.includes(id)}
            />
            <label htmlFor={id} className="filter__label ellipsis">
              {value}
            </label>
          </li>
        );
      })}
    </Filter>
  );
}

ListFilter.displayName = 'ListFilter';

ListFilter.propTypes = {
  name: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  selected: PropTypes.arrayOf(PropTypes.string).isRequired,
  values: PropTypes.arrayOf(PropTypes.string).isRequired,
};
