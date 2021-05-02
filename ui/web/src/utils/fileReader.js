export default function (blob) {
  return new Promise((resolve, reject) => {
    const r = new FileReader();
    r.onload = (e) => {
      resolve(new Uint8Array(e.target.result));
    };
    r.onerror = (e) => {
      reject(e);
    };
    r.readAsArrayBuffer(blob);
  });
}
