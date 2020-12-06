import React from 'react';

import styled from 'styled-components';

import Icon from 'components/Icon';

const Notification = styled.div`
  animation: slideIn .8s linear;
  margin: 1em 0;
  padding: 0.5em;
  color: white;
  background-color: #303030;
  display: block;
  transition: opcity 300ms ease-in-out;
`;
const Header = styled.div`
  position: relative;
  width: 100%;
  display: flex;
  border-bottom: 1px solid #414141;
`;
const Body = styled.div`
  position: relative;
  width: 100%;
  padding: 0.5em;
  white-space: pre-wrap;
`;
const Title = styled.div`
  flex: 1;
  padding-left: 0.5em;
`;

class component extends React.Component {
  render() {
    return (
      <Notification>
        <Header>
          <Icon icon={this.props.notification.type} color="green" size="1em" />
          <Title>{this.props.notification.title}</Title>
          <Icon
            icon="close"
            color="white"
            size="1em"
            hoverCursor="pointer"
            onClick={() => this.props.remove()}
          />
        </Header>
        <Body>{this.props.notification.message && this.props.notification.message.split('%0A').join('\n')}</Body>
      </Notification>
    );
  }
}

export default component;
