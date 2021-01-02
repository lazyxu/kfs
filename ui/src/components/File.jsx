import React from 'react';

import FileBase from 'components/FileBase';

export default React.memo(({
  ...props
}) => {
  return (
    <FileBase
      {...props}
      type="file"
    />
  );
});
