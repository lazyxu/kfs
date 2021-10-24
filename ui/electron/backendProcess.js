const path = require('path');
const { spawn } = require('child_process');
const { addProcess, getProcesses, removeProcess } = require('./processManager');

const cwd = path.join(__dirname, 'public', 'extraResources');

const spawnfile = './kfs-client';

module.exports = (port) => {
  const list = getProcesses().filter(p => p.spawnfile === spawnfile && p.spawnargs.length === 2);
  console.log(`clear ${list.length} old backend process!`);
  const otherList = list.filter(p => p.spawnargs[1] !== port);
  otherList.forEach(p => {
    removeProcess(p);
    p.kill();
  });
  console.log(`clear ${list.length} old backend process finished!`);
  if (list.length > otherList.length) {
    console.log(`backend process port not changed ${port}!`);
    return;
  }
  console.log(`backend process start with port ${port}!`);
  const backendProcess = spawn(spawnfile, [port], { cwd });
  addProcess(backendProcess);
  backendProcess.stdout.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });

  backendProcess.stderr.on('data', (data) => {
    console.error(`stderr: ${data}`);
    removeProcess(backendProcess);
  });

  backendProcess.on('close', (code) => {
    console.log(`backendProcess exited with code ${code}`);
    removeProcess(backendProcess);
  });
};
