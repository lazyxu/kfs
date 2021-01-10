import React from 'react';
import styled from 'styled-components';

import Icon from 'components/Icon';
import App from 'components/App';

import ResourceManager from 'apps/ResourceManager/ResourceManager';
import ConfigEditor from 'apps/SystemConfig';

import './App.css';
import './_variables.scss';

import { cd } from 'bus/fs';
import { setState, busState, inState } from 'bus/bus';
import 'bus/triggers';
import { join } from 'utils/filepath';

const Desktop = styled.div`
  height: 100%;
  width: 100%;
  position: fixed;
  color: var(--color);
  display: flex;
  flex-direction: column;
  user-select: element;
  :focus {
    outline: none;
  }
`;

@inState('windows')
class component extends React.Component {
  state = {
    loaded: false,
  }

  componentDidMount() {
    cd('/').then(() => this.setState({ loaded: true }));
  }

  render() {
    return (
      <Desktop
        tabIndex="-1"
        onKeyDown={(e) => {
          console.log(e.keyCode, e.metaKey);
          if (e.keyCode === 65 && e.metaKey === true) {
            setState({
              chosen: (_chosen) => {
                busState.files.map((f) => join(busState.pwd, f.name)).forEach((path) => {
                  _chosen[path] = 1;
                });
              },
            });
            e.preventDefault();
          }
        }}
        onMouseDown={(e) => e.button !== 2
          && setState({ contextMenu: null, contextMenuForFile: null })}
      >
        {!this.state.loaded && <span>loading...</span>}
        <App
          icon="wangpan"
          color="#cccccc"
          elm={ResourceManager}
          text="资源管理"
        />
        <App
          icon="peizhi"
          color="#cccccc"
          elm={ConfigEditor}
          newWindowOption="true"
          text="系统配置"
        />
        {Object.values(this.state.windows).sort((a, b) => a.zIndex - b.zIndex).map(w => {
          const Window = w.elm;
          return (
            <Window
              key={w.id}
              id={w.id}
              isOpen="true"
              close={() => {
                this.setState(prev => {
                  delete prev.windows[w.id];
                  return { windows: prev.windows };
                });
              }}
              zIndex={w.zIndex}
            />
          );
        })}
      </Desktop>
    );
  }
}

export default component;
