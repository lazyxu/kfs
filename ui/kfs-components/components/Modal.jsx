import React from 'react';
import styled from 'styled-components';

import Icon from './Icon.jsx';

const minWidth = 100;
const minHeight = 100;
const Modal = styled.div`
  position: absolute;
  display: grid;
  grid-template-columns: 3px 1fr 3px ;
  grid-template-rows: 3px 1fr 3px;
  min-width: ${minWidth}px;
  min-height: ${minHeight}px;
  max-width: 100%;
  max-height: 100%;
  overflow: hidden;
  z-index: 9999;
`;
const Empty = styled.div`
  flex: 1;
`;
const App = styled.div`
  height: 100%;
  width: 100%;
  border: 1px solid #717171;
  border-radius: var(--modal-radius);
  box-shadow: inset 0 0 1px 0 #000;
  background: var(--modal-body-color);
  display: grid;
  grid-template-rows: auto 1fr;
  overflow: hidden;
`;
const Header = styled.div`
  height: var(--modal-header-height);
  width: calc(100% - 0.5em);
  background-color: var(--modal-header-color);
  display: flex;
  padding-left: 0.5em;
`;
const Body = styled.div`
  overflow: hidden;
`;
const Text = styled.div`
  text-align: center;
`;
const RowResizeable = styled.div`
  height: 100%;
  width: 100%;
  background-color: transparent;
  :hover {
    cursor: row-resize;
  }
`;
const ColResizeable = styled.div`
  height: 100%;
  width: 100%;
  background-color: transparent;
  :hover {
    cursor: col-resize;
  }
`;
const NWSEResizeable = styled.div`
  height: 100%;
  width: 100%;
  background-color: transparent;
  :hover {
    cursor: nwse-resize;
  }
`;
const NESWResizeable = styled.div`
  height: 100%;
  width: 100%;
  background-color: transparent;
  :hover {
    cursor: nesw-resize;
  }
`;

export default React.memo(({
  zIndex, isOpen, disableSave, save, close, children, title, content, onClick, ...props
}) => {
  const cur = React.useRef();
  const [pos, setPos] = React.useState({
    top: 0, left: 0, height: 600, width: 800,
  });
  const listeners = React.useRef([]);
  const onMouseDown = fn => e => {
    cur.current = { clientX: e.clientX, clientY: e.clientY };
    document.addEventListener('mousemove', fn);
  };
  const onMouseMove = (e, fn) => {
    if (cur.current) {
      const { clientX, clientY } = e;
      const { clientX: curX, clientY: curY } = cur.current;
      const { clientWidth, clientHeight } = document.documentElement;
      if (clientX < 0 || clientX > clientWidth || clientY < 0 || clientY > clientHeight) {
        return;
      }
      const offX = clientX - curX;
      const offY = clientY - curY;
      if (offX === 0 || offY === 0) {
        return;
      }
      fn(offX, offY);
      cur.current = { clientX, clientY };
    }
  };
  const onMove = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        return {
          ...prev,
          top: offY + prev.top,
          left: offX + prev.left,
        };
      });
    });
  };
  const topResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        if (prev.height - offY < minHeight) {
          return prev;
        }
        return {
          ...prev,
          top: prev.top + offY,
          height: prev.height - offY,
        };
      });
    });
  };
  const bottomResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        return {
          ...prev,
          height: prev.height + offY,
        };
      });
    });
  };
  const leftResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        if (prev.width - offX < minHeight) {
          return prev;
        }
        return {
          ...prev,
          left: prev.left + offX,
          width: prev.width - offX,
        };
      });
    });
  };
  const rightResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        return {
          ...prev,
          width: prev.width + offX,
        };
      });
    });
  };
  const topLeftResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        if (prev.width - offX < minHeight || prev.height - offY < minHeight) {
          return prev;
        }
        return {
          top: prev.top + offY,
          height: prev.height - offY,
          left: prev.left + offX,
          width: prev.width - offX,
        };
      });
    });
  };
  const topRightResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        if (prev.height - offY < minHeight) {
          return prev;
        }
        return {
          ...prev,
          top: prev.top + offY,
          height: prev.height - offY,
          width: prev.width + offX,
        };
      });
    });
  };
  const bottomLeftResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        if (prev.height - offY < minHeight) {
          return prev;
        }
        return {
          ...prev,
          height: prev.height + offY,
          left: prev.left + offX,
          width: prev.width - offX,
        };
      });
    });
  };
  const bottomRightResize = e => {
    return onMouseMove(e, (offX, offY) => {
      setPos((prev) => {
        return {
          ...prev,
          height: prev.height + offY,
          width: prev.width + offX,
        };
      });
    });
  };
  listeners.current.push(onMove, topResize, bottomResize, leftResize, rightResize,
    topLeftResize, topRightResize, bottomLeftResize, bottomRightResize);
  const clearListeners = e => {
    if (cur.current) {
      listeners.current.forEach(l => {
        document.removeEventListener('mousemove', l);
      });
      cur.current = undefined;
    }
  };
  React.useEffect(() => {
    document.addEventListener('mouseup', clearListeners);
    return () => {
      document.removeEventListener('mouseup', clearListeners);
    };
  }, []);
  return (
    <Modal
      {...props}
      style={Object.assign(props.style || {}, {
        display: isOpen ? 'grid' : 'none',
        zIndex,
        top: `${pos.top}px`,
        left: `${pos.left}px`,
        width: `${pos.width}px`,
        height: `${pos.height}px`,
      })}
      onClick={onClick}
      onKeyDown={e => { e.stopPropagation(); }}
    >
      <NWSEResizeable onMouseDown={onMouseDown(topLeftResize)} />
      <RowResizeable onMouseDown={onMouseDown(topResize)} />
      <NESWResizeable onMouseDown={onMouseDown(topRightResize)} />
      <ColResizeable onMouseDown={onMouseDown(leftResize)} />
      <App>
        <Header onMouseDown={onMouseDown(onMove)}>
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
      </App>
      <ColResizeable onMouseDown={onMouseDown(rightResize)} />
      <NESWResizeable onMouseDown={onMouseDown(bottomLeftResize)} />
      <RowResizeable onMouseDown={onMouseDown(bottomResize)} />
      <NWSEResizeable onMouseDown={onMouseDown(bottomRightResize)} />
    </Modal>
  );
});
