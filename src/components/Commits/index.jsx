import React from 'react';
import PropTypes from 'prop-types';
import endOfToday from 'date-fns/endOfToday';
import endOfDay from 'date-fns/endOfDay';
import isToday from 'date-fns/isToday';
import formatDistance from 'date-fns/formatDistance';
import { clear as clearColor, get as getColor } from 'services/Color';
import Separator from 'components/Separator';
import './index.css';

function getDateDistance(date) {
  if (isToday(date)) {
    return 'Today';
  }

  return `${formatDistance(endOfDay(date), endOfToday())} ago`;
}

/**
 * Commits Functional Component.
 */
export default function Commits({ commits, onClick }) {
  if (!commits.length) {
    return <p>No entry found</p>;
  }

  clearColor();
  let currentDistance;

  const commitsItems = [];
  for (let i = 0, size = commits.length; i < size; ++i) {
    const commit = commits[i];

    const dateDistance = getDateDistance(new Date(commit.date * 1000));
    if (currentDistance !== dateDistance) {
      currentDistance = dateDistance;
      commitsItems.push(
        <li key={currentDistance}>
          <Separator text={`${dateDistance}`} />
        </li>,
      );
    }

    commitsItems.push(
      <li key={commit.hash}>
        <button
          className="button button-rounded padding-half"
          style={{ backgroundColor: getColor(commit.repository) }}
          onClick={() => onClick(commit.repository)}
        >
          {commit.repository}
        </button>

        <pre className="label padding-half no-margin success">
          {commit.type}
          {commit.component && (
            <strong className="primary">({commit.component})</strong>
          )}
        </pre>

        <a
          className="commit-link ellipsis"
          href={`https://${commit.remote}/${commit.repository}/commit/${commit.hash}`}
          target="_blank"
          rel="noopener noreferrer"
        >
          {commit.content}
        </a>
      </li>,
    );
  }

  return (
    <ol id="commits" className="no-padding no-margin">
      {commitsItems}
    </ol>
  );
}

Commits.displayName = 'Commits';

Commits.propTypes = {
  commits: PropTypes.arrayOf(
    PropTypes.shape({
      content: PropTypes.string.isRequired,
      hash: PropTypes.string.isRequired,
      repository: PropTypes.string.isRequired,
    }),
  ).isRequired,
  onClick: PropTypes.func.isRequired,
};
