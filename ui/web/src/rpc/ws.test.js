import { expect, test } from '@jest/globals'
import { list } from './ws';
import sysConfig from '../kfsConfig/default';

test('renders src/App.js', () => {
    return new Promise((resolve, reject) => {
        list(sysConfig, 'master', [], (total) => {
            console.log('total', total);
        }, (dirItem) => {
            console.log('dirItem', dirItem);
        }).then(code => {
            console.log('code', code);
            resolve(code);
        }).catch(e => {
            console.log('e', e);
            expect(e).toBeNull();
            reject(e);
        });
    });
});