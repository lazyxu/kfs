import React from 'react';

import styled from 'styled-components';
import { busState, setState } from 'bus/bus';
import Icon from 'components/Icon';

const Modal = styled.div`
  position: fixed;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: auto 1fr auto;
  max-width: calc(100% - 5em);
  max-height: calc(100% - 5em);
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: auto;
  border-radius: var(--modal-radius);
  border: 1px solid #717171;
  background: var(--modal-body-color);
  overflow: hidden;
  z-index: 9999;
  box-shadow: inset 0 0 1px 0 #000;
`;
const Empty = styled.div`
  flex: 1;
`;
const Header = styled.div`
  height: var(--modal-header-height);
  width: calc(100% - 0.5em);
  background-color: var(--modal-header-color);
  display: flex;
  border-radius: var(--modal-radius) var(--modal-radius) 0 0;
  padding-left: 0.5em;
`;
const Body = styled.div`
  padding: 4px;
  overflow: hidden;
`;
const Text = styled.div`
  text-align: center;
`;

export default class component extends React.Component {
  render() {
    const {
      isOpen, disableSave, save, close, children, title, content, ...props
    } = this.props;
    return (
      <Modal
        {...props}
        style={Object.assign(props.style || {}, { display: isOpen ? 'grid' : 'none' })}
        onKeyDown={e => { e.stopPropagation(); }}
      >
        <Header>
          <Text>
            {title}
          </Text>
          <Empty />
          <Icon
            icon="close"
            color="#cccccc"
            size="1.5em"
            hoverColor="white"
            hoverCursor="pointer"
            onClick={() => { close(); }}
          />
        </Header>
        <Body>
          {children}
        </Body>
      </Modal>
    );
  }
}
