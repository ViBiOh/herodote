import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';
import './index.css';

/**
 * Throbber for displaying background task.
 * @param {Object} props Props of the component.
 * @return {React.Component} Throbber with label and title if provided
 */
export default function Throbber({ label, title, className }) {
  const classes = classnames('throbber', className);

  return (
    <div className="container" title={title}>
      {label && <span>{label}</span>}

      <div className={classes}>
        <div className="bounce1" />
        <div className="bounce2" />
        <div className="bounce3" />
      </div>
    </div>
  );
}

Throbber.displayname = 'Throbber';

Throbber.propTypes = {
  label: PropTypes.string,
  title: PropTypes.string,
  className: PropTypes.string,
};

Throbber.defaultProps = {
  label: '',
  title: '',
  className: '',
};
