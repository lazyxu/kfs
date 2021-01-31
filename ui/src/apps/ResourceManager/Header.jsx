import React from 'react';
import styled from 'styled-components';

import Path from 'apps/ResourceManager/Path';
import Branch from 'apps/ResourceManager/Branch';

const Header = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
  background-color: var(--header-color);
  display: flex;
`;

export default React.memo(({
  name, ...props
}) => {
  return (
    <Header>
      <Branch />
      <Path />
    </Header>
  );
});
