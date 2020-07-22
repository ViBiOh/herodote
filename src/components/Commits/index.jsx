import React from 'react';
import { clear as clearColor, get as getColor } from 'services/Color';
import PropTypes from 'prop-types';
import './index.css';

/**
 * Commits Functional Component.
 */
export default function Commits({ results }) {
  if (!results.length) {
    return <p>No entry found</p>;
  }

  clearColor();

  return (
    <ol id="commits" className="no-padding no-margin">
      {results.map((result) => (
        <li key={result.hash}>
          <span
            className="label"
            style={{ backgroundColor: getColor(result.repository) }}
          >
            {result.repository}
          </span>

          <pre className="label no-margin success">
            {result.type}
            {result.component && (
              <strong className="primary">({result.component})</strong>
            )}
          </pre>
          <a
            className="commit-link ellipsis"
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

Commits.displayName = 'Commits';

Commits.propTypes = {
  results: PropTypes.arrayOf(
    PropTypes.shape({
      hash: PropTypes.string.isRequired,
      content: PropTypes.string.isRequired,
      repository: PropTypes.string.isRequired,
    }),
  ).isRequired,
};
