import Video from 'react-native-video';

export default function ({ width, height, uri }) {
    console.log("Video", width, height, uri)
    return (
        <Video style={{ width, height }} source={{ uri }}
        />
    );
}
