import React from 'react';

import FileName from './FileName.jsx';

import { useClick } from '../lib/use/index';

export default React.memo(({
  name, onClick, onDoubleClick, ...props
}) => {
  return (
    <FileName
      name={name}
      onMouseDown={useClick(onClick, onDoubleClick)}
      {...props}
    />
  );
});
