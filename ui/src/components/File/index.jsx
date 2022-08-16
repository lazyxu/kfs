import './index.scss';
import Icon from "components/Icon/Icon";
import {useClick} from "use";

export default ({name, type}) => {
    const onClick = e => {
        console.log('onClick')
    }
    const onDoubleClick = e => {
        console.log('onDoubleClick')
    }
    return (
        <div className='file-normal'>
            <div onMouseDown={useClick(onClick, onDoubleClick)}>
                <Icon icon={type === 'dir' ? 'floderblue' : 'file3'} className='file-icon'/>
            </div>
            <div className='file-name-wrapper'>
                <p className='file-name' onMouseDown={useClick(onClick, onDoubleClick)}>{name}</p>
            </div>
        </div>
    )
};
