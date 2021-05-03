import React from 'react';
import styled from 'styled-components';

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

document.addEventListener('click', function () {
  Object.values(onFinishList).forEach(onFinish => {
    onFinish();
  });
}, false);

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
    <Div
      onMouseDown={(e) => e.stopPropagation()}
      onClick={e => { e.stopPropagation(); }}
    >
      {Object.keys(options).map((o) => {
        const option = options[o];
        let fn = option;
        if (option.fn) {
          fn = option.fn;
          if (option.enabled) {
            return (
              <DisabledOption key={o}>{o}</DisabledOption>
            );
          }
        }
        if (option?.type?.name === 'component') {
          return (
            <Option key={o}>{option}</Option>
          );
        }
        return (
          <Option
            key={o}
            onMouseDown={(e) => {
              fn(e); onFinish && onFinish(); e.stopPropagation();
            }}
          >
            {o}
          </Option>
        );
      })}
    </Div>
  );
};
