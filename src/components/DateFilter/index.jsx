import React from 'react';
import PropTypes from 'prop-types';
import format from 'date-fns/format';
import add from 'date-fns/add';
import Filter from 'components/Filter';
import './index.css';

const now = new Date();
const isoDateFormat = 'yyyy-MM-dd';

const today = format(now, isoDateFormat);
const oneWeekAgo = format(add(now, { weeks: -1 }), isoDateFormat);

function getDateFilter(label, placeholder, name, value = '', onChange) {
  return (
    <li className="date-filter__values">
      <label htmlFor={name} className="full">
        {label}
      </label>

      <input
        id={name}
        type="date"
        placeholder={placeholder}
        className="no-margin"
        value={value}
        onChange={(e) => onChange(name, e.target.value)}
      />
    </li>
  );
}

/**
 * DateFilter Functional Component.
 */
export default function DateFilter({ onChange, dates }) {
  return (
    <Filter label="date" count={Object.entries(dates).filter(Boolean).length}>
      {getDateFilter('After', oneWeekAgo, 'after', dates.after, onChange)}
      {getDateFilter('Before', today, 'before', dates.before, onChange)}
    </Filter>
  );
}

DateFilter.displayName = 'DateFilter';

DateFilter.propTypes = {
  dates: PropTypes.shape({
    before: PropTypes.string,
    after: PropTypes.string,
  }),
  onChange: PropTypes.func.isRequired,
};
