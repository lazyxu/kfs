import React from 'react';

import FileName from 'components/FileName';

import { useClick } from 'use/use';

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
