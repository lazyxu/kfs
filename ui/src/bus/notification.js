import { setState } from 'bus/bus';

let id = 0;
export default async function notification(type, title, message) {
  setState({
    showNotifications: true,
    notifications: (notifications) => {
      notifications.push({
        type, id: id++, title, message,
      });
      console.log('---notifications---', notifications);
      return { notifications };
    },
  });
}

export async function error(title, message) {
  notification('error', title, message);
}

export async function warn(title, message) {
  notification('warn', title, message);
}

export async function info(title, message) {
  notification('info', title, message);
}
