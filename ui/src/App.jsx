import React from 'react';
import styled from 'styled-components';

import ViewDefault from 'containers/GridView';
import Header from 'containers/Header';
import StatusBar from 'containers/StatusBar';

import './App.css';
import './_variables.scss';

import { cd } from 'bus/fs';
import { setState, busState, inState } from 'bus/bus';
import 'bus/triggers';
import { join } from 'utils/filepath';

const App = styled.div`
  height: 100%;
  width: 100%;
  position: fixed;
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
      <App
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
        <StyledHeader />
        <StyledViewDefault />
        <StyledStatusBar />
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
      </App>
    );
  }
}

export default component;
