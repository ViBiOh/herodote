import React, { Component } from 'react';
import {
  search as algoliaSearch,
  facets as algoliaFacets,
} from 'services/Algolia';
import { search as urlSearch, push as urlPush } from 'helpers/URL';
import Filters from 'components/Filters';
import Commits from 'components/Commits';
import './index.css';

const filterRegex = /([^\s]+):([^\s]+)/gim;

/**
 * Herodote Component.
 */
export default class Herodote extends Component {
  /**
   * Creates an instance of Herodote.
   * @param {Object} props Component props
   */
  constructor(props) {
    super(props);

    this._isMounted = false;

    this.state = {
      query: '',
      filters: [],
      results: [],
      nextPage: 0,
      pagesCount: 0,
      facets: {},
    };
  }

  /**
   * React lifecycle.
   */
  componentDidMount() {
    this._isMounted = true;

    this.loadFacets();

    const { query } = urlSearch();
    this.filterBy(query);
  }

  /**
   * React lifecycle.
   */
  componentWillUnmount() {
    this._isMounted = false;
  }

  /**
   * React lifecycle.
   */
  componentDidCatch(error, info) {
    global.console.error(error, info.componentStack);
  }

  /**
   * Load facets in appropriate order
   * @return {[type]} [description]
   */
  loadFacets = async () => {
    const repository = await this.fetchFacets('repository');
    const type = await this.fetchFacets('type');
    const component = await this.fetchFacets('component');

    if (this._isMounted) {
      this.setState({ facets: { repository, type, component } });
    }
  };

  /**
   * Fetch facets of given name
   * @param  {String} name Name of facets
   */
  fetchFacets = async (name) => {
    const output = await algoliaFacets(name);
    if (output) {
      return output.facetHits;
    }

    return [];
  };

  /**
   * Filter content with given query
   * @param  {String} query Query string
   */
  filterBy = async (query = '', page) => {
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

    await this.fetchCommits(query.replace(filterRegex, ''), filtersValue, page);
    this.setState({ query, filters });
    urlPush({ query });
  };

  /**
   * Fetch commits from algolia backend
   * @param  {String} query   Search query
   * @param  {Array}  filters List of filters
   * @param  {Number} page    Page, 0 based
   */
  fetchCommits = async (query, filters, page = 0) => {
    const output = await algoliaSearch(query, {
      filters: filters.length ? `(${filters.join(') AND (')})` : '',
      page,
    });

    if (this._isMounted && output) {
      let results = output.hits;
      if (page) {
        results = this.state.results.concat(results);
      }

      this.setState({
        results,
        nextPage: output.page + 1,
        pagesCount: output.nbPages,
      });
    }
  };

  /**
   * Debounced trigger when user change input
   * @param  {Object} e Input
   */
  onSearchChange = (e) => {
    clearTimeout(this.timeout);

    const query = e.target.value;
    this.setState({ query });

    ((text) => {
      this.timeout = setTimeout(() => this.filterBy(text), 300);
    })(query);
  };

  /**
   * Fetch next page
   */
  onMoreClick = () => {
    this.filterBy(this.state.query, this.state.nextPage);
  };

  /**
   * Update query on filter change
   * @param  {Object} e     Event from the checkbox
   * @param  {String} name  Name of the filter
   * @param  {String} value Value of the filter
   */
  onFilterChange = (e, name, value) => {
    const filterValue = `${name}:${value}`;
    let { query } = this.state;

    if (e.target.checked) {
      this.filterBy(`${query} ${filterValue}`);
    } else {
      this.filterBy(query.replace(filterValue, ''));
    }
  };

  /**
   * React lifecycle.
   */
  render() {
    const { results, filters, nextPage, pagesCount, facets } = this.state;

    const facetsFilters = Object.entries(facets)
      .filter(([_, values]) => values && values.length)
      .map(([key, values]) => (
        <Filters
          key={key}
          name={key}
          values={values}
          onChange={this.onFilterChange}
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
            onChange={this.onSearchChange}
            value={this.state.query}
          />

          <hr />

          <Commits results={results} />

          {nextPage < pagesCount && (
            <button
              type="button"
              className="button padding margin"
              onClick={this.onMoreClick}
            >
              Load more
            </button>
          )}
        </article>
      </div>
    );
  }
}
