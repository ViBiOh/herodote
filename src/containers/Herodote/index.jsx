import React, { Component } from "react";
import { search } from "services/Algolia";
import CommitsList from "components/CommitsList";
import "./index.css";

const filterRegex = /([^\s]+:[^\s]+)/gim;

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
      q: "",
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
    this.filterBy("");
  }

  /**
   * React lifecycle.
   */
  componentWillUnmount() {
    this._isMounted = false;
  }

  /**
   * Filter content with given query
   * @param  {String} query Query string
   */
  filterBy = async (q, page = 0) => {
    const filters = [];
    q.replace(filterRegex, (all, filter) => {
      filters.push(filter);
    });

    const query = q.replace(filterRegex, "");

    try {
      const output = await search(query, page, {
        filters: filters.join(" AND "),
      });

      if (this._isMounted && output) {
        let list = output.hits;
        if (page > 0) {
          list = [...this.state.results, ...list];
        }

        this.setState({
          q,
          results: list,
          nextPage: page + 1,
          pagesCount: output.nbPages,
        });
      }
    } catch (e) {
      global.console.error(`unable to filter by '${q}': ${e}`);
    }
  };

  /**
   * Debounced tigger when user change input
   * @param  {Object} e Input
   */
  onSearchChange = (e) => {
    clearTimeout(this.timeout);

    ((text) => {
      this.timeout = setTimeout(() => this.filterBy(text), 300);
    })(e.target.value);
  };

  /**
   * Fetch next page
   */
  onMoreClick = () => {
    this.filterBy(this.state.q, this.state.nextPage);
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
