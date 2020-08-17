import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import {
  facets as algoliaFacets,
  enabled as algoliaEnabled,
} from 'services/Algolia';
import { filters as apiFilters, enabled as apiEnabled } from 'services/Backend';
import toggleFilter from 'helpers/Utils';
import AlgoliaLogo from 'components/AlgoliaLogo';
import ListFilter from 'components/ListFilter';
import DateFilter from 'components/DateFilter';
import Error from 'components/Error';
import Throbber from 'components/Throbber';
import './index.css';

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

      <DateFilter
        onChange={(key, value) =>
          onChange(toggleFilter(query, key, value, true))
        }
        dates={dates}
      />

      {Object.entries(facets)
        .filter(([_, values]) => values && values.length)
        .map(([key, values]) => (
          <ListFilter
            key={key}
            name={key}
            values={values}
            onChange={(value) => onChange(toggleFilter(query, key, value))}
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
