let processes = [];

export const getProcesses = () => processes;
export const addProcess = p => processes.push(p);
export const removeProcess = p => {
  processes = processes.filter(v => v !== p);
};
