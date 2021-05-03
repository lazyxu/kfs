import React from 'react';

import { FileIcon } from 'kfs-components';

import { useClick } from '../lib/use/index';

export default React.memo(function ({
  onClick, onDoubleClick, ...props
}) {
  return (
    <FileIcon
      onMouseDown={useClick(onClick, onDoubleClick)}
      {...props}
    />
  );
});
