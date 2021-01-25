import React from 'react';
import styled from 'styled-components';

import Icon from 'components/Icon';
import { newWindow } from 'components/Modal';

const App = styled.div`
  margin-top: 1em;
  margin-left: 1em;
  width: 4em;
  text-align: center;
`;
const Text = styled.div`
  font-size: 1em;
  padding: 0;
  overflow : hidden;
  text-overflow: ellipsis;
  border-radius: 0.3em;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  user-select: none;
  overflow-wrap: break-word;
  margin: 0;
`;

export default ({
  text, elm, icon, color, newWindowOption,
}) => {
  return (
    <App>
      <Icon
        icon={icon}
        color={color}
        size="4em"
        hoverColor="white"
        hoverCursor="pointer"
        onClick={e => {
          newWindow(elm, newWindowOption);
        }}
      />
      <Text>
        {text}
      </Text>
    </App>
  );
};
