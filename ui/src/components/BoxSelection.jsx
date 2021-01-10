import React from 'react';
import styled from 'styled-components';

const Bg = styled.div`
  position: absolute;
  height: 100%;
  width: 100%;
`;

export default React.memo(({
  onPosChange,
}) => {
  const [pos, setPos] = React.useState(undefined);
  const zIndex = React.useRef(0);
  const bg = React.useRef();
  let left;
  let top;
  let height;
  let width;
  if (pos) {
    if (pos.x1 < pos.x2) {
      left = pos.x1;
      width = pos.x2 - pos.x1;
    } else {
      left = pos.x2;
      width = pos.x1 - pos.x2;
    }
    if (pos.y1 < pos.y2) {
      top = pos.y1;
      height = pos.y2 - pos.y1;
    } else {
      top = pos.y2;
      height = pos.y1 - pos.y2;
    }
  }
  const Selection = typeof width === 'undefined'
    ? styled.div`
      display: none;
    `
    : styled.div`
      position: fixed;
      left: ${typeof left !== 'undefined' ? `${left}px` : ''};
      top: ${typeof top !== 'undefined' ? `${top}px` : ''};
      width: ${typeof width !== 'undefined' ? `${width}px` : ''};
      height: ${typeof height !== 'undefined' ? `${height}px` : ''};
      border: 1px solid #5d5e60;
      background-color: #343537;
      opacity: .7;
      z-index: var(--z-body-select);
    `;
  return (
    <Bg
      style={{ zIndex: zIndex.current }}
      onMouseDown={(e) => {
        if (e.button === 2) {
          return;
        }
        const { clientX, clientY } = e;
        setPos({ x1: clientX, y1: clientY });
        zIndex.current = 100;
      }}
      onMouseMove={(e) => {
        const { clientX, clientY } = e;
        setPos((prevState) => {
          if (prevState) {
            prevState.x2 = clientX;
            prevState.y2 = clientY;
            const { x, y } = bg.current.getBoundingClientRect();
            const x1 = prevState.x1 - x;
            const x2 = prevState.x2 - x;
            const y1 = prevState.y1 - y;
            const y2 = prevState.y2 - y;
            onPosChange({
              x1: Math.min(x1, x2),
              x2: Math.max(x1, x2),
              y1: Math.min(y1, y2),
              y2: Math.max(y1, y2),
            });
          }
          return prevState;
        });
      }}
      onMouseUp={(e) => {
        setPos(undefined);
        zIndex.current = 0;
      }}
      ref={bg}
    >
      <Selection />
    </Bg>
  );
});
