import { Appbar, Surface } from "react-native-paper";
import ThumbnailList from '../Photos/ThumbnailList';

export default function ({ navigation, route }) {
    const { list } = route.params;
    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.BackAction onPress={() => navigation.pop()} />
                <Appbar.Content title="自拍" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <ThumbnailList metadataList={list} />
        </Surface>
    );
}
