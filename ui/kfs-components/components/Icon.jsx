import React from 'react';
import styled from 'styled-components';

export default function ({
  hoverCursor, hoverColor, margin, icon, size, color, onClick, padding, style, ...props
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
    <Icon
      aria-hidden="true"
      onClick={(e) => { onClick && onClick(e); e.stopPropagation(); }}
      style={style}
      width="100%"
      height="100%"
      viewBox="0 0 200 200"
      preserveAspectRatio="xMinYMin meet"
      {...props}
    >
      <use data-tag="choose-able" xlinkHref={`#icon-${icon}`} />
    </Icon>
  );
}
