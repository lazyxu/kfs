import React from 'react';

import FileName from 'components/FileName';
import { useClick } from 'use/use';

export default React.memo(({
  name, style, onClick, onDoubleClick,
}) => {
  return (
    <FileName
      name={name}
      attributes={{
        onMouseDown: useClick(onClick, onDoubleClick),
        style,
      }}
    />
  );
});
