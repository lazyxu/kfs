import React from 'react';

import PathSymbol from 'components/PathSymbol';

function* getParticalPath(path) {
  let lastPos = 0;
  let i;
  for (i = path.length - 1; i > 0; i--) {
    if (path.charAt(i) !== '/') {
      break;
    }
  }
  path = path.substring(0, i + 1);
  const { length } = path;
  for (i = 0; i < length; i++) {
    if (i > 0 && path.charAt(i - 1) === '/') {
      lastPos = i;
    }
    if (path.charAt(i + 1) === '/') {
      yield { path: path.substring(0, i + 1), symbol: path.substring(lastPos, i + 1) };
    }
    if (path.charAt(i) === '/') {
      yield { path: path.substring(0, i + 1), symbol: '/' };
    }
  }
  if (path.charAt(length - 1) !== '/') {
    yield { path, symbol: path.substring(lastPos, length) };
  }
}

function getPathArray(path) {
  const gen = getParticalPath(path);
  let res = gen.next();
  const array = [];
  while (!res.done) {
    array.push(res.value);
    res = gen.next(res.value);
  }
  let i = array.length - 1;
  if (i > 0) {
    array[i--].dragable = false;
  }
  if (i > 0) {
    array[i--].dragable = false;
  }
  for (; i >= 0; i--) {
    array[i].dragable = true;
  }
  return array;
}

class component extends React.Component {
  render() {
    const pathArray = getPathArray(this.props.pwd);
    return (
      <div>
        {pathArray.map((s) => (
          <PathSymbol
            key={s.path}
            {...s}
          />
        ))}
      </div>
    );
  }
}

export default component;
