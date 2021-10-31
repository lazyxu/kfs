const remote = require('@electron/remote');

const { Menu, MenuItem } = remote;

const menu = new Menu();
menu.append(new MenuItem({
  label: '放大',
  click: () => {
    console.log('item 1 clicked');
  },
}));
menu.append(new MenuItem({ type: 'separator' }));
menu.append(new MenuItem({ label: '缩小', type: 'checkbox', checked: true }));

export function onContextMenu(e) {
  e.preventDefault();
  menu.popup({ window: remote.getCurrentWindow() });
}
