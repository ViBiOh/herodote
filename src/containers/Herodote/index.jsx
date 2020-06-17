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
      results: [],
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
  filterBy = async (q) => {
    const filters = [];
    q.replace(filterRegex, (all, filter) => {
      filters.push(filter);
    });

    const query = q.replace(filterRegex, "");

    const results = await search(query, { filters: filters.join(" AND ") });
    if (this._isMounted) {
      this.setState({ results });
    }
  };

  /**
   * Debounced tigger when user change input
   * @param  {String} e Input
   */
  onSearchChange = (e) => {
    clearTimeout(this.timeout);

    ((text) => {
      this.timeout = setTimeout(() => this.filterBy(text), 300);
    })(e.target.value);
  };

  /**
   * React lifecycle.
   */
  render() {
    const { results } = this.state;

    return (
      <article>
        <input
          type="text"
          placeholder="Filter commit..."
          onChange={this.onSearchChange}
        />
        <CommitsList results={results} />
      </article>
    );
  }
}
