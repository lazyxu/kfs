import React from 'react';

import styled from 'styled-components';

const List = styled.ul`
  list-style-type: none;
  padding-left: 1em;
`;
const ListItem = styled.li`
`;
const File = styled.div`
  padding: 0.15em;
  cursor: pointer;
  :hover {
    background: #2f2f2f;
  }
  :focus {
    background: #094771;
    outline: none;
  }
`;

export default React.memo(({
  ...props
}) => {
  return (
      <List>
        <ListItem>
          <File>workspace</File>
        </ListItem>
        <ListItem>
          <List>
            <ListItem>
              <File>index.js</File>
            </ListItem>
            <ListItem>
              <File>index.js</File>
            </ListItem>
          </List>
        </ListItem>
        <ListItem>
          <File>repos</File>
        </ListItem>
      </List>
  );
});
