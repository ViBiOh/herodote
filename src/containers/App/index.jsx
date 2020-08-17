import React, { useEffect, useState } from 'react';
import getConfig from 'services/Config';
import { init as initAlgolia } from 'services/Algolia';
import { init as initAPI } from 'services/Backend';
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
  const [dates, setDates] = useState({});

  useEffect(() => {
    (async () => {
      const rawConfig = await getConfig();
      initAPI(rawConfig);
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
          <Filters
            query={query}
            onChange={setQuery}
            filters={filters}
            dates={dates}
          />
        )}

        {config && (
          <Herodote
            query={query}
            onChange={setQuery}
            filters={filters}
            setFilters={setFilters}
            dates={dates}
            setDates={setDates}
          />
        )}
      </div>
    </div>
  );
}

App.displayName = 'App';
