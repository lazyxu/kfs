import React from 'react';
import styled from 'styled-components';

import Dir from 'components/Dir';
import File from 'components/File';
import DefaultContextMenu from 'apps/ResourceManager/DefaultContextMenu';
import FileContextMenu from 'apps/ResourceManager/FileContextMenu';
import BoxSelection from 'components/BoxSelection';

import {
  inState, busState, setState, busValue, StoreContext,
} from 'bus/bus';
import { join } from 'utils/filepath';

const View = styled.div`
  position: relative;
  height: 100%;
  width: 100%;
  overflow: scroll;
  display: flex;
  flex-flow:row wrap;
  align-content:flex-start;
  background-color: transparent;
`;

@inState('files', 'chosen', 'boxChosen')
class component extends React.Component {
  static contextType = StoreContext;

  render() {
    const {
      files, chosen, boxChosen,
    } = this.state;
    return (
      <View
        onContextMenu={(e) => {
          e.preventDefault();
          const { fileListView } = busValue;
          console.log('onContextMenu', e.target, e.target.getAttribute('data-tag'));
          if (e.target === fileListView || e.target.getAttribute('data-tag') !== 'choose-able') {
            let { clientX, clientY } = e;
            (clientX > fileListView.clientWidth - 200) && (clientX = fileListView.clientWidth - 200);
            (clientY > fileListView.clientHeight - 120) && (clientY = fileListView.clientHeight - 120);
            setState({
              contextMenu: { x: clientX, y: clientY },
              contextMenuForFile: null,
            });
          }
        }}
        onMouseDown={(e) => {
          if (e.button === 2) {
            return;
          }
          if (e.target.getAttribute('data-tag') !== 'choose-able') {
            if (!e.metaKey) {
              setState({
                chosen: (_chosen) => {
                  Object.keys(_chosen).forEach((item) => {
                    delete _chosen[item];
                  });
                },
              });
            }
          }
        }}
        onMouseUp={(e) => {
          const keys = Object.keys(this.state.boxChosen);
          if (keys.length !== 0) {
            setState({
              chosen: (_chosen) => {
                keys.forEach((key) => _chosen[key] = 1);
              },
              boxChosen: {},
            });
          }
        }}
        ref={(fileListView) => {
          this.context.state.fileListView = fileListView;
        }}
      >
        <DefaultContextMenu />
        <FileContextMenu />
        <BoxSelection
          onPosChange={({
            x1, x2, y1, y2,
          }) => {
            const boxChosen = {};
            window.sss = this.context.state;
            const { childNodes, scrollTop } = this.context.state.fileListView;
            for (let i = 0; i < childNodes.length; i++) {
              const e = childNodes[i];
              const path = e.getAttribute('data-path');
              if (path) {
                if (!(e.offsetLeft + e.clientWidth < x1
                  || e.offsetTop + e.clientHeight - scrollTop < y1
                  || e.offsetLeft > x2
                  || e.offsetTop - scrollTop > y2)) {
                  boxChosen[path] = 1;
                }
              }
            }
            setState({ boxChosen });
          }}
        />
        {files.map((f) => {
          const path = join(busState.pwd, f.name);
          return f.type === 'file'
            ? <File key={`${f.type}-${path}`} {...f} chosen={chosen[path] || boxChosen[path]} dir={busState.pwd} path={path} />
            : <Dir key={`${f.type}-${path}`} {...f} chosen={chosen[path] || boxChosen[path]} dir={busState.pwd} path={path} />;
        })}
      </View>
    );
  }
}

export default component;
