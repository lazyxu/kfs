import React from 'react';

import styled from 'styled-components';
import { upload } from 'bus/fs';
import { join } from 'utils/filepath';
import { busState } from 'bus/bus';
import myFileReader from 'utils/fileReader';
import map from 'promise.map';

const Div = styled.div`
`;
const Input = styled.input`
  opacity: 0;
  width: 0px;
  height: 0;
  display: none;
`;

const MAX_BLOCK_BYTES = 20 * 1024 * 1024;
const PARALLEL_BLOBK_COUNT = 5;
class component extends React.Component {
  componentDidMount() {
    this.input.addEventListener('change', async (e) => {
      console.log(this.input.files);
      const blob = this.input.files[0];
      const path = join(busState.pwd, blob.name);
      const total = blob.size;

      const bytes = await myFileReader(blob);
      upload(path, bytes);
      // if (total <= MAX_BLOCK_BYTES) {
      //   const bytes = await myFileReader(blob);
      //   upload(path, bytes);
      // } else {
      //   const blocks = [];
      //   const blockCount = total / MAX_BLOCK_BYTES;
      //   console.log('block count', blockCount);
      //   for (let i = 0; i < blockCount; i++) {
      //     blocks[i] = blob.slice(i * MAX_BLOCK_BYTES, (i + 1) * MAX_BLOCK_BYTES);
      //   }
      //   const hashList = await map(blocks, async (block, i) => {
      //     console.log(`---upload block ${i}/${blockCount}---`);
      //     const bytes = await myFileReader(block);
      //     const hash = await upload('', bytes);
      //     console.log(`---upload block ${i}/${blockCount} cb---`);
      //     return hash;
      //   }, 2);
      //   console.log('upload combined block');
      //   upload(path, undefined, hashList);
      // }
    });
  }

  render() {
    return [
      <Div
        key="0"
        onMouseDown={(e) => { e.stopPropagation(); }}
        onMouseUp={(e) => { e.stopPropagation(); }}
        onClick={(e) => { console.log(e); this.input.click(); }}
      >
        上传文件
      </Div>,
      <Input
        key="1"
        onClick={(e) => { console.log('file', e); }}
        ref={(input) => this.input = input}
        type="file"
      />,
    ];
  }
}

export default component;
