import React from 'react';
import styled from 'styled-components';

import { Notification } from 'kfs-components';

import { inState } from 'bus/bus';

const Notifications = styled.div`
  position: fixed;
  right: 0.5em;
  top: 1.5em;
  width: 20em;
`;

@inState('notifications')
class component extends React.Component {
  render() {
    return (
      <Notifications>
        {this.state.notifications.map((n) => (
          <Notification
            key={n.id}
            notification={n}
            remove={() => {
              this.setState((prevState) => ({
                notifications: prevState.notifications.filter((notification) => n !== notification),
              }));
            }}
          />
        ))}
      </Notifications>
    );
  }
}

export default component;
