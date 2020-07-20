import React from 'react';
import PropTypes from 'prop-types';

/**
 * Error Functional Component.
 */
export default function Error({ error }) {
  return (
    <div className="danger">
      <h2>Error</h2>
      <pre>{`${error}`}</pre>
    </div>
  );
}

Error.displayName = 'Error';

Error.propTypes = {
  error: PropTypes.shape({}).isRequired,
};
