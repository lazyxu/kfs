
import babel from '@rollup/plugin-babel';
import external from 'rollup-plugin-peer-deps-external';
import del from 'rollup-plugin-delete';
import pkg from './package.json';

export default {
    input: "components/index.js",
    output: [
        { file: 'dist/index.cjs.js', format: 'cjs' },
        { file: 'dist/index.esm.js', format: 'esm' }
    ],
    plugins: [
        external(),
        babel({
            exclude: 'node_modules/**',
            babelHelpers: 'bundled',
        }),
        del({ targets: ['dist/*'] }),
    ],
    external: Object.keys(pkg.peerDependencies || {}),
};
