import './index.scss';
import Icon from "components/Icon/Icon";
import { useClick } from "use";;

import { list } from "rpc/ws";
import useResourceManager from 'hox/resourceManager';

export default ({ name, type }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const onClick = e => {
        console.log('onClick')
    }
    const onDoubleClick = e => {
        console.log('onDoubleClick')
    }
    const open = name => {
        console.log(name);
        (async () => {
            let dirItems;
            let { filePath, branchName } = resourceManager;
            filePath.push(name);
            await list(branchName, filePath, (total) => {
                dirItems = new Array(total);
            }, (dirItem, i) => {
                dirItems[i] = dirItem;
            });
            setResourceManager(prev => {
                return { ...prev, branchName, filePath, dirItems };
            });
        })()
    }
    return (
        <div className='file-normal'>
            <div onMouseDown={useClick(onClick, () => {
                open(name);
            })}>
                <Icon icon={type === 'dir' ? 'floderblue' : 'file3'} className='file-icon' />
            </div>
            <div className='file-name-wrapper'>
                <p className='file-name' onMouseDown={useClick(onClick, onDoubleClick)}>{name}</p>
            </div>
        </div>
    )
};
