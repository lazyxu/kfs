import { useEffect, useRef } from 'react';

export default function ({ width, height, source }) {
    const ref = useRef();
    useEffect(() => {    
        ref.current?.load();
    }, [source]);
    return (
        <video ref={ref} controls style={{ width, height }} data-setup='{}'>
            <source src={source} />
            您的浏览器不支持 HTML5 video 标签。
        </video>
    );
}
