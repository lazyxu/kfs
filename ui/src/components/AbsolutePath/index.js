import useResourceManager from 'hox/resourceManager';
import './index.scss';
import PathElement from "./PathElement";

export default () => {
    const [resourceManager] = useResourceManager();
    return (
        <div className='filePath'>
            <PathElement name={resourceManager.branchName} filePath={[]} separator="$"/>
            {resourceManager.filePath.map((elemName, i) =>
                <PathElement key={i} name={elemName} filePath={resourceManager.filePath.slice(0, i + 1)}
                             separator={i === resourceManager.filePath.length - 1 ? "" : ">"}/>
            )}
        </div>
    )
};
