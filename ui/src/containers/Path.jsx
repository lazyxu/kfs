import React from 'react';

import styled from 'styled-components';

import { inState } from 'bus/bus';
import PathArray from 'containers/PathArray';

@inState('pwd')
class component extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      type: 'plain',
    };
  }

  render() {
    const Path = styled.div`
      padding: 0 0.2em;
      margin: 0.3em;
      /* border: 1px solid #31729c;
      border-radius: 0.5em; */
      background-color: var(--header-color);
      /* width: 10em; */
      overflow: scroll;
    `;
    return (
      <Path>
        {this.state.type === 'plain' && <PathArray pwd={this.state.pwd} />}
      </Path>
    );
  }
}

export default component;
