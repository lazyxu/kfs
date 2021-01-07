import React from 'react';

import styled from 'styled-components';

import Icon from 'components/Icon';
import ConfigEditor from 'containers/ConfigEditor';
import Path from 'containers/Path';
import { newWindow } from 'components/Modal';

const Header = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
  background-color: var(--header-color);
  display: flex;
`;
const Empty = styled.div`
  flex: 1;
`;

export default React.memo(({
  name, ...props
}) => {
  return (
    <Header>
      <Path />
      <Empty />
      <Icon
        icon="setting"
        color="#cccccc"
        size="1.8em"
        hoverColor="white"
        hoverCursor="pointer"
        onClick={e => {
          newWindow(ConfigEditor);
        }}
      />
    </Header>
  );
});
