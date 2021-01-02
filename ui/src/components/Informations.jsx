import React from 'react';

import styled from 'styled-components';

const Informations = styled.div`
  position: fixed;
  right: 0.5em;
  bottom: 2em;
  width: 20em;
`;
const Information = styled.div`
  border-bottom: 1px solid #414141;
  padding: 0.5em;
  color: white;
  background-color: #303030;
  display: block;
`;
export default React.memo(({
  informations,
}) => {
  return (
    <Informations>
      {informations.map((info, i) => (
        <Information key={i}>
          {info}
        </Information>
      ))}
    </Informations>
  );
});
