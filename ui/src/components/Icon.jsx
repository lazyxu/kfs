import React from 'react';
import styled from 'styled-components';

export default function ({
  hoverCursor, hoverColor, margin, icon, size, color, onClick, padding,
}) {
  const Icon = styled.svg`
    :hover {
      cursor: ${hoverCursor};
      fill:   ${hoverColor};
    }
    margin: ${margin};
    height: ${size};
    width: ${size};
    vertical-align: -0.15em;
    fill: ${color};
    padding: ${padding};
  `;
  return (
    // eslint-disable-next-line no-unused-expressions
    <Icon aria-hidden="true" onClick={(e) => { onClick && onClick(e); e.stopPropagation(); }}>
      <use xlinkHref={`#icon-${icon}`} />
    </Icon>
  );
}
