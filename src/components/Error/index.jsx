import React from 'react';
import PropTypes from 'prop-types';

/**
 * Error Functional Component.
 */
export default function Error({ error }) {
  return (
    <div className="danger">
      <h2>Error</h2>
      <pre>{`${error.message}`}</pre>
    </div>
  );
}

Error.displayName = 'Error';

Error.propTypes = {
  error: PropTypes.shape({
    message: PropTypes.string.isRequired,
  }).isRequired,
};
