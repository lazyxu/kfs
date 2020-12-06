import React from 'react';

import styled from 'styled-components';

import Icon from 'components/Icon';
import Path from 'containers/Path';

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
class component extends React.Component {
  render() {
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
        />
      </Header>
    );
  }
}

export default component;
