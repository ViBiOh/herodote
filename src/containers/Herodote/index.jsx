import React, { useState, useEffect } from 'react';
import {
  search as algoliaSearch,
  facets as algoliaFacets,
} from 'services/Algolia';
import { search as urlSearch, push as urlPush } from 'helpers/URL';
import debounced from 'helpers/Debounce';
import Filters from 'components/Filters';
import Commits from 'components/Commits';
import './index.css';

const filterRegex = /([^\s]+):([^\s]+)/gim;
const doubleSpaces = /\s{2,}/gim;
const debouncedCall = debounced();

function filterBy(query = '') {
  const filtersGroup = {};
  const filters = [];

  query = query.trim();
  query.replace(filterRegex, (all, key, value) => {
    filtersGroup[key] = (filtersGroup[key] || []).concat(value);
    filters.push(`${key}:${value}`);
  });

  const filtersValue = Object.entries(filtersGroup).map(([key, values]) =>
    values.map((value) => `${key}:${value}`).join(' OR '),
  );

  return [
    query.replace(filterRegex, ''),
    {
      filters: filtersValue.length ? `(${filtersValue.join(') AND (')})` : '',
    },
    filters,
  ];
}

async function fetchCommits(query, options, page) {
  const output = await algoliaSearch(query, { ...options, page });
  if (output) {
    return [output.hits, { next: output.page + 1, count: output.nbPages }];
  }

  return;
}

async function fetchFacets(name) {
  const output = await algoliaFacets(name, '');
  if (output) {
    return output.facetHits;
  }

  return [];
}

/**
 * Herodote Functional Component.
 */
export default function Herodote() {
  const [query, setQuery] = useState('');
  const [algoliaParams, setAlgoliaParams] = useState({});
  const [page, setPage] = useState(0);
  const [filters, setFilters] = useState([]);
  const [results, setResults] = useState([]);
  const [pagination, setPagination] = useState({ next: 0, count: 0 });
  const [facets, setFacets] = useState({});

  useEffect(() => {
    (async () => {
      const repository = await fetchFacets('repository');
      const type = await fetchFacets('type');
      const component = await fetchFacets('component');
      setFacets({ repository, type, component });
    })();

    const { query: searchQuery } = urlSearch();
    if (searchQuery) {
      setQuery(searchQuery);
    }
  }, []);

  useEffect(() => {
    const [algoliaQuery, algoliaOptions, selectedFilters] = filterBy(query);
    setAlgoliaParams({ query: algoliaQuery, options: algoliaOptions });
    setFilters(selectedFilters);
    setPage(0);

    urlPush({ query });
  }, [query]);

  useEffect(() => {
    debouncedCall(async () => {
      const [hits, newPagination] = await fetchCommits(
        algoliaParams.query,
        algoliaParams.options,
        page,
      );

      if (page) {
        setResults(results.concat(hits));
      } else {
        setResults(hits);
      }
      setPagination(newPagination);
    });
  }, [algoliaParams, page]); // eslint-disable-line react-hooks/exhaustive-deps

  /**
   * Handle filter change click event
   * @param  {Object} e     Click event on a checkbox
   * @param  {String} name  Facet's name
   * @param  {String} value Facet's value
   */
  const onFilterChange = (e, name, value) => {
    const filterValue = `${name}:${value}`;
    if (e.target.checked) {
      setQuery(`${query} ${filterValue}`.trim());
    } else {
      setQuery(
        query.replace(filterValue, '').replace(doubleSpaces, ' ').trim(),
      );
    }
  };

  const facetsFilters = Object.entries(facets)
    .filter(([_, values]) => values && values.length)
    .map(([key, values]) => (
      <Filters
        key={key}
        name={key}
        values={values}
        onChange={onFilterChange}
        selected={filters}
      />
    ));

  return (
    <div className="flex full">
      <aside className="padding">{facetsFilters}</aside>

      <article className="padding full">
        <input
          type="text"
          placeholder="Filter commit..."
          className="search padding full"
          onChange={(e) => setQuery(e.target.value)}
          value={query}
        />

        <hr />

        <Commits results={results} />

        {pagination.next < pagination.count && (
          <button
            type="button"
            className="button padding margin"
            onClick={() => setPage(pagination.next)}
          >
            Load more
          </button>
        )}
      </article>
    </div>
  );
}

Herodote.displayName = 'Herodote';
