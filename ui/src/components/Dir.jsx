import React from 'react';

import FileBase from 'components/FileBase';
import { cd } from 'bus/fs';

import { busState } from 'bus/bus';
import { join } from 'utils/filepath';

class component extends React.Component {
  state = {}

  render() {
    return (
      <FileBase
        {...this.props}
        type="dir"
        onDoubleClick={() => cd(join(busState.pwd, this.props.name))}
      />
    );
  }
}

export default component;
