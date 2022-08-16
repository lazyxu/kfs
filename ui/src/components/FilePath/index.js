import useResourceManager from 'hox/resourceManager';
import './index.scss';

export default () => {
    const [resourceManager] = useResourceManager();
    return (
        <div className='filePath'>
            <div className='pathElement'>
                <div className='pathName'>{resourceManager.branchName}</div>
                <div className='pathNameRight'> > </div>
            </div>
            {resourceManager.filePath.map(elemName => (
                <div key={elemName} className='pathElement'>
                    <div className='pathName'>{elemName}</div>
                    <div className='pathNameRight'> > </div>
                </div>
            ))}
        </div>
    )
};
