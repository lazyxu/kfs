import { inState, setState } from 'bus/bus';
import Icon from 'components/Icon';
import React from 'react';
import styled from 'styled-components';
import Notifications from 'containers/Notifications';

const TaskBar = styled.div`
  background-color: var(--task-bar-color);
  height: var(--task-bar-height);
  z-index: var(--z-task-bar);
`;

@inState('showNotifications')
class component extends React.Component {
  render() {
    return (
      <TaskBar>
        <Icon
          icon="koala"
          color="#cccccc"
          size="1.3em"
          hoverColor="white"
          style={{ marginLeft: '1em' }}
        />
        <Icon
          icon="notice"
          color="#cccccc"
          size="1.3em"
          hoverColor="white"
          style={{ float: 'right', marginRight: '0.5em' }}
          onClick={() => this.setState(
            (prevState) => ({ showNotifications: !prevState.showNotifications }),
          )}
        />
        {this.state.showNotifications && <Notifications />}
      </TaskBar>
    );
  }
}

export default component;
