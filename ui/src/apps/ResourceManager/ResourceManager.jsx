import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';
import ViewDefault from 'apps/ResourceManager/GridView';
import Header from 'apps/ResourceManager/Header';
import StatusBar from 'apps/ResourceManager/StatusBar';

import Store, { StoreContext } from 'bus/bus';

const App = styled.div`
  height: 100%;
  width: 100%;
  position: relative;
  color: var(--color);
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: auto 1fr auto;
  user-select: element;
  :focus {
    outline: none;
  }
`;
const StyledHeader = styled(Header)`
  background-color: red;
  grid-column: 1;
  grid-row: 1;
  z-index: var(--z-header);
`;
const StyledViewDefault = styled(ViewDefault)`
  grid-column: 1;
  grid-row: 2;
  z-index: var(--z-body);
`;
const StyledStatusBar = styled(StatusBar)`
  grid-column: 1;
  grid-row: 3;
  z-index: var(--z-footer);
`;

export default React.memo(({
  ...props
}) => {
  const context = React.useRef();
  if (!context.current) {
    // componentWillMount
    context.current = new Store({
      pwd: '/',
      files: [],
      notifications: [],
      chosen: {},
      fileSize: null,
      filesPositions: [],
      boxChosen: {},
      cutFiles: [],
      copyFiles: [],
    });
  }
  React.useEffect(() => {
    if (props.isOpen) {
      console.log(context.current);
      context.current.cd(context.current.pwd);
    }
  }, [props.isOpen]);
  return (
    <StoreContext.Provider value={context.current}>
      <Modal
        title="资源管理"
        {...props}
      >
        <App>
          <StyledHeader />
          <StyledViewDefault />
          <StyledStatusBar />
        </App>
      </Modal>
    </StoreContext.Provider>
  );
});
