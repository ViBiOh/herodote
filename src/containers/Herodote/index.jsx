import React, { Component } from "react";
import { search as algoliaSearch } from "services/Algolia";
import { search as urlSearch, push as urlPush } from "helpers/URL";
import CommitsList from "components/CommitsList";
import "./index.css";

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
      query: "",
      results: [],
      nextPage: 0,
      pagesCount: 0,
    };
  }

  /**
   * React lifecycle.
   */
  componentDidMount() {
    this._isMounted = true;

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
   * Filter content with given query
   * @param  {String} query Query string
   */
  filterBy = async (query = "", page) => {
    const filtersValues = {};
    query.replace(filterRegex, (all, key, value) => {
      filtersValues[key] = (filtersValues[key] || []).concat(value);
    });

    const filters = Object.entries(filtersValues).map(([key, values]) =>
      values.map((value) => `${key}:${value}`).join(" OR ")
    );

    await this.fetchCommits(query.replace(filterRegex, ""), filters, page);

    this.setState({ query });
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
      filters: filters.length ? `(${filters.join(") AND (")})` : "",
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
    console.log(
      "this.state.nextPage",
      typeof this.state.nextPage,
      this.state.nextPage
    );
    this.filterBy(this.state.query, this.state.nextPage);
  };

  /**
   * React lifecycle.
   */
  render() {
    const { results, nextPage, pagesCount } = this.state;

    return (
      <article className="padding">
        <input
          type="text"
          placeholder="Filter commit..."
          className="search padding"
          onChange={this.onSearchChange}
          value={this.state.query}
        />

        <hr />

        <CommitsList results={results} />

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
    );
  }
}
