import React from 'react';

import FileBase from 'components/FileBase';
import { cd } from 'bus/fs';

export default React.memo(({
  ...props
}) => {
  return (
    <FileBase
      {...props}
      type="dir"
      onDoubleClick={() => cd(props.path)}
    />
  );
});
