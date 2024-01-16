import { ResizeMode, Video } from 'expo-av';
import * as React from 'react';

export default function App({ width, height, uri }) {
    const video = React.useRef(null);
    return (
        <Video
            ref={video}
            style={{ width, height }}
            source={{ uri }}
            useNativeControls
            resizeMode={ResizeMode.CONTAIN}
        />
    );
}