import React from 'react';

import FileBase from 'components/FileBase';
import { cd } from 'bus/fs';

import { busState } from 'bus/bus';
import { join } from 'utils/filepath';

export default React.memo(({
  ...props
}) => {
  return (
    <FileBase
      {...props}
      type="dir"
      onDoubleClick={() => cd(join(busState.pwd, props.name))}
    />
  );
});
