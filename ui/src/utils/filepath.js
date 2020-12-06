export function dirname(path) {
  const index = path.lastIndexOf('/');
  return path.substring(0, index) || '/';
}

export function basename(path) {
  const index = path.lastIndexOf('/');
  return path.substring(index + 1);
}

export function join(...elems) {
  return elems.join('/').replace(/\/+/g, '/');
}
