import { post } from './axios';

export async function readObject(hash) {
  return post('/api/readObject', {
    hash,
  });
}
