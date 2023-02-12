import './index.scss';
import Icon from "components/Icon/Icon";
import {useClick} from "use";
import {open} from "api/api";
import {modeIsDir} from "api/utils/api";
import useResourceManager from 'hox/resourceManager';
import useSysConfig from 'hox/sysConfig';
import useContextMenu from "../../hox/contextMenu";

export default ({dirItem, filesElm}) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    const [contextMenu, setContextMenu] = useContextMenu();
    const onClick = e => {
        console.log('onClick')
    }
    const onDoubleClick = e => {
        console.log('onDoubleClick')
    }
    let {filePath, branchName} = resourceManager;
    const {Name, Mode} = dirItem;
    filePath = filePath.concat(Name);
    return (
        <div className='file-normal' onContextMenu={(e) => {
            e.preventDefault();
            e.stopPropagation();
            const {clientX, clientY} = e;
            let {x, y, width, height} = filesElm.current.getBoundingClientRect();
            setContextMenu({
                type: 'file',
                dirItem,
                clientX, clientY,
                x, y, width, height,
            })
        }}>
            <div onMouseDown={useClick(null, () => {
                open(sysConfig, setResourceManager, branchName, filePath);
            })}>
                <Icon icon={modeIsDir(Mode) ? 'floderblue' : 'file3'} className='file-icon'/>
            </div>
            <div className='file-name-wrapper'>
                <p kfs-attr="file" className='file-name' onMouseDown={useClick(onClick, onDoubleClick)}>{Name}</p>
            </div>
        </div>
    )
};
