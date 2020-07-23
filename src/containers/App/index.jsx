import React, { useEffect, useState } from 'react';
import getConfig from 'services/Config';
import { init as initAlgolia } from 'services/Algolia';
import { search as urlSearch } from 'helpers/URL';
import Header from 'components/Header';
import Filters from 'containers/Filters';
import Herodote from 'containers/Herodote';
import Throbber from 'components/Throbber';

/**
 * App Functional Component.
 */
export default function App() {
  const [config, setConfig] = useState();
  const [query, setQuery] = useState('');
  const [filters, setFilters] = useState([]);

  useEffect(() => {
    (async () => {
      const rawConfig = await getConfig();
      initAlgolia(rawConfig);
      setConfig(rawConfig);
    })();
  }, []);

  useEffect(() => {
    const { query: searchQuery } = urlSearch();
    if (searchQuery) {
      setQuery(searchQuery);
    }
  }, []);

  return (
    <div className="content">
      <Header />

      <div className="padding full">
        {!config && <Throbber label="Loading configuration..." />}
        {config && (
          <Filters query={query} onChange={setQuery} filters={filters} />
        )}
        {config && <Herodote query={query} setFilters={setFilters} />}
      </div>
    </div>
  );
}

App.displayName = 'App';
