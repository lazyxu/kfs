import React from 'react';

import FileBase from 'components/FileBase';

import { StoreContext } from 'bus/bus';

export default React.memo(({
  ...props
}) => {
  const context = React.useContext(StoreContext);
  return (
    <FileBase
      {...props}
      type="dir"
      onDoubleClick={() => context.cd(context.state.branch, props.path)}
    />
  );
});
