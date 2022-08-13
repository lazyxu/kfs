// eslint-disable-next-line import/no-extraneous-dependencies
const { app, BrowserWindow, dialog, ipcMain, Menu } = require('electron');

if (app.isPackaged) {
  process.env.NODE_ENV = 'production';
  process.env.REACT_APP_PLATFORM = 'electron';
}

if (process.env.NODE_ENV === 'development') {
  require('electron-reloader')(module);
}


const path = require('path');
const { getProcesses } = require('./processManager');

const publicPath = 'electron-' + process.env.NODE_ENV;
let mainWindow;

app.setName("考拉云盘");

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1280,
    height: 800,
    // title: "考拉云盘",
    titleBarStyle: 'hidden',
    // backgroundColor: '#FFF',
    webPreferences: {
      // preload: path.join(publicPath, 'preload.js'),
      nodeIntegration: true,
      contextIsolation: false,
      nativeWindowOpen: true,
      remote: true,
      sandbox: false,
      nodeIntegrationInSubFrames: true, //for subContent nodeIntegration Enable
      // webviewTag:true //for webView
    },
    // icon: image,
  });
  mainWindow.webContents.openDevTools();

  global.mainWindow = mainWindow;

  mainWindow.loadFile(path.join(publicPath, 'index.html'));

  const isMac = process.platform === 'darwin';

  const template = [
    // { role: 'appMenu' }
    ...(isMac ? [{
      label: app.name,
      submenu: [
        { role: 'about' },
        { type: 'separator' },
        { role: 'services' },
        { type: 'separator' },
        { role: 'hide' },
        { role: 'hideOthers' },
        { role: 'unhide' },
        { type: 'separator' },
        { role: 'quit' }
      ]
    }] : []),
    // { role: 'fileMenu' }
    {
      label: 'File',
      submenu: [
        isMac ? { role: 'close' } : { role: 'quit' }
      ]
    },
    // { role: 'editMenu' }
    {
      label: 'Edit',
      submenu: [
        { role: 'undo' },
        { role: 'redo' },
        { type: 'separator' },
        { role: 'cut' },
        { role: 'copy' },
        { role: 'paste' },
        ...(isMac ? [
          { role: 'pasteAndMatchStyle' },
          { role: 'delete' },
          { role: 'selectAll' },
          { type: 'separator' },
          {
            label: 'Speech',
            submenu: [
              { role: 'startSpeaking' },
              { role: 'stopSpeaking' }
            ]
          }
        ] : [
          { role: 'delete' },
          { type: 'separator' },
          { role: 'selectAll' }
        ])
      ]
    },
    // { role: 'viewMenu' }
    {
      label: 'View',
      submenu: [
        { role: 'reload' },
        { role: 'forceReload' },
        { role: 'toggleDevTools' },
        { type: 'separator' },
        { role: 'resetZoom' },
        { role: 'zoomIn' },
        { role: 'zoomOut' },
        { type: 'separator' },
        { role: 'togglefullscreen' }
      ]
    },
    // { role: 'windowMenu' }
    {
      label: 'Window',
      submenu: [
        { role: 'minimize' },
        { role: 'zoom' },
        ...(isMac ? [
          { type: 'separator' },
          { role: 'front' },
          { type: 'separator' },
          { role: 'window' }
        ] : [
          { role: 'close' }
        ])
      ]
    },
    {
      role: 'help',
      submenu: [
        {
          label: 'Learn More',
          click: async () => {
            const { shell } = require('electron')
            await shell.openExternal('https://electronjs.org')
          }
        }
      ]
    }
  ]

  const menu = Menu.buildFromTemplate(template)
  Menu.setApplicationMenu(menu)

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

app.whenReady().then(() => {
  createWindow()
});

app.on('window-all-closed', () => {
  getProcesses().forEach(p => p.kill());
  app.quit();
});

ipcMain.on('select-dirs', async (event, arg) => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openDirectory']
  })
  console.log('directories selected', result.filePaths);
});
