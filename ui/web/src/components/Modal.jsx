import React from 'react';

import { Modal } from 'kfs-components';

import { busState, setState } from 'bus/bus';

export default React.memo(({
  id, ...props
}) => {
  return (
    <Modal
      onClick={e => {
        setState({
          windows: (windows => {
            fucusWindow(windows, id);
          }),
        });
      }}
      {...props}
    />
  );
});

function fucusWindow(windows, id) {
  const cur = windows[id];
  const vals = Object.values(windows).sort((a, b) => b.zIndex - a.zIndex);
  if (vals.length === 1) {
    return;
  }
  if (cur === vals[0]) {
    return;
  }
  const maxZIndex = vals[0].zIndex;
  for (let i = 1; i < vals.length; i++) {
    vals[i - 1].zIndex = vals[i].zIndex;
    if (vals[i] === cur) {
      break;
    }
  }
  cur.zIndex = maxZIndex;
}
let id = 0;
busState.windows = {};
export function newWindow(elm, single = false) {
  setState({
    windows: (windows => {
      if (single) {
        const w = Object.values(windows).filter(w => w.elm === elm)[0];
        if (w) {
          fucusWindow(windows, w.id);
          return;
        }
      }
      id++;
      windows[id] = {
        elm,
        id,
        zIndex: id,
      };
    }),
  });
  return id;
}
