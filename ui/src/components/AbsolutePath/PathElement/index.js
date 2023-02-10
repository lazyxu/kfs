import useResourceManager from 'hox/resourceManager';
import {useClick} from "use";
import {open} from "api/api";
import useSysConfig from "hox/sysConfig";

export default ({ name, filePath, separator }) => {
    const [resourceManager, setResourceManager] = useResourceManager();
    const {sysConfig} = useSysConfig();
    return (
        <div className='pathElement' onMouseDown={useClick(() => {
            open(sysConfig, setResourceManager, resourceManager.branchName, filePath);
        })}>
            <div className='pathName'>{name}</div>
            <div className='pathNameRight'>{separator}</div>
        </div>
    )
};
