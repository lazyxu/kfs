import Video from 'react-native-video';

export default function ({ width, height, source }) {
    return (
        <Video style={{width, height }} source={ source }
        />
    );
}
