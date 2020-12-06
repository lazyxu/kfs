import React from 'react';

import FileBase from 'components/FileBase';

class component extends React.Component {
  render() {
    return (
      <FileBase
        {...this.props}
        type="file"
      />
    );
  }
}

export default component;
