import React from 'react';
import styled from 'styled-components';

import UploadFile from 'components/UploadFile';
import { isObject } from 'util';

const Option = styled.div`
  color: #ffffff;
  padding: 0.2em 5em 0.2em 1.5em;
  cursor: default;
  :last-child{
    margin-bottom: 0;
  }
  :hover {
    background-color: #1367cd;
  }
`;
const DisabledOption = styled.div`
  color: gray;
  padding: 0.2em 5em 0.2em 1.5em;
  cursor: default;
  :last-child{
    margin-bottom: 0;
  }
`;

const onFinishList = {};
let id = 0;

document.onclick = e => {
  Object.values(onFinishList).forEach(onFinish => {
    onFinish();
  });
};

export default ({
  x, y, options, onFinish,
}) => {
  React.useEffect(() => {
    id++;
    onFinishList[id] = onFinish;
    return () => {
      delete onFinishList[id];
    };
  }, []);
  const Div = styled.div`
      background-color: #322f32;
      border: 1px solid #504d51;
      border-radius: 0.5em;
      padding: 0.5em 0;
      position: fixed;
      left: ${x}px;
      top: ${y}px;
      z-index: 100000;
      user-select: none;
    `;
  return (
    <Div onMouseDown={(e) => e.stopPropagation()}>
      {Object.keys(options).map((o) => {
        const option = options[o];
        let fn;
        let enabled = true;
        if (typeof option === 'function') {
          fn = option;
        } else {
          enabled = option.enabled;
          fn = option.fn;
        }
        return enabled ? (
          <Option
            key={o}
            onMouseDown={(e) => { fn(e); onFinish && onFinish(); e.stopPropagation(); }}
          >
            {o === '上传文件'
              ? <UploadFile text={o} />
              : o}
          </Option>
        ) : <DisabledOption key={o}>{o}</DisabledOption>;
      })}
    </Div>
  );
};
