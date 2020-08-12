import React, { useEffect, useState } from 'react';
import {
  search as algoliaSearch,
  enabled as algoliaEnabled,
} from 'services/Algolia';
import { search as apiSearch, enabled as apiEnabled } from 'services/Backend';
import { push as urlPush } from 'helpers/URL';
import { useDebounce } from 'helpers/Hooks';
import Commits from 'components/Commits';
import Error from 'components/Error';
import Throbber from 'components/Throbber';
import ThrobberButton from 'components/ThrobberButton';
import './index.css';

const filterRegex = /([^\s]+):([^\s]+)/gim;

function parseQuery(query = '') {
  const filters = [];
  query = query.trim();

  query.replace(filterRegex, (all, key, value) => {
    filters.push(`${key}:${value}`);
  });

  return [query.replace(filterRegex, ''), filters];
}

async function fetchCommits(query, filters, page) {
  if (apiEnabled()) {
    const output = await apiSearch(query, filters, page);
    if (output) {
      return [
        output.results || [],
        { next: output.page, count: output.pageCount },
      ];
    }
  } else if (algoliaEnabled()) {
    const output = await algoliaSearch(query, filters, page);
    if (output) {
      return [output.hits, { next: output.page + 1, count: output.nbPages }];
    }
  }

  return [];
}

/**
 * Herodote Functional Component.
 */
export default function Herodote({ query, filters, setFilters }) {
  const [pending, setPending] = useState(true);
  const [morePending, setMorePending] = useState(false);
  const [error, setError] = useState('');

  const [page, setPage] = useState(0);
  const [results, setResults] = useState([]);
  const [pagination, setPagination] = useState({ next: 0, count: 0 });

  useEffect(() => {
    setPending(true);
    setPage(0);

    urlPush({ query });
  }, [query]);

  useDebounce(async () => {
    try {
      const [newQuery, newFilters] = parseQuery(query);
      setFilters(newFilters);

      const [items, newPagination] = await fetchCommits(
        newQuery,
        newFilters,
        page,
      );

      if (page) {
        setResults(results.concat(items));
      } else {
        setResults(items);
      }

      setPagination(newPagination);
      setMorePending(false);
      setPending(false);
    } catch (e) {
      setError(e);
    }
  }, [query, filters, page]);

  if (error) {
    return <Error error={error} />;
  }

  if (pending) {
    return <Throbber label="Loading commits..." />;
  }

  return (
    <article>
      <Commits commits={results} />

      {pagination.next < pagination.count && (
        <ThrobberButton
          pending={morePending}
          onClick={() => setMorePending(true) || setPage(pagination.next)}
          type="button"
          className="button button-rounded padding margin margin-auto bg-primary"
        >
          Load more
        </ThrobberButton>
      )}
    </article>
  );
}

Herodote.displayName = 'Herodote';
