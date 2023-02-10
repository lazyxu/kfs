import './index.scss';
import Icon from "components/Icon/Icon";
import {useClick} from "use";
import {open} from "api/api";
import useResourceManager from 'hox/resourceManager';
import useSysConfig from 'hox/sysConfig';
import useContextMenu from "../../hox/contextMenu";

function downloadURI(uri, name) {
    let link = document.createElement("a");
    link.download = name;
    link.href = uri;
    link.click();
}

function downloader(data, name) {
    let blob = new Blob([data]);
    let url = window.URL.createObjectURL(blob);
    downloadURI(url, name);
    window.URL.revokeObjectURL(url);
}

export default ({name, type, filesElm}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    const [contextMenu, setContextMenu] = useContextMenu();
    const onClick = e => {
        console.log('onClick')
    }
    const onDoubleClick = e => {
        console.log('onDoubleClick')
    }
    const onOpen = name => {
        console.log(name);
        (async () => {
            let {filePath, branchName} = resourceManager;
            filePath = [...filePath, name];
            await open(sysConfig, setResourceManager, branchName, filePath);
        })()
    }
    return (
        <div className='file-normal' onContextMenu={(e) => {
            e.preventDefault();
            e.stopPropagation();
            const {clientX, clientY} = e;
            let {x, y, width, height} = filesElm.current.getBoundingClientRect();
            setContextMenu({
                type: 'file',
                name,
                clientX, clientY,
                x, y, width, height,
            })
        }}>
            <div onMouseDown={useClick(onClick, () => {
                onOpen(name);
            })}>
                <Icon icon={type === 'dir' ? 'floderblue' : 'file3'} className='file-icon'/>
            </div>
            <div className='file-name-wrapper'>
                <p kfs-attr="file" className='file-name' onMouseDown={useClick(onClick, onDoubleClick)}>{name}</p>
            </div>
        </div>
    )
};
