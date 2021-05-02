import React from 'react';

import FileIcon from 'components/FileIcon';

import { useClick } from 'use/use';

export default React.memo(function ({
  xlinkHref, onClick, onDoubleClick, ...props
}) {
  return (
    <FileIcon
      xlinkHref={xlinkHref}
      onMouseDown={useClick(onClick, onDoubleClick)}
      {...props}
    />
  );
});
