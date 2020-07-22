import React from 'react';
import PropTypes from 'prop-types';
import './index.css';

/**
 * Separator Functional Component.
 */
export default function Separator({ text }) {
  return <div className="separator full">{text}</div>;
}

Separator.displayName = 'Separator';

Separator.propTypes = {
  text: PropTypes.string.isRequired,
};
