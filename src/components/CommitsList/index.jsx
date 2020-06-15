import React from "react";
import { clear as clearColor, get as getColor } from "services/Color";
import PropTypes from "prop-types";

/**
 * CommitsList Functional Component.
 */
export default function CommitsList({ results }) {
  if (!results.length) {
    return <p>No entry found</p>;
  }

  clearColor();

  return (
    <ol id="commits" className="no-padding">
      {results.map((result) => (
        <li key={result.hash}>
          <span
            className="label"
            style={{ backgroundColor: getColor(result.repository) }}
          >
            {result.repository}
          </span>

          <pre className="label no-margin">
            {result.type}
            {result.component && `(${result.component})`}
          </pre>
          <a
            href={`https://${result.remote}/${result.repository}/commit/${result.hash}`}
            target="_blank"
            rel="noopener noreferrer"
          >
            {result.content}
          </a>
        </li>
      ))}
    </ol>
  );
}

CommitsList.displayName = "CommitsList";

CommitsList.propTypes = {
  results: PropTypes.arrayOf(
    PropTypes.shape({
      hash: PropTypes.string.isRequired,
      content: PropTypes.string.isRequired,
      repository: PropTypes.string.isRequired,
    })
  ).isRequired,
};
