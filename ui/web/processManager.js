let processes = [];

module.exports = {
  getProcesses: () => processes,
  addProcess: p => processes.push(p),
  removeProcess: p => {
    processes = processes.filter(v => v !== p);
  }
};
