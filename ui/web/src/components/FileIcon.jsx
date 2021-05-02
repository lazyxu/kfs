import React from 'react';
import styled from 'styled-components';

const Icon = styled.svg`
  height: 100%;
  width: 100%;
  vertical-align: -0.15em;
  fill: #dddddd;
  background-color: transparent;
`;
const FileIcon = styled.div`
  padding: 0.3em;
  margin: 0 0.3em;
  height: 4em;
  width: cal(100% - 0.3em);
  border-radius: 0.3em;
`;

export default React.memo(({
  xlinkHref, ...props
}) => {
  return (
    <FileIcon
      data-tag="choose-able"
      {...props}
    >
      <Icon data-tag="choose-able" aria-hidden="true">
        <use data-tag="choose-able" xlinkHref={xlinkHref} />
      </Icon>
    </FileIcon>
  );
});
