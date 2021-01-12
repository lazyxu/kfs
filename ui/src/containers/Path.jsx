import React from 'react';
import styled from 'styled-components';

import PathArray from 'containers/PathArray';

import { ctxInState, StoreContext } from 'bus/bus';

@ctxInState(StoreContext, 'pwd')
class component extends React.Component {
  constructor(props, ctx) {
    super(props, ctx);

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
