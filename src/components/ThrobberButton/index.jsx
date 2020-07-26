import React from 'react';
import PropTypes from 'prop-types';
import Throbber from 'components/Throbber';

/**
 * ThrobberButton Functional Component.
 */
export default function ThrobberButton({
  pending,
  onClick,
  children,
  ...buttonProps
}) {
  return (
    <button {...buttonProps} onClick={(e) => (pending ? null : onClick(e))}>
      {pending ? <Throbber /> : children}
    </button>
  );
}

ThrobberButton.displayName = 'ThrobberButton';

ThrobberButton.propTypes = {
  children: PropTypes.oneOfType([
    PropTypes.node,
    PropTypes.arrayOf(PropTypes.node),
  ]).isRequired,
  onClick: PropTypes.func.isRequired,
  pending: PropTypes.bool.isRequired,
};
