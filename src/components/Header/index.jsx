import React from 'react';

/**
 * Header Functional Component.
 */
export default function Header() {
  return (
    <header className="padding">
      <h1 className="no-margin no-padding">
        <a href="/" className="no-style clear">
          Herodote
        </a>
      </h1>
    </header>
  );
}

Header.displayName = 'Header';
