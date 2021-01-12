import React from 'react';
import styled from 'styled-components';

import { warn } from 'bus/notification';
import { StoreContext } from 'bus/bus';

class component extends React.Component {
  static contextType = StoreContext

  state = {
    dragOver: false,
  }

  render() {
    const Button = styled.button`
      padding: 0.2em;
      color: white;
      background-color: transparent;
      border: 1px dashed transparent;
      border-color: ${this.state.dragOver && 'white'};
      border-radius: 0.3em;
      outline: none;
      :hover {
        border: 1px solid white;
        cursor: pointer;
      }
      display: inline-block;
    `;
    const { path, symbol, dragable } = this.props;
    return (
      <Button
        key={path}
        dragable={dragable}
        onClick={() => this.context.cd(path)}
        onDragEnter={(e) => {
          if (dragable) {
            e.preventDefault();
            if (!this.state.dragOver) {
            // console.log('onDragEnter', e.dataTransfer.getData('text/plain'));
              this.setState({ dragOver: true });
            }
          }
        }}
        onDragOver={(e) => {
          if (dragable) {
            e.preventDefault();
          }
        }}
        onDragLeave={(e) => {
          if (dragable) {
            e.preventDefault();
            if (this.state.dragOver) {
            // console.log('onDragLeave', e.dataTransfer.getData('text/plain'));
              this.setState({ dragOver: false });
            }
          }
        }}
        onDrop={(e) => {
          if (dragable) {
            let files = e.dataTransfer.getData('text/plain');
            files = JSON.parse(files);
            let i;
            for (i = path.length - 1; i > 0; i--) {
              if (path.charAt(i) !== '/') {
                break;
              }
            }
            const dst = path.substring(0, i + 1);
            console.log('onDrop', files, dst);
            if (files.includes(dst)) {
              warn('移动文件至文件夹', '移动文件夹至本身');
            } else {
              this.context.mv(files, dst);
            }
            if (this.state.dragOver) {
              this.setState({ dragOver: false });
            }
            e.preventDefault();
          }
        }}
      >
        {symbol}
      </Button>
    );
  }
}

export default component;
