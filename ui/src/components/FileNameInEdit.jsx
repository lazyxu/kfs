import React from 'react';
import styled from 'styled-components';

const FileName = styled.p`
  font-size: 1em;
  padding: 0;
  overflow : hidden;
  text-overflow: ellipsis;
  border-radius: 0.3em;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  user-select: none;
  overflow-wrap: break-word;
  margin: 0;
`;

export default React.memo(({
  name, style,
}) => {
  return (
    <FileName style={style}>
      {name}
    </FileName>
  );
});
