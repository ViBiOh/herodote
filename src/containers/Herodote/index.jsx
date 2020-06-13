import React, { Component } from "react";
import classNames from "classnames";
import { search } from "services/Algolia";
import "./index.css";

const COLOR_MAPPING = ["blue", "green", "red"];

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
    const results = await search("");

    const colorMapping = {};
    let index = 0;

    results.forEach((result) => {
      result.label = colorMapping[result.repository];
      if (!result.label) {
        colorMapping[result.repository] = result.label = COLOR_MAPPING[index++];
      }
    });

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

    return (
      <ol id="commits" className="no-padding">
        {results.map((result) => (
          <li key={result.hash}>
            <span className={classNames("label", `label-${result.label}`)}>
              {result.repository}
            </span>
            <span>{result.content}</span>
          </li>
        ))}
      </ol>
    );
  }
}
