import { post } from './axios';

export async function readDirectory(hash) {
  return post('/api/readDirectory', {
    hash,
  });
}

export async function createDirectory(parent, name) {
  return post('/api/createDirectory', {
    parent, name,
  });
}
