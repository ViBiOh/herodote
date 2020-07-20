import React from 'react';
import { render } from '@testing-library/react';
import Filters from './index';

function defaultProps() {
  return {
    name: 'test',
    onChange: () => null,
    selected: [],
    values: [],
  };
}

it('should always render as a div', () => {
  const props = defaultProps();
  const { container } = render(<Filters {...props} />);
  expect(container.firstChild.nodeName).toBe('DIV');
});

it('should contains title with given name', () => {
  const props = defaultProps();
  props.name = 'repository';
  const { getByText } = render(<Filters {...props} />);

  expect(getByText('repository')).toBeInTheDocument();
});

it('should contains an entry for each values', () => {
  const props = defaultProps();
  props.values = [{ value: 'one' }, { value: 'two' }, { value: 'three' }];
  const { getByText } = render(<Filters {...props} />);

  props.values.forEach(({ value }) => {
    const item = getByText(value);
    expect(item).toBeInTheDocument();
    expect(item.nodeName).toBe('LABEL');
  });
});
