import React from 'react';

import FileIcon from 'components/FileIcon';
import { useClick } from 'use/use';

export default React.memo(function ({
  xlinkHref, style, onClick, onDoubleClick,
}) {
  return (
    <FileIcon
      xlinkHref={xlinkHref}
      onMouseDown={useClick(onClick, onDoubleClick)}
      style={style}
    />
  );
});
