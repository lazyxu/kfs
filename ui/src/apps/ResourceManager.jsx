import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';
import ViewDefault from 'containers/GridView';
import Header from 'containers/Header';
import StatusBar from 'containers/StatusBar';

import { getConfig, setConfig } from 'adaptor/config';

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
  const [textarea, setTextarea] = React.useState({
    text: '',
    valid: false,
  });
  React.useEffect(() => {
    if (props.isOpen) {
      console.log('load config', JSON.stringify(getConfig(), undefined, 2));
      setTextarea({
        text: JSON.stringify(getConfig(), undefined, 2),
        valid: true,
      });
    }
  }, [props.isOpen]);
  return (
    <Modal
      title="资源管理器"
      {...props}
    >
      <App>
        <StyledHeader />
        <StyledViewDefault />
        <StyledStatusBar />
      </App>
    </Modal>
  );
});
