const fs = require('fs-extra');
const { spawn } = require('child_process');
const os = require('os');
const pathLib = require('path');

function runCommand(command, ...args) {
  return new Promise((resolve, reject) => {
    let message;
    const handler = spawn(command, ...args);
    // handler.stdout.on('data', (data) => {
    //   console.log(`stdout: ${data}`);
    // });

    handler.stderr.on('data', (data) => {
      console.log('stderr', data.toString());
      message = data.toString();
    });

    handler.on('exit', (code) => {
      console.log('exit', code);
      if (code === 0) {
        resolve();
      } else {
        reject({ code, message });
      }
    });
  });
}

function middleware(asyncFunc) {
  return (call, cb) => {
    asyncFunc(call).then((res) => cb(null, res)).catch((e) => {
      console.error(e);
      cb(e);
    });
  };
}

async function getFileStat(path, fileName) {
  fileName && (path = pathLib.join(path, fileName));
  fileName = pathLib.basename(path);
  const stat = await fs.stat(path);
  let type;
  stat.isDirectory() && (type = 'dir');
  stat.isFile() && (type = 'file');
  if (type) {
    return {
      name: fileName,
      type,
      size: stat.size,
      atimeMs: stat.atimeMs,
      mtimeMs: stat.mtimeMs,
      ctimeMs: stat.ctimeMs,
      birthtimeMs: stat.birthtimeMs,
    };
  }
  return undefined;
}

module.exports.ls = middleware(async (call) => {
  console.log('ls', call.request);
  let { path } = call.request;
  if (path.charAt(0) === '~') {
    path = os.homedir() + path.substring(1);
  }
  path = pathLib.normalize(path);
  global.mainWindow.setTitle(path);
  const fileNames = await fs.readdir(path);
  const files = await Promise.all(fileNames.map((fileName) => getFileStat(path, fileName)));
  return { path, files: files.filter((f) => f !== undefined) };
});

module.exports.mv = middleware(async (call) => {
  console.log('mv', call.request);
  const { src, dst } = call.request;
  await runCommand('mv', [...src, dst]);
  return {};
});

module.exports.createFile = middleware(async (call) => {
  console.log('createFile', call.request);
  let { path } = call.request;
  if (!path || path.length === 0) {
    const fileNameBase = pathLib.join(call.metadata.get('pwd')[0], 'new file');
    path = fileNameBase;
    let i = 0;
    if (await fs.pathExists(path)) {
      do {
        path = `${fileNameBase} (${++i})`;
      } while (await fs.pathExists(path));
    }
  }
  console.log('createFile', path);
  await fs.createFile(path);
  const stat = await getFileStat(path);
  return stat;
});


module.exports.mkdir = middleware(async (call) => {
  console.log('mkdir', call.request);
  let { path } = call.request;
  if (!path || path.length === 0) {
    const fileNameBase = pathLib.join(call.metadata.get('pwd')[0], 'new dir');
    path = fileNameBase;
    let i = 0;
    if (await fs.pathExists(path)) {
      do {
        path = `${fileNameBase} (${++i})`;
      } while (await fs.pathExists(path));
    }
  }
  console.log('mkdir', path);
  await fs.mkdir(path);
  const stat = await getFileStat(path);
  return stat;
});

module.exports.remove = middleware(async (call) => {
  console.log('remove', call.request);
  const { path } = call.request;
  const invalidPathList = path.filter((path) => !path || path === '/');
  if (invalidPathList.length !== 0) {
    throw new Error(`cannot remove these files: ${invalidPathList.join('\n')}`);
  }
  await Promise.all(path
    .filter((path) => path && path !== '/')
    .map((path) => fs.remove(path)));
  return {};
});

module.exports.download = middleware(async (call) => {
  console.log('download', call.request);
  const { path } = call.request;
  if (path.length === 1) {
    const filePath = path[0];
    const stat = fs.statSync(filePath);
    if (stat.isFile()) {
      const singleFileContent = fs.readFileSync(filePath).toString();
      console.log(singleFileContent);
      return { singleFileContent };
    }
  }
  const URI = [];
  return { URI };
});
