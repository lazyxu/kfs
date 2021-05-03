const fs = require('fs-extra');

fs.writeFileSync('components/index.js', fs.readdirSync('components')
  .filter(filename => /(\.\/|\.jsx)/g.test(filename))
  .map(filename => filename.replace(/(\.\/|\.jsx)/g, ""))
  .map(name => `export { default as ${name} } from './${name}.jsx';`)
  .reduce((prev, cur) => prev + '\n' + cur, ''))
