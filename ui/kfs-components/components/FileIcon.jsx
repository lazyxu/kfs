import React from 'react';
import styled from 'styled-components';
import Icon from './Icon.jsx';

const FileIcon = styled.div`
  padding: 0.3em;
  margin: 0 0.3em;
  height: 4em;
  width: cal(100% - 0.3em);
  border-radius: 0.3em;
`;

export default function ({
  ...props
}) {
  return (
    <FileIcon>
      <Icon {...props}/>
    </FileIcon>
  );
}
