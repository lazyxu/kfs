import React from 'react';

import styled from 'styled-components';

import Icon from 'components/Icon';
import ConfigEditor from 'containers/ConfigEditor';
import Path from 'containers/Path';
import { busState, setState } from 'bus/bus';

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
  const [isOpen, setOpen] = React.useState(false);
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
        onClick={e => { setOpen(true); }}
      />
      <ConfigEditor
        isOpen={isOpen}
        close={e => { setOpen(false); }}
      />
    </Header>
  );
});
