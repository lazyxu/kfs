import React from 'react';
import styled from 'styled-components';

import App from 'components/App';

import ResourceManager from 'apps/ResourceManager/ResourceManager';
import ConfigEditor from 'apps/SystemConfig';

import './App.css';
import './_variables.scss';

import { inState } from 'bus/bus';
import 'bus/triggers';
import TaskBar from 'containers/TaskBar';

import 'adaptor/ws';
import MirrorSync from 'apps/MirrorBackup';

const Index = styled.div`
  height: 100%;
  width: 100%;
  position: fixed;
  color: var(--color);
  display: flex;
  flex-direction: column;
`;
const Desktop = styled.div`
  height: 100%;
  width: 100%;
  position: relative;
  color: var(--color);
  display: flex;
  flex-direction: column;
`;

@inState('windows')
class component extends React.Component {
  render() {
    return (
      <Index>
        <TaskBar />
        <Desktop>
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
          <App
            icon="shangchuan"
            color="#cccccc"
            elm={MirrorSync}
            newWindowOption="true"
            text="镜像备份"
          />
          {/* <App
            icon="yuntongbu"
            color="#cccccc"
            elm={MirrorSync}
            newWindowOption="true"
            text="同步盘"
          /> */}
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
      </Index>
    );
  }
}

export default component;
