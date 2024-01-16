import { Appbar, Surface } from "react-native-paper";
import ThumbnailListAll from "../Photos/ThumbnailListAll";

export default function ({ navigation, route }) {
    const { driver } = route.params;
    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.BackAction onPress={() => navigation.pop()} />
                <Appbar.Content title={driver.name} />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <ThumbnailListAll metadataList={driver.metadataList} />
        </Surface>
    );
}
