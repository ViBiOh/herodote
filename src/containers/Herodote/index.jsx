import React, { useEffect, useState } from 'react';
import { search as algoliaSearch } from 'services/Algolia';
import { push as urlPush } from 'helpers/URL';
import { useDebounce } from 'helpers/Hooks';
import Commits from 'components/Commits';
import Error from 'components/Error';
import Throbber from 'components/Throbber';
import './index.css';

const filterRegex = /([^\s]+):([^\s]+)/gim;

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

  return [];
}

/**
 * Herodote Functional Component.
 */
export default function Herodote({ query, setFilters }) {
  const [pending, setPending] = useState(true);
  const [error, setError] = useState('');

  const [algoliaParams, setAlgoliaParams] = useState({});
  const [page, setPage] = useState(0);
  const [results, setResults] = useState([]);
  const [pagination, setPagination] = useState({ next: 0, count: 0 });

  useEffect(() => {
    setPending(true);
    const [algoliaQuery, algoliaOptions, selectedFilters] = filterBy(query);
    setAlgoliaParams({ query: algoliaQuery, options: algoliaOptions });
    setPage(0);
    urlPush({ query });

    setFilters(selectedFilters);
  }, [setFilters, query]);

  useDebounce(
    300,
    async () => {
      try {
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
        setPending(false);
      } catch (e) {
        setError(e);
      }
    },
    [algoliaParams, page],
  );

  if (error) {
    return <Error error={error} />;
  }

  if (pending) {
    return <Throbber label="Loading commits..." />;
  }

  return (
    <article>
      <Commits results={results} />

      {pagination.next < pagination.count && (
        <button
          type="button"
          className="button button-rounded padding margin margin-auto bg-primary"
          onClick={() => setPage(pagination.next)}
        >
          Load more
        </button>
      )}
    </article>
  );
}

Herodote.displayName = 'Herodote';
