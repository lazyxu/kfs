import React from 'react';

import styled from 'styled-components';
import { busState, setState } from 'bus/bus';
import Icon from 'components/Icon';

const Modal = styled.div`
  position: fixed;
  display: grid;
  grid-template-columns: 1fr;
  grid-template-rows: auto 1fr auto;
  max-width: calc(100% - 100px);
  max-height: calc(100% - 100px);
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

export default React.memo(({
  isOpen, disableSave, save, close, children, title, content, ...props
}) => {
  const cur = React.useRef();
  const [pos, setPos] = React.useState({ top: 50, left: 50 });
  const [size, setSize] = React.useState({ width: 800, height: 600 });
  const onMove = e => {
    if (cur.current) {
      const { clientX, clientY } = e;
      const { clientWidth, clientHeight } = document.documentElement;
      if (clientX < 0 || clientX > clientWidth || clientY < 0 || clientY > clientHeight) {
        cur.current = undefined;
        document.removeEventListener('mousemove', onMove);
        return;
      }
      const offX = clientX - cur.current.clientX;
      const offY = clientY - cur.current.clientY;
      if (offX === 0 || offY === 0) {
        return;
      }
      setPos((prev) => {
        return {
          top: offY + prev.top,
          left: offX + prev.left,
        };
      });
      cur.current = { clientX, clientY };
    }
  };
  return (
    <Modal
      {...props}
      style={Object.assign(props.style || {}, {
        display: isOpen ? 'grid' : 'none',
        top: `${pos.top}px`,
        left: `${pos.left}px`,
        width: `${size.width}px`,
        height: `${size.height}px`,
      })}
      onKeyDown={e => { e.stopPropagation(); }}
    >
      <Header
        onMouseDown={e => {
          cur.current = { clientX: e.clientX, clientY: e.clientY };
          document.removeEventListener('mousemove', onMove);
          document.addEventListener('mousemove', onMove);
        }}
        onMouseUp={e => {
          if (cur.current) {
            document.removeEventListener('mousemove', onMove);
            cur.current = undefined;
          }
        }}
      >
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
});
