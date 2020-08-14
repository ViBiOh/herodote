import React, { useState, useEffect } from 'react';
import {
  facets as algoliaFacets,
  enabled as algoliaEnabled,
} from 'services/Algolia';
import { filters as apiFilters, enabled as apiEnabled } from 'services/Backend';
import PropTypes from 'prop-types';
import AlgoliaLogo from 'components/AlgoliaLogo';
import ListFilter from 'components/ListFilter';
import DateFilter from 'components/DateFilter';
import Error from 'components/Error';
import Throbber from 'components/Throbber';
import './index.css';

const doubleSpaces = /\s{2,}/gim;

async function fetchFacets(name) {
  if (apiEnabled()) {
    const output = await apiFilters(name);
    if (output) {
      return output.results || [];
    }
  } else if (algoliaEnabled()) {
    const output = await algoliaFacets(name, '');
    if (output) {
      return output.facetHits.map((o) => o.value);
    }
  }

  return [];
}

async function loadFacets() {
  const repository = await fetchFacets('repository');
  const type = await fetchFacets('type');
  const component = await fetchFacets('component');

  return { repository, type, component };
}

/**
 * Filters Functional Component.
 */
export default function Filters({ query, onChange, filters, dates }) {
  const [pending, setPending] = useState(true);
  const [error, setError] = useState('');

  const [facets, setFacets] = useState({});

  useEffect(() => {
    (async () => {
      try {
        setFacets(await loadFacets());
        setPending(false);
      } catch (e) {
        setError(e);
      }
    })();
  }, []);

  useEffect(() => {
    onChange(query);
  }, [onChange, query]);

  if (error) {
    return <Error error={error} />;
  }

  if (pending) {
    return <Throbber label="Loading filters..." />;
  }

  /**
   * Handle date change click event
   * @param  {Object} e     Click event on a checkbox
   * @param  {String} name  Facet's name
   * @param  {String} value Facet's value
   */
  const onDateChange = (name, value) => {
    if (value) {
      const filterValue = `${name}:${value}`;

      if (query.indexOf(`${name}:`) !== -1) {
        onChange(query.replace(new RegExp(`${name}:[^\\s]+`), filterValue));
      } else {
        onChange(`${query} ${filterValue}`.trim());
      }
    } else {
      onChange(query.replace(new RegExp(`${name}:[^\\s]+`), '').trim());
    }
  };

  /**
   * Handle filter change click event
   * @param  {Object} e     Click event on a checkbox
   * @param  {String} name  Facet's name
   * @param  {String} value Facet's value
   */
  const onFilterChange = (e, name, value) => {
    const filterValue = `${name}:${value}`;
    if (e.target.checked) {
      onChange(`${query} ${filterValue}`.trim());
    } else {
      onChange(
        query.replace(filterValue, '').replace(doubleSpaces, ' ').trim(),
      );
    }
  };

  return (
    <aside id="filters" className="flex full">
      {!apiEnabled() && algoliaEnabled() && (
        <AlgoliaLogo
          className="algolia-logo"
          height={38}
          title="Search by Algolia"
        />
      )}

      <input
        type="text"
        aria-label="Filter commits"
        placeholder="Filter commits..."
        className="no-border no-margin padding search"
        onChange={(e) => onChange(e.target.value)}
        value={query}
      />

      <DateFilter onChange={onDateChange} dates={dates} />

      {Object.entries(facets)
        .filter(([_, values]) => values && values.length)
        .map(([key, values]) => (
          <ListFilter
            key={key}
            name={key}
            values={values}
            onChange={onFilterChange}
            selected={filters}
          />
        ))}
    </aside>
  );
}

Filters.displayName = 'Filters';

Filters.propTypes = {
  dates: PropTypes.shape({}).isRequired,
  filters: PropTypes.arrayOf(PropTypes.string).isRequired,
  onChange: PropTypes.func.isRequired,
  query: PropTypes.string.isRequired,
};
