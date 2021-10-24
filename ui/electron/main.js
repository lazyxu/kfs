// eslint-disable-next-line import/no-extraneous-dependencies
const { app, BrowserWindow, dialog, ipcMain } = require('electron');
const path = require('path');
const { getProcesses } = require('./processManager');

let mainWindow;

const { nativeImage } = require('electron');
const image = nativeImage.createFromPath(path.join(__dirname,
  process.env.ELECTRON_START_URL ? '../desktop/public/icon512.png' : 'public/icon512.png'));

app.setName("考拉云盘");
app.dock.setIcon(image);

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 800,
    height: 600,
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
  global.mainWindow = mainWindow;

  mainWindow.loadURL(process.env.ELECTRON_START_URL || require('url').format({
    pathname: path.join(__dirname, 'public/index.html'),
    protocol: 'file:',
    slashes: true
  }));

  const remoteMain = require('@electron/remote/main');
  remoteMain.initialize();
  remoteMain.enable(mainWindow.webContents);

  const { app, Menu } = require('electron');

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
