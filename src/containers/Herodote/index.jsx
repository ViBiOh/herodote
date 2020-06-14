import React, { Component } from 'react';
import { search } from 'services/Algolia';
import CommitsList from 'components/CommitsList';
import './index.css';

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

    this.state = {
      results: [],
    };
  }

  /**
   * React lifecycle.
   */
  async componentDidMount() {
    const results = await search('');

    this.setState({ results });
  }

  /**
   * React lifecycle.
   */
  render() {
    const { results } = this.state;
    if (!results.length) {
      return <p>No entry found</p>;
    }

    return <CommitsList results={results} />;
  }
}
