import React from 'react';
import ContextMenu from './index';
import useContextMenu from "../../hox/contextMenu";
import useSysConfig from "../../hox/sysConfig";
import useResourceManager from "../../hox/resourceManager";
import {open} from "../../api/api";

export default () => {
  const [contextMenu, setContextMenu] = useContextMenu();
  const [resourceManager, setResourceManager] = useResourceManager();
  const { sysConfig } = useSysConfig();
  console.log(contextMenu)
  if (contextMenu === null || contextMenu.type !== "file") {
    return <div />
  }
  let left = contextMenu.clientX;
  let top = contextMenu.clientY;
  let maxWidth = 200;
  let maxHeight = 100;
  if (left + maxWidth > contextMenu.x + contextMenu.width) {
    left = contextMenu.x + contextMenu.width - maxWidth;
  }
  if (top + maxHeight > contextMenu.y + contextMenu.height) {
    top = contextMenu.y + contextMenu.height - maxHeight;
  }
  let { filePath, branchName } = resourceManager;
  filePath.push(contextMenu.name);
  return <ContextMenu
    left={left}
    top={top}
    options={{
      打开: () => {
        open(sysConfig, setResourceManager, branchName, filePath);
      },
      // 下载: () => this.context.download(pathList),
      // 剪切: () => setState({
      //   clipboard: {
      //     cut: true,
      //     file: {
      //       branch: this.context.state.branch,
      //       pathList,
      //     },
      //   },
      // }),
      // 复制: () => setState({
      //   clipboard: {
      //     copy: true,
      //     file: {
      //       branch: this.context.state.branch,
      //       pathList,
      //     },
      //   },
      // }),
      // 删除: () => this.context.remove(pathList),
      // 重命名: () => this.context.setState({ chosen: (_chosen) => _chosen[pathList[0]] = 2 }),
      属性: { enabled: false, fn: () => { } },
      历史版本: { enabled: false, fn: () => { } },
    }}
    onFinish={() => setContextMenu(null)}
  />
}
