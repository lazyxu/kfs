import React from 'react';

import styled from 'styled-components';
import Modal from 'components/Modal';

import { getConfig, setConfig } from 'adaptor/config';

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
  const [textarea, setTextarea] = React.useState({
    text: '',
    valid: false,
  });
  React.useEffect(() => {
    if (props.isOpen) {
      console.log('load config', JSON.stringify(getConfig(), undefined, 2));
      setTextarea({
        text: JSON.stringify(getConfig(), undefined, 2),
        valid: true,
      });
    }
  }, [props.isOpen]);
  return (
    <Modal
      title="镜像备份"
      {...props}
    >
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
    </Modal>
  );
});
